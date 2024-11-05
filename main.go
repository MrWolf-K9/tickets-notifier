package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"

	client "go.vxn.dev/sms-manager/pkg/client"
)

var (
	apiKey       string   = ""
	ticketsUrl   string   = ""
	openTickets  int      = 0
	text         string   = ""
	numbers      []string = []string{""}
	re                    = regexp.MustCompile("")
)

func main() {
	log.Println("Starting ocelari-notifier service")
	openTickets = getOpenTickets()
	log.Println("Open tickets ", openTickets)
	for {
		time.Sleep(30 * time.Second)
		updatedOpenTickets := getOpenTickets()
		if updatedOpenTickets == -1 {
			log.Println("Failed to get open tickets number")
			continue
		}
		if updatedOpenTickets == openTickets {
			continue
		}
		if updatedOpenTickets > openTickets {
			notifyUsers()
			openTickets = updatedOpenTickets
			log.Println("New ticket!! going from {} to {} number of tickets", openTickets, updatedOpenTickets)
			// TODO: add more info, date time
			continue
		}
		if updatedOpenTickets < openTickets {
			openTickets = updatedOpenTickets
			log.Println("Some tickets were removed, going from {} to {} number of tickets", openTickets, updatedOpenTickets)
		}

	}
}

func notifyUsers() {
	for _, num := range numbers {
		notify(num, text)
	}
}

func getOpenTickets() int {
	webPageData := getWebPageData(ticketsUrl)
	tickets := countDivsWithClass(webPageData, "EventRowItem_content_logo")
	return tickets
}

func countDivsWithClass(pageContent string, targetClass string) int {
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
	req := client.Request{
		APIKey:      apiKey,
		PhoneNumber: number,
		Message:     text,
		GatewayType: "high",
		// Sender:      "Wolf steel-notify",
	}
	resp := &client.Response{}

	if err := client.DoRequest(req, "Send", resp); err != nil {
		log.Fatal(err)
	}
	if resp.Message != "ERROR" {
		fmt.Println("status         : " + resp.Message)
		fmt.Printf("request ID     : %s\n\r", resp.RequestID)
		fmt.Println("phone number(s): " + resp.PhoneNumber)
		fmt.Printf("custom ID      : %d\n\r", resp.CustomID)
	} else {
		fmt.Printf("response error code: %d\n\r", resp.ErrorCode)
		fmt.Println("response error msg: " + resp.ErrorMessage)
	}
}
