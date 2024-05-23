# Watcher Browser Service

This is a serverless service to be ran on an aws lambda to provide functions to scrape websites

## Development

### Quick guide

-   First run `make getDebugTools`
-   To build and run locally:

    -   Set `CHROME_PATH` in .env to chrome instillation
    -   Run: `make runDebug` in one terminal
    -   In another terminal run `curl -XPOST "http://localhost:8080/2015-03-31/functions/function/invocations" -d 'JSON HERE'`
    -   Replace with actual json
    -   Wait for the response!

-   To run docker image locally, first run:

    -   Run: `make runDebugDocker`
    -   In another terminal run :`curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d 'JSON HERE'`
    -   Replace with actual json
    -   Wait for the response!

### Json

-   url: url of product
-   price_xpath: xpath to price element
-   image_xpath: xpath to image element
-   in_stock_string: the text when a product is in stock
-   out_of_stock_string: the text when a product is out of stock
-   Empty json: `{"url": "", "price_xpath": "", "image_xpath": "", "in_stock_string": "", "out_of_stock_string": ""}`

Curl full example:

```
curl -XPOST "http://localhost:$PORT/2015-03-31/functions/function/invocations" \
    -d '{"url": "https://www.pokemoncenter.com/en-gb/product/701E11880/bulbasaur-pokemon-soda-pop-plush-5-in", \
        "price_xpath": "/html/body/div[2]/main/div/div[2]/div[2]/p/span", \
        "image_xpath": "/html/body/div[2]/main/div/div[2]/div[1]/div/div/div/div[2]/div/div/div[5]/figure/div/div[1]/img", \
        "in_stock_string": "add to basket", \
        "out_of_stock_string": "out of stock" \
    }' > response.json
```
