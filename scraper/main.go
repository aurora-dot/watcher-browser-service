package main

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/joho/godotenv"
)

type MyEvent struct {
	URL                 string `json:"url"`
	PreviousPrice       string `json:"previous_price"`
	PreviousStockStatus string `json:"previous_stock_status"`
	PriceXPATH          string `json:"price_xpath"`
	InStockString       string `json:"in_stock_string"`
	OutStockString      string `json:"out_stock_string"`
}

type MyResponse struct {
	HTML    string `json:"html"`
	Price   string `json:"price"`
	InStock bool   `json:"in_stock"`
	Image   string `json:"image"`
	Error   string `json:"error"`
}

func getStock(page *rod.Page, inStockString string, outStockString string) (*bool, error) {
	inStockElement := page.MustElementR("button", "/"+inStockString+"/i")
	outOfStockElement := page.MustElementR("button", "/"+outStockString+"/i")

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

func takeScreenShot(page *rod.Page) (string, error) {
	bytes, err := page.Screenshot(false, &proto.PageCaptureScreenshot{})

	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes), nil
}

func scrape(ctx context.Context, event *MyEvent) (*MyResponse, error) {
	CHROME_PATH := os.Getenv("CHROME_PATH")
	u := launcher.New().Bin(CHROME_PATH).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(event.URL)
	page.MustWaitStable()

	// check price xpath and get string
	// check stock string and get string

	// some sort of check here so we don't do all this stuff every time if the page hasn't changed

	stockStatus, err := getStock(page, event.InStockString, event.OutStockString)

	if err != nil {
		log.Fatal(err)
		return &MyResponse{}, err
	}

	base64EncodedImage, err := takeScreenShot(page)

	if err != nil {
		log.Fatal(err)
		return &MyResponse{}, err
	}

	return &MyResponse{HTML: "Test", Price: "Test", Image: base64EncodedImage, InStock: *stockStatus}, nil

	// if check didnt pass, send error
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lambda.Start(scrape)
}
