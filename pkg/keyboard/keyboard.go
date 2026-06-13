// Package keyboard provides a fluent builder for Telegram inline keyboards.
//
// Usage:
//
//	kb := keyboard.New().
//	    Row(keyboard.Btn("Buy", "btn_buy"), keyboard.Btn("Info", "btn_info")).
//	    Row(keyboard.URL("Website", "https://example.com")).
//	    Build()
//	c.Send("Choose an option:", kb)
//
// Register callback handlers with:
//
//	b.Handle("\fbtn_buy", buyHandler)
package keyboard

import "gopkg.in/telebot.v3"

// Btn creates a callback data button.
// unique identifies the callback; data is optional payload passed to the handler.
func Btn(text, unique string, data ...string) telebot.Btn {
	b := telebot.Btn{Text: text, Unique: unique}
	if len(data) > 0 {
		b.Data = data[0]
	}
	return b
}

// URL creates a button that opens a URL.
func URL(text, url string) telebot.Btn {
	return telebot.Btn{Text: text, URL: url}
}

// Builder builds an inline keyboard row by row.
type Builder struct {
	m    *telebot.ReplyMarkup
	rows []telebot.Row
}

// New returns an empty keyboard builder.
func New() *Builder {
	return &Builder{m: &telebot.ReplyMarkup{}}
}

// Row appends a row of buttons.
func (b *Builder) Row(btns ...telebot.Btn) *Builder {
	b.rows = append(b.rows, b.m.Row(btns...))
	return b
}

// Build finalises the keyboard and returns the *telebot.ReplyMarkup ready to
// pass as a send option.
func (b *Builder) Build() *telebot.ReplyMarkup {
	b.m.Inline(b.rows...)
	return b.m
}
