// Package events defines interfaces and types used for fetching and processing events,
// as well as the structure of an event and its types.
package events

// Fetcher defines an interface for fetching events.
// The Fetch method accepts a limit and returns a slice of Event objects and an error if any.
type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

// Processor defines an interface for processing events.
// The Process method accepts an Event and returns an error if the processing fails.
type Processor interface {
	Process(e Event) error
}

// Type represents the type of an event. It is an enumerated integer type.
type Type int

const (
	// Unknown is the default event type, representing an event with an unrecognized or undefined type.
	Unknown Type = iota
	// Message represents an event of type Message, typically used for text-based messages.
	Message
)

// Event represents a generic event with a type, text content, and additional metadata.
// The Type field specifies the event type, the Text field contains the event message,
// and the Meta field can hold additional contextual information.
type Event struct {
	Type Type
	Text string
	Meta interface{}
}
