package keyboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBtn(t *testing.T) {
	b := Btn("Buy", "btn_buy")
	assert.Equal(t, "Buy", b.Text)
	assert.Equal(t, "btn_buy", b.Unique)
	assert.Equal(t, "", b.Data)
}

func TestBtn_WithData(t *testing.T) {
	b := Btn("Item", "btn_item", "item_42")
	assert.Equal(t, "Item", b.Text)
	assert.Equal(t, "btn_item", b.Unique)
	assert.Equal(t, "item_42", b.Data)
}

func TestURL(t *testing.T) {
	b := URL("Website", "https://example.com")
	assert.Equal(t, "Website", b.Text)
	assert.Equal(t, "https://example.com", b.URL)
	assert.Equal(t, "", b.Unique)
}

func TestBuilder_Build_NotNil(t *testing.T) {
	kb := New().
		Row(Btn("Yes", "btn_yes"), Btn("No", "btn_no")).
		Build()
	require.NotNil(t, kb)
	require.NotNil(t, kb.InlineKeyboard)
}

func TestBuilder_SingleRow(t *testing.T) {
	kb := New().
		Row(Btn("A", "a"), Btn("B", "b")).
		Build()
	assert.Len(t, kb.InlineKeyboard, 1)
	assert.Len(t, kb.InlineKeyboard[0], 2)
}

func TestBuilder_MultipleRows(t *testing.T) {
	kb := New().
		Row(Btn("Row1Col1", "r1c1"), Btn("Row1Col2", "r1c2")).
		Row(Btn("Row2Col1", "r2c1")).
		Row(URL("Site", "https://example.com")).
		Build()
	assert.Len(t, kb.InlineKeyboard, 3)
	assert.Len(t, kb.InlineKeyboard[0], 2)
	assert.Len(t, kb.InlineKeyboard[1], 1)
	assert.Len(t, kb.InlineKeyboard[2], 1)
}

func TestBuilder_Empty(t *testing.T) {
	kb := New().Build()
	require.NotNil(t, kb)
}

func TestBuilder_Chaining(t *testing.T) {
	b := New()
	result := b.Row(Btn("A", "a"))
	assert.Equal(t, b, result, "Row should return the same builder for chaining")
}

// --- Reply keyboard ---

func TestReply_Build_NotNil(t *testing.T) {
	kb := Reply().Row("Profile", "Settings").Row("Help").Build()
	require.NotNil(t, kb)
	require.NotNil(t, kb.ReplyKeyboard)
}

func TestReply_RowCount(t *testing.T) {
	kb := Reply().Row("A", "B").Row("C").Build()
	assert.Len(t, kb.ReplyKeyboard, 2)
	assert.Len(t, kb.ReplyKeyboard[0], 2)
	assert.Len(t, kb.ReplyKeyboard[1], 1)
}

func TestReply_ButtonText(t *testing.T) {
	kb := Reply().Row("Hello").Build()
	assert.Equal(t, "Hello", kb.ReplyKeyboard[0][0].Text)
}

func TestReply_ResizeEnabledByDefault(t *testing.T) {
	kb := Reply().Build()
	assert.True(t, kb.ResizeKeyboard)
}

func TestReply_OneTime(t *testing.T) {
	kb := Reply().OneTime().Build()
	assert.True(t, kb.OneTimeKeyboard)
}

func TestReply_Persistent(t *testing.T) {
	kb := Reply().Persistent().Build()
	assert.True(t, kb.IsPersistent)
}

func TestReply_Placeholder(t *testing.T) {
	kb := Reply().Placeholder("Type here...").Build()
	assert.Equal(t, "Type here...", kb.Placeholder)
}

func TestReply_Chaining(t *testing.T) {
	b := Reply()
	result := b.Row("A")
	assert.Equal(t, b, result)
}

func TestRemove(t *testing.T) {
	kb := Remove()
	require.NotNil(t, kb)
	assert.True(t, kb.RemoveKeyboard)
}
