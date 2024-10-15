// Package telegram provides a client for interacting with the Telegram Bot API.
// It includes methods to fetch updates from the bot and send messages.
package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"tg_alarm_bot/lib/e"
)

// Client represents a client for the Telegram Bot API.
// It holds the host, base API path, and an HTTP client.
type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	// getUpdatesMethod is the API method name for fetching updates from the bot.
	getUpdatesMethod = "getUpdates"
	// sendMessageMethod is the API method name for sending messages through the bot.
	sendMessageMethod = "sendMessage"
)

// New creates a new Client instance with the provided host and token.
// It sets the base API path using the provided token.
func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// Updates fetches updates (messages, events) from the bot.
// It takes an offset and limit as parameters, representing the message starting point and the number of updates to retrieve.
// Returns a slice of Update objects or an error if the request fails.
func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, e.Wrap("can't get updates", err)
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap("can't get updates", err)
	}

	return res.Result, nil
}

// SendMessage sends a message to a specific chat identified by chatID.
// It takes the chatID and the message text as parameters.
// Returns an error if the message could not be sent.
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

// doRequest performs an HTTP request to the Telegram Bot API.
// It constructs the URL based on the method and query parameters, and returns the response body as bytes or an error.
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}

	return body, nil
}

// newBasePath constructs the base API path by prepending "bot" to the token.
func newBasePath(token string) string {
	return "bot" + token
}
