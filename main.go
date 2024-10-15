package main

import (
	"flag"
	"log"
	tgClient "tg_alarm_bot/client/telegram"
	event_consumer "tg_alarm_bot/consumer/event-consumer"
	"tg_alarm_bot/events/telegram"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	eventProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()))

	log.Printf("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

// mustToken parses the token flag from the command-line arguments.
// If the token flag ("-t") is not specified, the function logs a fatal error and exits the program.
// If the token is provided, it returns the token as a string.
func mustToken() string {
	token := flag.String("t", "", "token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
