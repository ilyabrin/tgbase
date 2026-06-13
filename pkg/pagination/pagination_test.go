package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

// --- Pages ---

func TestPages(t *testing.T) {
	cases := []struct {
		total, size, want int
	}{
		{0, 10, 0},
		{10, 0, 0},
		{10, 10, 1},
		{11, 10, 2},
		{20, 10, 2},
		{21, 10, 3},
		{1, 10, 1},
	}
	for _, tc := range cases {
		p := New(tc.total, tc.size)
		assert.Equal(t, tc.want, p.Pages(), "total=%d size=%d", tc.total, tc.size)
	}
}

// --- Offset ---

func TestOffset(t *testing.T) {
	p := New(100, 10)
	assert.Equal(t, 0, p.Offset(1))
	assert.Equal(t, 10, p.Offset(2))
	assert.Equal(t, 20, p.Offset(3))
	assert.Equal(t, 0, p.Offset(0), "page < 1 clamps to 1")
	assert.Equal(t, 0, p.Offset(-5))
}

// --- Keyboard ---

func TestKeyboard_SinglePage_ReturnsNil(t *testing.T) {
	p := New(5, 10)
	assert.Nil(t, p.Keyboard(1, "items"))
}

func TestKeyboard_ZeroTotal_ReturnsNil(t *testing.T) {
	p := New(0, 10)
	assert.Nil(t, p.Keyboard(1, "items"))
}

func TestKeyboard_FirstPage_NoPrev(t *testing.T) {
	p := New(25, 10) // 3 pages
	kb := p.Keyboard(1, "items")
	require.NotNil(t, kb)
	require.Len(t, kb.InlineKeyboard, 1)

	row := kb.InlineKeyboard[0]
	// page 1: no ← , indicator, →
	assert.Len(t, row, 2)
	assert.Equal(t, "1 / 3", row[0].Text)
	assert.Equal(t, "→", row[1].Text)
	assert.Equal(t, "items_next", row[1].Unique)
	assert.Equal(t, "2", row[1].Data)
}

func TestKeyboard_LastPage_NoNext(t *testing.T) {
	p := New(25, 10) // 3 pages
	kb := p.Keyboard(3, "items")
	require.NotNil(t, kb)
	row := kb.InlineKeyboard[0]
	// page 3: ← , indicator
	assert.Len(t, row, 2)
	assert.Equal(t, "←", row[0].Text)
	assert.Equal(t, "items_prev", row[0].Unique)
	assert.Equal(t, "3 / 3", row[1].Text)
}

func TestKeyboard_MiddlePage_BothButtons(t *testing.T) {
	p := New(30, 10) // 3 pages
	kb := p.Keyboard(2, "items")
	require.NotNil(t, kb)
	row := kb.InlineKeyboard[0]
	assert.Len(t, row, 3)
	assert.Equal(t, "←", row[0].Text)
	assert.Equal(t, "2 / 3", row[1].Text)
	assert.Equal(t, "→", row[2].Text)
}

func TestKeyboard_CallbackDataIsTargetPage(t *testing.T) {
	p := New(30, 10)
	kb := p.Keyboard(2, "users")
	row := kb.InlineKeyboard[0]
	// ← button data = "1" (page to go to)
	assert.Equal(t, "1", row[0].Data)
	// → button data = "3"
	assert.Equal(t, "3", row[2].Data)
}

// --- Page ---

type cbCtx struct {
	telebot.Context
	cb *telebot.Callback
}

func (c *cbCtx) Callback() *telebot.Callback { return c.cb }

func TestPage_Valid(t *testing.T) {
	c := &cbCtx{cb: &telebot.Callback{Data: "\fusers_prev|3"}}
	page, err := Page(c)
	require.NoError(t, err)
	assert.Equal(t, 3, page)
}

func TestPage_NoSeparator(t *testing.T) {
	c := &cbCtx{cb: &telebot.Callback{Data: "\fusers_prev"}}
	_, err := Page(c)
	assert.Error(t, err)
}

func TestPage_NilCallback(t *testing.T) {
	c := &cbCtx{cb: nil}
	_, err := Page(c)
	assert.Error(t, err)
}
