package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"sync"
	tg_client "tg_alarm_bot/client/telegram"
	event_consumer "tg_alarm_bot/consumer/event-consumer"
	source_consumer "tg_alarm_bot/consumer/source-consumer"
	"tg_alarm_bot/events/telegram"
	tg_sources "tg_alarm_bot/sources/telegram"
)

const (
	tgBotHost = "api.telegram.org"     // Telegram API host address.
	dataPath  = "./data/channels.json" // Path to the JSON file containing channel configurations.
	batchSize = 100                    // Number of events to process in a single batch.
)

var (
	wg sync.WaitGroup
)

func main() {
	// Create a new Telegram client using the provided bot token.
	tg := tg_client.New(tgBotHost, mustToken())

	// Load the list of channels from the specified JSON file.
	channels, err := loadChannels(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("service started")

	// Initialize the event processor for handling incoming Telegram bot events.
	eventProcessor := telegram.New(tg)
	// Initialize the event consumer to fetch and process events in batches.
	eventConsumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	// Start a goroutine to run the event consumer.
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := eventConsumer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// For each channel, start a source consumer to fetch and process messages.
	for _, c := range channels {
		// Initialize the source processor for handling messages from the channel.
		sourceProcessor := tg_sources.New(c.Name, c.URL, c.SearchRegexp, c.PhrasesToRemove, c.ToChannel, tg)

		// Start a goroutine to run the source consumer.
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Initialize and start the source consumer, and log a fatal error if it fails.
			sourceConsumer := source_consumer.New(sourceProcessor, sourceProcessor)
			if err := sourceConsumer.Start(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	wg.Wait()
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

// loadChannels loads the Telegram channel configurations from a JSON file.
// It reads the file and unmarshals the JSON data into a slice of Source structs.
// Returns the loaded channels or an error if the loading fails.
func loadChannels(filename string) ([]tg_sources.Source, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var channels []tg_sources.Source
	if err := json.Unmarshal(byteValue, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}
