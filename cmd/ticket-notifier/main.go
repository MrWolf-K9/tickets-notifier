package main

import (
	"log"
	"time"

	service "ticket-notifier/cmd/service"
)

var openTickets int = 0

func main() {
	log.Println("Starting ocelari-notifier service")
	service.LoadConfigAndSet("../../config.yaml")
	openTickets = service.GetOpenTickets()
	log.Println("Open tickets ", openTickets)
	for {
		time.Sleep(30 * time.Second)
		updatedOpenTickets := service.GetOpenTickets()
		if updatedOpenTickets == -1 {
			log.Println("Failed to get open tickets number")
			continue
		}
		if updatedOpenTickets == openTickets {
			continue
		}
		if updatedOpenTickets > openTickets {
			service.NotifyUsers()
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
