# Watcher Browser Service

This is a serverless service to be ran on an aws lambda to provide functions to scrape websites

# Development

-   To run locally, first run:
    -   ```
        wget https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie && \
        mkdir ~/.aws-lambda-rie && \
        mv .aws-lambda-rie ~/.aws-lambda-rie/aws-lambda-rie && \
        chmod +x ~/.aws-lambda-rie/aws-lambda-rie
        ```
    -   Build the docker image: `docker run -p 9000:8080 watcher-local-build:latest`
    -   Run the image
        -   ```
            docker run --platform linux/amd64 -v ~/.aws-lambda-rie:/aws-lambda -p 9000:8080 \
                --entrypoint /aws-lambda/aws-lambda-rie \
                watcher-local-build:latest \
                    /src/app
            ```
    -   In another terminal, run `url -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"payload":"hello world!"}'`
    -   Wait for the response!
