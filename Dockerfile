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

ENTRYPOINT [ "/src/app" ] 
