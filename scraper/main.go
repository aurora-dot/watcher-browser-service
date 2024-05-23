package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

type MyEvent struct {
	Url              string `json:"url"`
	PriceXpath       string `json:"price_xpath"`
	ImageXpath       string `json:"image_xpath"`
	InStockString    string `json:"in_stock_string"`
	OutOfStockString string `json:"out_stock_string"`
}

type MyResponse struct {
	Price   string `json:"price"`
	InStock bool   `json:"in_stock"`
	Image   string `json:"image"`
	HTML    string `json:"html"`
	Error   string `json:"error"`
}

func getStock(page *rod.Page, inStockString string, outOfStockString string) (*bool, error) {
	// We use /%s/i to make the search, case insensitive
	inStockElement := page.MustElementR("button", fmt.Sprintf("/%s/i", inStockString))
	outOfStockElement := page.MustElementR("button", fmt.Sprintf("/%s/i", outOfStockString))

	stockStatus := new(bool)

	if inStockElement != nil && outOfStockElement != nil {
		return stockStatus, errors.New("stock: both in and out of stock")
	} else if inStockElement == nil && outOfStockElement == nil {
		return stockStatus, errors.New("stock: neither in and out of stock")
	}

	if inStockElement != nil {
		*stockStatus = true
	} else if outOfStockElement != nil {
		*stockStatus = false
	} else {
		return stockStatus, errors.New("stock: uncaught error")
	}

	return stockStatus, nil
}

func getImageUrl(page *rod.Page, imageXpath string) (string, error) {
	element := page.MustElementX(imageXpath)

	if element == nil {
		return "", errors.New("image: no such element")
	}

	image := element.MustResource()

	if image == nil {
		return "", errors.New("image: couldn't get image resource")
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(image), nil
}

func scrape(ctx context.Context, event *MyEvent) (*MyResponse, error) {
	CHROME_PATH := os.Getenv("CHROME_PATH")
	u := launcher.New().Bin(CHROME_PATH).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(event.Url)
	page.MustWaitStable()

	stockStatus, err := getStock(page, event.InStockString, event.OutOfStockString)

	if err != nil {
		log.Fatal(err)
		return &MyResponse{}, err
	}

	imageUrl, err := getImageUrl(page, event.ImageXpath)

	if err != nil {
		log.Fatal(err)
		return &MyResponse{}, err
	}

	return &MyResponse{HTML: "Test", Price: "Test", Image: imageUrl, InStock: *stockStatus}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lambda.Start(scrape)
}
