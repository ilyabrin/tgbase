package deeplink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v3"
)

// fakeCtx wraps telebot.Context with a controlled Message.
type fakeCtx struct {
	telebot.Context
	msg *telebot.Message
}

func (f *fakeCtx) Message() *telebot.Message { return f.msg }

func msgCtx(text string) *fakeCtx {
	return &fakeCtx{msg: &telebot.Message{Text: text}}
}

// --- Parse ---

func TestParse_WithPayload(t *testing.T) {
	assert.Equal(t, "ref_42", Parse(msgCtx("/start ref_42")))
}

func TestParse_PlainStart(t *testing.T) {
	assert.Equal(t, "", Parse(msgCtx("/start")))
}

func TestParse_EmptyPayloadWithSpace(t *testing.T) {
	assert.Equal(t, "", Parse(msgCtx("/start   ")))
}

func TestParse_PayloadWithSpaces(t *testing.T) {
	// Only the first space is the delimiter; rest is part of payload.
	assert.Equal(t, "foo bar", Parse(msgCtx("/start foo bar")))
}

func TestParse_NilMessage(t *testing.T) {
	c := &fakeCtx{msg: nil}
	assert.Equal(t, "", Parse(c))
}

// --- Build ---

func TestBuild_WithAt(t *testing.T) {
	assert.Equal(t, "https://t.me/mybot?start=ref_42", Build("@mybot", "ref_42"))
}

func TestBuild_WithoutAt(t *testing.T) {
	assert.Equal(t, "https://t.me/mybot?start=ref_42", Build("mybot", "ref_42"))
}

func TestBuild_EmptyPayload(t *testing.T) {
	assert.Equal(t, "https://t.me/mybot?start=", Build("mybot", ""))
}
