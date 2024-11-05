package main

import 
	service "ticket-notifier/cmd/service"

func main() {
	service.LoadConfigAndSet("config.yaml")
	service.NotifyUsers()
}
