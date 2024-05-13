FROM ubuntu:20.04

USER root

RUN apt update -y && apt upgrade -y
RUN apt install golang-go -y

WORKDIR /src

COPY scraper/main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o app main.go

RUN rm go.mod go.sum

RUN echo "ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula select true" | debconf-set-selections && \
    apt-get -y -qq install software-properties-common &&\
    apt-add-repository "deb http://archive.canonical.com/ubuntu $(lsb_release -sc) partner" && \
    apt-get -y -qq --no-install-recommends install \
    fontconfig \
    fonts-freefont-ttf \
    fonts-gfs-neohellenic \
    fonts-indic \
    fonts-ipafont-gothic \
    fonts-kacst \
    fonts-liberation \
    fonts-noto-cjk \
    fonts-noto-color-emoji \
    fonts-roboto \
    fonts-thai-tlwg \
    fonts-ubuntu \
    fonts-wqy-zenhei

RUN apt install curl unzip

RUN mkdir -p "/opt/chrome/" \
    && curl -Lo "/opt/chrome/chrome-linux.zip" "https://www.googleapis.com/download/storage/v1/b/chromium-browser-snapshots/o/Linux_x64%2F1299153%2Fchrome-linux.zip?generation=1715336417866122&alt=media" \
    && unzip -q "/opt/chrome/chrome-linux.zip" -d "/opt/chrome/" && mv /opt/chrome/chrome-linux/* /opt/chrome/ \
    && rm -rf /opt/chrome/chrome-linux "/opt/chrome/chrome-linux.zip"

RUN curl -Lo aws-lambda-rie https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie \
    && chmod +x aws-lambda-rie && mv aws-lambda-rie /usr/local/bin/

COPY ./entry_script.sh /entry_script.sh
RUN chmod +x /entry_script.sh

ENTRYPOINT [ "/entry_script.sh","lambda_function.handler" ]

