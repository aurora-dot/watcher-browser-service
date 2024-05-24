FROM ubuntu:20.04

USER root

ARG DEBIAN_FRONTEND=noninteractive
ARG TZ=Europe/London
ARG CHROME_VERSION=1299153

ARG UID=1000
ARG GID=1000

ENV TZ=$TZ
ENV DEBIAN_FRONTEND=$DEBIAN_FRONTEND
ENV LANG="C.UTF-8"
ENV DEBUG_COLORS=true
ENV CHROME_PATH=/task/chrome/chrome
ENV DEBUG=false

RUN apt-get update && apt-get install -y 
RUN apt-get install ca-certificates gnupg -y
RUN apt-get install golang-go -y

# Chrome dependencies
RUN apt-get install -y software-properties-common xvfb libu2f-udev gconf-service \
    libasound2 libatk1.0-0 libc6 libcairo2 libcups2 libdbus-1-3 libexpat1 libfontconfig1 \
    libgcc1 libgconf-2-4 libgdk-pixbuf2.0-0 libglib2.0-0 libgtk-3-0 libnspr4 libpango-1.0-0 \
    libpangocairo-1.0-0 libstdc++6 libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 \
    libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 libxss1 libxtst6 ca-certificates \
    fonts-liberation libappindicator1 libnss3 libnss3-dev lsb-release xdg-utils libgbm-dev

# Fonts
RUN echo "ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula select true" | debconf-set-selections && \
    apt-get -y -qq install software-properties-common &&\
    apt-add-repository "deb http://archive.canonical.com/ubuntu $(lsb_release -sc) partner" && \
    apt-get -y -qq --no-install-recommends install \
    fontconfig fonts-freefont-ttf fonts-gfs-neohellenic fonts-indic fonts-ipafont-gothic \
    fonts-kacst fonts-liberation fonts-noto-cjk fonts-noto-color-emoji fonts-roboto \
    fonts-thai-tlwg fonts-ubuntu fonts-wqy-zenhei

RUN apt-get install curl unzip -y

RUN groupadd -g "${GID}" worker \
    && useradd --create-home --no-log-init -u "${UID}" -g "${GID}" worker
RUN mkdir /task
RUN chown -R worker:worker /task
RUN chmod 755 /task

USER worker

WORKDIR /task

# Get chrome
RUN mkdir -p "/task/chrome/" \
    && curl -Lo "/task/chrome/chrome-linux.zip" "https://www.googleapis.com/download/storage/v1/b/chromium-browser-snapshots/o/Linux_x64%2F${CHROME_VERSION}%2Fchrome-linux.zip?alt=media" \
    && unzip -q "/task/chrome/chrome-linux.zip" -d "/task/chrome/" && mv /task/chrome/chrome-linux/* /task/chrome/ \
    && rm -rf /task/chrome/chrome-linux "/task/chrome/chrome-linux.zip"

RUN echo CHROME_PATH=${CHROME_PATH} > .env

# Copy source and build scraper
COPY scraper/main.go go.mod go.sum ./

RUN GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o app main.go && \
    rm main.go go.mod go.sum

ENTRYPOINT [ "task/app" ]
