# Watcher Browser Service

This is a serverless service to be ran on an aws lambda to provide functions to scrape websites

# Development

-   First run `make getDebugTools`
-   To build and run locally:
    -   Set `CHROME_PATH` in .env to chrome instillation
    -   Run: `make runDebug` in one terminal
    -   In another terminal run `curl -XPOST "http://localhost:8080/2015-03-31/functions/function/invocations" -d '{"Name": "World"}'`
    -   Wait for the response!
-   To run docker image locally, first run:
    -   Build the docker image: `docker run -p 9000:8080 watcher-local-build:latest`
    -   Run the image: `make runDebugDocker`
    -   In another terminal, run `url -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"payload":"hello world!"}'`
    -   Wait for the response!
