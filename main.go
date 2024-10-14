package main

import (
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {}

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
