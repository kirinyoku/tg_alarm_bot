package telegram

// UpdatesResponse represents the response structure for the getUpdates API method.
// It contains a boolean Ok to indicate success and a slice of Update objects.
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update represents a single update (message, event) from the bot.
// It contains the update ID and the message text.
type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}
