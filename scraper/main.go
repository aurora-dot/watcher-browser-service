package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-rod/rod"
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
	// u := launcher.New().Bin("/opt/chrome/chrome").MustLaunch()
	page := rod.New().MustConnect().MustPage("https://www.wikipedia.org/")
	page.MustWaitStable().MustScreenshot("a.png")

	return &MyResponse{Content: "Test", Hash: "Test"}, nil
}

func main() {
	lambda.Start(scrape)
}
