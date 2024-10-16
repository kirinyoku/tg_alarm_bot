// Package source_consumer implements a message consumption system that continuously
// fetches and processes messages from external sources. It is designed to use a Fetcher
// to retrieve messages and a Processor to handle them, with error handling and retry logic.
package source_consumer

import (
	"log"
	"tg_alarm_bot/sources"
	"time"
)

// Consumer represents a structure that fetches and processes messages.
// It relies on an external Fetcher to retrieve the messages and a Processor to handle them.
type Consumer struct {
	fetcher   sources.Fetcher
	processor sources.Processor
}

// New creates a new Consumer instance with the provided Fetcher and Processor.
// It returns a Consumer with both fetcher and processor initialized.
func New(fetcher sources.Fetcher, processor sources.Processor) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
	}
}

// Start begins an infinite loop that continuously fetches and processes messages.
// If an error occurs during fetching or processing, it logs the error and continues.
// If no messages are fetched, it waits for 10 seconds before retrying.
func (c Consumer) Start() error {
	for {
		messages, err := c.fetcher.Fetch()
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(messages) == 0 {
			time.Sleep(10 * time.Second)
			continue
		}

		if err := c.handleMessages(messages); err != nil {
			log.Print(err)
			continue
		}
	}
}

// handleMessages processes each message in the slice using the Processor.
// If processing a message fails, it logs the error and continues with the next message.
func (c *Consumer) handleMessages(messages []sources.Message) error {
	for _, message := range messages {
		if err := c.processor.Process(message); err != nil {
			log.Printf("can't handle message: %s", err.Error())
			continue
		}
	}

	return nil
}
