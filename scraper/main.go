package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/joho/godotenv"
)

type MyEvent struct {
	URL   string `json:"URL"`
	XPATH string `json:"XPATH"`
}

type MyResponse struct {
	Content string `json:"content"`
	Hash    string `json:"hash"`
	Image   string `json:image`
}

func scrape(ctx context.Context, event *MyEvent) (*MyResponse, error) {
	CHROME_PATH := os.Getenv("CHROME_PATH")
	u := launcher.New().Bin(CHROME_PATH).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(event.URL)
	page.MustWaitStable()
	bytes, err := page.Screenshot(false, &proto.PageCaptureScreenshot{})

	if err != nil {
		log.Fatal(err)
	}

	base64Encoding := "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes)
	return &MyResponse{Content: "Test", Hash: "Test", Image: base64Encoding}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lambda.Start(scrape)
}
