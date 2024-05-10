FROM public.ecr.aws/lambda/go:1 as lambda

# Copy over go source code
COPY scraper/main.go /src/main.go
COPY go.mod /src/go.mod
COPY go.sum /src/go.sum

WORKDIR /src

# Install packages
RUN yum install go xz atk cups-libs gtk3 libXcomposite alsa-lib tar \
    libXcursor libXdamage libXext libXi libXrandr libXScrnSaver \
    libXtst pango at-spi2-atk libXt xorg-x11-server-Xvfb \
    xorg-x11-xauth dbus-glib dbus-glib-devel unzip bzip2 -y -q

# Build go lambda function
RUN GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ${LAMBDA_TASK_ROOT} /src/main.go

# Install chrome
RUN mkdir -p /opt/chrome/
RUN curl -Lo "/opt/chrome/chrome-linux.zip" "https://www.googleapis.com/download/storage/v1/b/chromium-browser-snapshots/o/Linux_x64%2F1299153%2Fchrome-linux.zip?generation=1715336417866122&alt=media"
RUN unzip -q "/opt/chrome/chrome-linux.zip" -d "/opt/chrome/"
RUN mv /opt/chrome/chrome-linux/* /opt/chrome/
RUN rm -rf /opt/chrome/chrome-linux "/opt/chrome/chrome-linux.zip"

CMD [ "app" ]
