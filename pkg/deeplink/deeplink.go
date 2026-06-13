// Package deeplink helps parse and build Telegram deep links.
//
// Deep links arrive as /start <payload> messages. Use Parse to extract the
// payload in a /start handler, and Build to generate shareable links.
//
//	b.Handle("/start", func(c telebot.Context) error {
//	    payload := deeplink.Parse(c)
//	    if payload == "" {
//	        return c.Send("Welcome!")
//	    }
//	    return c.Send("You came from: " + payload)
//	})
//
//	link := deeplink.Build("mybot", "ref_42") // → https://t.me/mybot?start=ref_42
package deeplink

import (
	"fmt"
	"strings"

	"gopkg.in/telebot.v3"
)

// Parse extracts the payload from a /start deep link message.
// Returns an empty string for a plain /start with no payload.
func Parse(c telebot.Context) string {
	msg := c.Message()
	if msg == nil {
		return ""
	}
	_, payload, _ := strings.Cut(msg.Text, " ")
	return strings.TrimSpace(payload)
}

// Build returns the deep link URL for the given bot username and payload.
// botUsername may include or omit the leading "@".
func Build(botUsername, payload string) string {
	username := strings.TrimPrefix(botUsername, "@")
	return fmt.Sprintf("https://t.me/%s?start=%s", username, payload)
}
