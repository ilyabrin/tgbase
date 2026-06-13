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

// --- Reply keyboard ---

// ReplyBuilder builds a reply keyboard (the persistent button row shown above
// the message input field).
//
// Usage:
//
//	kb := keyboard.Reply().
//	    Row("Profile", "Settings").
//	    Row("Help").
//	    Build()
//	c.Send("Choose:", kb)
//
// To remove the keyboard:
//
//	c.Send("Done!", keyboard.Remove())
type ReplyBuilder struct {
	m    *telebot.ReplyMarkup
	rows []telebot.Row
}

// Reply returns a new ReplyBuilder with resize enabled by default.
func Reply() *ReplyBuilder {
	return &ReplyBuilder{m: &telebot.ReplyMarkup{ResizeKeyboard: true}}
}

// Row appends a row of text buttons.
func (b *ReplyBuilder) Row(texts ...string) *ReplyBuilder {
	btns := make([]telebot.Btn, len(texts))
	for i, t := range texts {
		btns[i] = b.m.Text(t)
	}
	b.rows = append(b.rows, b.m.Row(btns...))
	return b
}

// OneTime hides the keyboard after the first use.
func (b *ReplyBuilder) OneTime() *ReplyBuilder {
	b.m.OneTimeKeyboard = true
	return b
}

// Persistent keeps the keyboard visible between messages.
func (b *ReplyBuilder) Persistent() *ReplyBuilder {
	b.m.IsPersistent = true
	return b
}

// Placeholder sets the input field hint text shown while the keyboard is open.
func (b *ReplyBuilder) Placeholder(text string) *ReplyBuilder {
	b.m.Placeholder = text
	return b
}

// Build finalises the keyboard and returns the *telebot.ReplyMarkup.
func (b *ReplyBuilder) Build() *telebot.ReplyMarkup {
	b.m.Reply(b.rows...)
	return b.m
}

// Remove returns a markup that removes any existing reply keyboard.
func Remove() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{RemoveKeyboard: true}
}
