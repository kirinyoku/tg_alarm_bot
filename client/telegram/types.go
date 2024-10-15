package telegram

// UpdatesResponse represents the response structure for the getUpdates API method.
// It contains a boolean Ok to indicate success and a slice of Update objects.
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update represents a single update (message or event) received from the bot.
// It contains the update ID and a pointer to an IncomingMessage structure, which holds the details of the message.
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

// IncomingMessage represents the content of an incoming message in an update.
// It includes the message text, the sender information, and the chat details.
type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

// From represents the sender of the message.
// It contains the username of the sender.
type From struct {
	Username string `json:"username"`
}

// Chat represents the chat from which the message was sent.
// It contains the chat ID, which is used to identify the chat.
type Chat struct {
	ID int `json:"id"`
}
