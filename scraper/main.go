package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	"github.com/joho/godotenv"
)

type MyEvent struct {
	Url              string `json:"url"`
	PriceXpath       string `json:"price_xpath"`
	ImageXpath       string `json:"image_xpath"`
	InStockString    string `json:"in_stock_string"`
	OutOfStockString string `json:"out_of_stock_string"`
}

type MyResponse struct {
	Price   string `json:"price"`
	InStock bool   `json:"in_stock"`
	Image   string `json:"image"`
	HTML    string `json:"html"`
	Error   string `json:"error"`
}

func getStock(page *rod.Page, inStockString string, outOfStockString string) (*bool, error) {
	fmt.Println("Started getStock")

	// We use /%s/i to make the search, case insensitive
	inStockElement, inStockErr := page.ElementR("button", fmt.Sprintf("/%s/i", inStockString))
	fmt.Println("Got inStockElement")

	outOfStockElement, outOfStockErr := page.ElementR("button", fmt.Sprintf("/%s/i", outOfStockString))
	fmt.Println("Got outOfStockElement")

	stockStatus := new(bool)

	if inStockErr != nil && outOfStockErr != nil {
		return stockStatus, errors.New("stock: both in and out of stock")
	} else if inStockErr == nil && outOfStockErr == nil {
		return stockStatus, errors.New("stock: neither in and out of stock")
	}

	if inStockElement != nil && inStockErr == nil {
		*stockStatus = true
	} else if outOfStockElement != nil && outOfStockErr == nil {
		*stockStatus = false
	} else {
		return stockStatus, errors.New("stock: uncaught error")
	}

	return stockStatus, nil
}

func getPrice(page *rod.Page, priceXpath string) (string, error) {
	element, err := page.ElementX(priceXpath)

	if err != nil {
		log.Fatal(err)
		return "", errors.New("price: no such element")
	}

	text := element.MustText()

	if text == "" {
		return "", errors.New("price: no text content for element")
	}

	return text, nil

}

func getImageAsBase64(page *rod.Page, imageXpath string) (string, error) {
	element, err := page.ElementX(imageXpath)

	if err != nil {
		log.Fatal(err)
		return "", errors.New("image: no such element")
	}

	image := element.MustResource()

	if image == nil {
		return "", errors.New("image: couldn't get image resource")
	}

	return fmt.Sprintf("data:image/%s;base64,%s", http.DetectContentType(image), base64.StdEncoding.EncodeToString(image)), nil
}

func scrape(ctx context.Context, event *MyEvent) (*MyResponse, error) {

	if event.ImageXpath == "" || event.PriceXpath == "" || event.Url == "" || event.InStockString == "" || event.OutOfStockString == "" {
		return &MyResponse{}, errors.New("request: doesn't have all json attributes")
	}

	fmt.Println("Started")

	CHROME_PATH := os.Getenv("CHROME_PATH")
	u := launcher.New().Bin(CHROME_PATH).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	page := stealth.MustPage(browser)
	page.MustNavigate(event.Url).MustWaitStable()

	fmt.Println("Got page")

	stockStatus, err := getStock(page, event.InStockString, event.OutOfStockString)

	if err != nil {
		log.Fatal(err)
		return &MyResponse{}, err
	}

	fmt.Println("stockStatus")

	price, err := getPrice(page, event.PriceXpath)

	if err != nil {
		log.Fatal(err)
		return &MyResponse{}, err
	}

	fmt.Println("price")

	imageAsBase64, err := getImageAsBase64(page, event.ImageXpath)

	if err != nil {
		log.Fatal(err)
		return &MyResponse{}, err
	}

	fmt.Println("Image: " + imageAsBase64)

	return &MyResponse{HTML: page.MustHTML(), Price: price, Image: imageAsBase64, InStock: *stockStatus}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	lambda.Start(scrape)
}
