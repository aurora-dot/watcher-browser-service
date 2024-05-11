FROM browserless/chrome

USER root

RUN apt update -y && apt upgrade -y
RUN apt install golang-go -y

WORKDIR /src

COPY scraper/main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o app main.go


RUN mkdir -p ~/.aws-lambda-rie && \
    curl -Lo ~/.aws-lambda-rie/aws-lambda-rie https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie && \
    chmod +x ~/.aws-lambda-rie/aws-lambda-rie
# RUN yum install docker

# # Install ATK from CentOS 7
# RUN rpm -ivh --nodeps http://mirror.centos.org/centos/7/os/x86_64/Packages/atk-2.28.1-2.el7.x86_64.rpm
# RUN rpm -ivh --nodeps http://mirror.centos.org/centos/7/os/x86_64/Packages/at-spi2-atk-2.26.2-1.el7.x86_64.rpm
# RUN rpm -ivh --nodeps http://mirror.centos.org/centos/7/os/x86_64/Packages/at-spi2-core-2.28.0-1.el7.x86_64.rpm
# RUN rpm -ivh --nodeps  http://mirror.centos.org/centos/7/os/x86_64/Packages/mesa-libgbm-18.3.4-10.el7.x86_64.rpm
# RUN rpm -ivh --nodeps   http://mirror.centos.org/centos/7/os/x86_64/Packages/libwayland-server-1.15.0-1.el7.x86_64.rpm
# RUN rpm -ivh --nodeps  http://mirror.centos.org/centos/7/os/x86_64/Packages/glibc-2.17-317.el7.x86_64.rpm

# yum install pango-devel.x86_64 pango.x86_64
# yum install libXrandr.x86_64 libXrandr-devel.x86_64


ENTRYPOINT [ "/src/app" ] 
