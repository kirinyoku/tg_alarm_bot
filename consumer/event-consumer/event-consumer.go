// Package event_consumer defines a Consumer that fetches and processes events in batches.
// It handles the continuous fetching of events and processing them with the provided processor.
package event_consumer

import (
	"log"
	"tg_alarm_bot/events"
	"time"
)

// Consumer is responsible for fetching and processing events.
// It uses a Fetcher to retrieve events and a Processor to handle them.
// The batchSize determines how many events are fetched at a time.
type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

// New creates a new Consumer with the provided Fetcher, Processor, and batchSize.
// The fetcher is used to retrieve events, the processor handles them, and batchSize
// controls the number of events to fetch at once.
func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

// Start begins the continuous loop for fetching and processing events.
// It fetches events in batches, processes each event, and handles errors.
// If no events are fetched, the consumer sleeps for 1 second before trying again.
func (c *Consumer) Start() error {
	for {
		events, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(events) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(events); err != nil {
			log.Print(err)
			continue
		}
	}
}

// handleEvents processes each event in the provided slice of events.
// It logs each new event and attempts to process it. If an error occurs while processing,
// the error is logged and processing continues with the next event.
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %q, %d, %v", event.Text, event.Type, event.Meta)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}

	return nil
}
