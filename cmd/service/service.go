package service

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"

	client "go.vxn.dev/sms-manager/pkg/client"
)

var (
	apiKey     string   = ""
	ticketsUrl string   = ""
	text       string   = ""
	numbers    []string = []string{""}
	re                  = regexp.MustCompile("")
)

func LoadConfigAndSet(path string) {
	log.Println("Loading config from config.yaml")
	config, err := LoadConfig(path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	apiKey = config.APIKey
	ticketsUrl = config.TicketsURL
	text = config.MessageText
	numbers = config.PhoneNumbers
	re = regexp.MustCompile(config.RegexString)

	log.Println("Configuration loaded")
}

func NotifyUsers() {
	for _, num := range numbers {
		notify(num, text)
	}
}

func GetOpenTickets() int {
	webPageData := getWebPageData(ticketsUrl)
	tickets := countDivsWithClass(webPageData)
	return tickets
}

func countDivsWithClass(pageContent string) int {
	// return strings.Count(pageContent, targetClass)
	matches := re.FindAllStringSubmatch(pageContent, -1)
	return len(matches)
}

func getWebPageData(url string) string {
	log.Println("Getting webPage data: ", url)
	// Set up context with timeout
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout for the task
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var htmlContent string

	// Run tasks
	err := chromedp.Run(ctx,
		chromedp.Navigate(url), // Replace with the URL you want to fetch
		// TODO: optimalize
		chromedp.Sleep(15*time.Second),
		chromedp.WaitReady("body", chromedp.ByQuery), // Wait until the body tag is ready
		chromedp.OuterHTML("html", &htmlContent),     // Get the full HTML content
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("WebPage data obtained")

	// Print the fetched HTML content
	return htmlContent
}

func notify(number string, text string) {
	log.Println("Sending notification sms to ", number)
	req := client.Request{
		APIKey:      apiKey,
		PhoneNumber: number,
		Message:     text,
		GatewayType: "high",
	}
	resp := &client.Response{}

	if err := client.DoRequest(req, "Send", resp); err != nil {
		log.Fatal(err)
	}
	if resp.Message != "ERROR" {
		log.Println("status         : " + resp.Message)
		log.Printf("request ID     : %s\n\r", resp.RequestID)
		log.Println("phone number(s): " + resp.PhoneNumber)
		log.Printf("custom ID      : %d\n\r", resp.CustomID)
	} else {
		log.Printf("response error code: %d\n\r", resp.ErrorCode)
		log.Println("response error msg: " + resp.ErrorMessage)
	}
}
