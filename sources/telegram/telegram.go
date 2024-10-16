// Package telegram provides functionality for fetching and processing messages from public Telegram channels.
// It includes tools to filter messages based on specific patterns, clean up unwanted content, and forward the
// filtered messages to a designated Telegram channel via a Telegram client.
package telegram

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"tg_alarm_bot/client/telegram"
	"tg_alarm_bot/lib/e"
	"tg_alarm_bot/sources"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

// Source represents a Telegram source that fetches and processes messages.
// It includes configuration for fetching, filtering, and sending messages to a specific Telegram channel.
type Source struct {
	Name            string               `json:"name"`              // Name of the source.
	URL             string               `json:"url"`               // URL of the Telegram public channel.
	SearchRegexp    string               `json:"search_regexp"`     // Regular expression to search for specific patterns in messages.
	PhrasesToRemove []string             `json:"phrases_to_remove"` // List of phrases to remove from the messages before sending.
	ToChannel       int                  `json:"to_channel"`        // ID of the destination Telegram channel to forward messages to.
	seen            map[string]time.Time `json:"-"`                 // Map of seen messages with their timestamp to avoid duplicates.
	expiry          time.Duration        `json:"-"`                 // Expiry duration for messages to be considered 'seen'.
	tg              *telegram.Client     `json:"-"`                 // Telegram client to send messages.
	startTime       time.Time            `json:"-"`                 // Time when the source started, used to filter old messages.
}

// New creates a new Source instance with the provided parameters.
// It initializes the seen map and sets the expiry duration to 24 hours by default.
func New(name string, url string, pattern string, phrases []string, to int, tg *telegram.Client) *Source {
	return &Source{
		Name:            name,
		URL:             url,
		SearchRegexp:    pattern,
		PhrasesToRemove: phrases,
		ToChannel:       to,
		seen:            make(map[string]time.Time),
		expiry:          24 * time.Hour,
		tg:              tg,
		startTime:       time.Now(),
	}
}

// Fetch retrieves and filters messages from the Telegram source URL.
// It uses an HTTP GET request to fetch the data and filters the messages based on the search regular expression.
// Old messages in the 'seen' map that exceed the expiry time are deleted.
func (s *Source) Fetch() ([]sources.Message, error) {
	res, err := http.Get(s.URL)
	if err != nil {
		return nil, e.Wrap("can't fetch data from telegram source", err)
	}

	defer res.Body.Close()

	messages, err := s.filter(res.Body)
	if err != nil {
		return nil, e.Wrap("can't fetch data from telegram source", err)
	}

	for id, timestamp := range s.seen {
		if time.Since(timestamp) > s.expiry {
			delete(s.seen, id)
		}
	}

	return messages, nil
}

// Process sends a given message to the configured Telegram channel using the Telegram client.
// It forwards the message text to the target channel.
func (s *Source) Process(message sources.Message) error {
	return s.tg.SendMessage(
		s.ToChannel,
		message.Text,
	)
}

// filter parses the HTML content of the Telegram page and extracts messages matching the search pattern.
// It uses goquery to select and process message elements and returns a list of matching messages.
func (s *Source) filter(r io.Reader) ([]sources.Message, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var messages []sources.Message

	// Find and iterate over message elements in the HTML document.
	doc.Find(".tgme_widget_message").Each(func(i int, sel *goquery.Selection) {
		messageID, _ := sel.Attr("data-post")
		messageText := sel.Find(".tgme_widget_message_text").First().Text()

		// If the message has a reply, extract the text from the reply block.
		if replyBlock := sel.Find(".tgme_widget_message_reply"); replyBlock.Length() > 0 {
			messageText = replyBlock.Next().Find(".tgme_widget_message_text").Text()
		}

		// Extract and parse the message timestamp.
		postTime, _ := sel.Find(".tgme_widget_message_date time").Attr("datetime")
		parsedTime, err := time.Parse(time.RFC3339, postTime)
		if err != nil || parsedTime.Before(s.startTime) {
			// Skip the message if the timestamp is invalid or before the start time.
			return
		}

		// Compile the search regular expression and check if the message text matches.
		rx := regexp.MustCompile(s.SearchRegexp)
		if rx.MatchString(messageText) && utf8.RuneCountInString(messageText) < 150 {
			// If the message is new or expired, add it to the list of messages and mark it as seen.
			if _, exists := s.seen[messageID]; !exists || time.Since(s.seen[messageID]) > s.expiry {
				s.seen[messageID] = time.Now()
				messages = append(messages, sources.Message{
					ID:   messageID,
					Text: s.cleanMessage(messageText),
				})
			}
		}
	})

	return messages, nil
}

// cleanMessage removes unwanted phrases from the message text and trims whitespace.
// It iterates over the PhrasesToRemove and applies them to the message.
func (s *Source) cleanMessage(text string) string {
	for _, phrase := range s.PhrasesToRemove {
		text = strings.ReplaceAll(text, phrase, "")
	}

	return strings.TrimSpace(text)
}
