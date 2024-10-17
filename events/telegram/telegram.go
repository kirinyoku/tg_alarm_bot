// Package telegram provides a processor for handling events and interacting with the Telegram Bot API.
// It includes functions for fetching and processing events and converting Telegram updates into internal event types.
package telegram

import (
	"errors"
	"tg_alarm_bot/client/telegram"
	"tg_alarm_bot/events"
	"tg_alarm_bot/lib/e"
)

// Processor handles fetching and processing of Telegram updates.
// It maintains the client for Telegram communication and the current update offset.
type Processor struct {
	tg     *telegram.Client
	offset int
}

// Meta contains metadata for a message, including the chat ID and the username of the sender.
type Meta struct {
	ChatID   int
	Username string
}

var (
	// ErrUnknownEventType is returned when an event with an unrecognized type is encountered.
	ErrUnknownEventType = errors.New("unknown event type")
	// ErrUnknownMetaType is returned when the event's meta field cannot be cast to the expected Meta type.
	ErrUnknownMetaType = errors.New("unknown meta type")
)

// New creates a new Processor with the provided Telegram client.
func New(client *telegram.Client) *Processor {
	return &Processor{
		tg: client,
	}
}

// Fetch retrieves a list of events by fetching updates from the Telegram Bot API.
// It returns a slice of events and updates the offset to process subsequent events.
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, utoe(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

// Process processes a single event by checking its type and handling it accordingly.
// Currently, it only supports processing message events.
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process event", ErrUnknownEventType)
	}
}

// processMessage handles the processing of message events.
// It retrieves metadata from the event and sends a response message using the Telegram client.
func (p *Processor) processMessage(event events.Event) error {
	m, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	return p.tg.SendMessage(m.ChatID, "This bot does not interact directly.", "")
}

// meta extracts metadata from the event's Meta field and casts it to the Meta type.
// Returns an error if the cast is unsuccessful.
func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("cant't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

// utoe (Update to Event) converts a Telegram update into an internal Event structure.
// It maps the update type and message content to the corresponding fields in the Event struct.
func utoe(u telegram.Update) events.Event {
	uType := fetchType(u)

	res := events.Event{
		Type: uType,
		Text: fetchText(u),
	}

	if uType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}

	return res
}

// fetchType determines the event type based on the content of the Telegram update.
// If the update contains a message, it returns events.Message; otherwise, it returns events.Unknown.
func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}

	return events.Message
}

// fetchText retrieves the text content from the Telegram update.
// If the update does not contain a message, it returns an empty string.
func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}

	return u.Message.Text
}
