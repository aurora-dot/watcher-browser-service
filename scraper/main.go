package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

type MyEvent struct {
	URL   string `json:"URL"`
	XPATH string `json:"XPATH"`
}

type MyResponse struct {
	Content string `json:"content"`
	Hash    string `json:"hash"`
}

func scrape(ctx context.Context, event *MyEvent) (*MyResponse, error) {
	CHROME_PATH := os.Getenv("CHROME_PATH")
	u := launcher.New().Bin(CHROME_PATH).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage("https://www.wikipedia.org/")
	page.MustWaitStable().MustScreenshot("a.png")

	return &MyResponse{Content: "Test", Hash: "Test"}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lambda.Start(scrape)
}
