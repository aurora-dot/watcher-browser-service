package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
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
}

func getStock(page *rod.Page, inStockString string, outOfStockString string) (*bool, error) {
	hasInStockElement, _, inStockErr := page.HasR("button", fmt.Sprintf("/%s/i", inStockString))
	hasOutOfStockElement, _, outOfStockErr := page.HasR("button", fmt.Sprintf("/%s/i", outOfStockString))

	stockStatus := new(bool)

	if inStockErr != nil {
		log.Println(inStockErr)
		return stockStatus, errors.New("stock: internal in stock error")
	}

	if outOfStockErr != nil {
		log.Println(outOfStockErr)
		return stockStatus, errors.New("stock: internal out of stock error")
	}

	if hasInStockElement && hasOutOfStockElement {
		return stockStatus, errors.New("stock: both in and out of stock")
	} else if !hasInStockElement && !hasOutOfStockElement {
		return stockStatus, errors.New("stock: neither in or out of stock, this could be due to being redirected to their 'verify you are not a robot' page")
	}

	if hasInStockElement && !hasOutOfStockElement {
		*stockStatus = true
	} else if !hasInStockElement && hasOutOfStockElement {
		*stockStatus = false
	} else {
		return stockStatus, errors.New("stock: uncaught error")
	}

	return stockStatus, nil
}

func getPrice(page *rod.Page, priceXpath string) (string, error) {
	hasElement, element, err := page.HasX(priceXpath)

	if err != nil {
		log.Println(err)
		return "", errors.New("price: internal price error")
	}

	if !hasElement {
		return "", errors.New("price: cannot fetch element from xpath")
	}

	text := element.MustText()

	if text == "" {
		return "", errors.New("price: no text content for element")
	}

	return text, nil

}

func getImageAsBase64(page *rod.Page, imageXpath string) (string, error) {
	hasElement, element, err := page.HasX(imageXpath)

	if err != nil {
		log.Println(err)
		return "", errors.New("image: internal image error")
	}

	if !hasElement {
		return "", errors.New("image: cannot fetch element from xpath")
	}

	image := element.MustResource()

	if image == nil {
		return "", errors.New("image: couldn't get image resource")
	}

	return fmt.Sprintf("data:image/%s;base64,%s", http.DetectContentType(image), base64.StdEncoding.EncodeToString(image)), nil
}

func setupBrowser() *rod.Page {
	CHROME_PATH := os.Getenv("CHROME_PATH")
	browserArgs := launcher.New().
		UserDataDir("/tmp/profile").
		Leakless(true).
		Devtools(false).
		Headless(true).
		NoSandbox(true).
		Set("--no-zygote").
		Set("--disable-dev-shm-usage").
		Set("--disable-setuid-sandbox").
		Set("--disable-dev-shm-usage").
		Set("--disable-gpu").
		Set("--no-zygote").
		Set("--single-process")

	wsURL := browserArgs.Bin(CHROME_PATH).MustLaunch()
	browser := rod.New().ControlURL(wsURL).MustConnect().DefaultDevice(devices.IPhoneX)
	page := stealth.MustPage(browser)
	page.MustSetExtraHeaders(
		"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language", "en-GB,en;q=0.5",
		"Accept-Encoding", "gzip, deflate, br, zstd",
	)
	return page
}

func scrape(ctx context.Context, event *MyEvent) (*MyResponse, error) {
	if event.ImageXpath == "" || event.PriceXpath == "" || event.Url == "" || event.InStockString == "" || event.OutOfStockString == "" {
		return &MyResponse{}, errors.New("request: doesn't have all json attributes")
	}

	log.Println("Started scrape")

	page := setupBrowser()
	page.MustNavigate(event.Url).MustWaitStable()

	log.Println("Got page")

	DEBUG := strings.ToLower(os.Getenv("DEBUG"))
	if DEBUG == "true" {
		if err := os.WriteFile("page.html", []byte(page.MustHTML()), 0666); err != nil {
			log.Fatal(err)
		}

	}

	stockStatus, err := getStock(page, event.InStockString, event.OutOfStockString)

	if err != nil {
		log.Println(err)
		return &MyResponse{}, err
	}

	log.Println("Finished getStock")

	price, err := getPrice(page, event.PriceXpath)

	if err != nil {
		log.Println(err)
		return &MyResponse{}, err
	}

	log.Println("Finished price")

	imageAsBase64, err := getImageAsBase64(page, event.ImageXpath)

	if err != nil {
		log.Println(err)
		return &MyResponse{}, err
	}

	log.Println("Finished getImageAsBase64")

	page.MustClose()

	return &MyResponse{HTML: page.MustHTML(), Price: price, Image: imageAsBase64, InStock: *stockStatus}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lambda.Start(scrape)
}
