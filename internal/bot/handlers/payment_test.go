package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"

	"tgbase/internal/database"
)

// --- StarProduct ---

func TestStarProduct_Validate(t *testing.T) {
	valid := StarProduct{
		Title:       "Premium",
		Description: "Unlock premium features",
		Payload:     "premium_1month",
		Stars:       100,
	}

	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, valid.Validate())
	})
	t.Run("empty title", func(t *testing.T) {
		p := valid
		p.Title = ""
		assert.Error(t, p.Validate())
	})
	t.Run("empty description", func(t *testing.T) {
		p := valid
		p.Description = ""
		assert.Error(t, p.Validate())
	})
	t.Run("empty payload", func(t *testing.T) {
		p := valid
		p.Payload = ""
		assert.Error(t, p.Validate())
	})
	t.Run("zero stars", func(t *testing.T) {
		p := valid
		p.Stars = 0
		assert.Error(t, p.Validate())
	})
	t.Run("negative stars", func(t *testing.T) {
		p := valid
		p.Stars = -1
		assert.Error(t, p.Validate())
	})
}

func TestStarProduct_Invoice(t *testing.T) {
	p := StarProduct{Title: "Premium", Description: "Unlock", Payload: "prem_1m", Stars: 50}
	inv := p.invoice()

	assert.Equal(t, "Premium", inv.Title)
	assert.Equal(t, "Unlock", inv.Description)
	assert.Equal(t, "prem_1m", inv.Payload)
	assert.Equal(t, "XTR", inv.Currency)
	assert.Empty(t, inv.Token, "Stars invoices must have empty provider token")
	require.Len(t, inv.Prices, 1)
	assert.Equal(t, 50, inv.Prices[0].Amount)
	assert.Equal(t, "Premium", inv.Prices[0].Label)
}

// --- Handler constructors ---

func TestSendInvoice_ReturnsHandler(t *testing.T) {
	h := SendInvoice(StarProduct{Title: "T", Description: "D", Payload: "P", Stars: 1})
	assert.NotNil(t, h)
}

func TestPreCheckout_ReturnsHandler(t *testing.T) {
	assert.NotNil(t, PreCheckout())
}

func TestPaymentSuccess_ReturnsHandler(t *testing.T) {
	assert.NotNil(t, PaymentSuccess(nil))
}

// --- PaymentSuccess DB insertion ---

type mockPaymentContext struct {
	telebot.Context
	senderID int64
	payment  *telebot.Payment
	sent     []string
}

func (m *mockPaymentContext) Sender() *telebot.User { return &telebot.User{ID: m.senderID} }
func (m *mockPaymentContext) Message() *telebot.Message {
	return &telebot.Message{Payment: m.payment}
}
func (m *mockPaymentContext) Send(what interface{}, _ ...interface{}) error {
	if s, ok := what.(string); ok {
		m.sent = append(m.sent, s)
	}
	return nil
}

func newTestDB(t *testing.T) database.Database {
	t.Helper()
	db, err := database.NewSQLiteDB(context.Background(), ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(context.Background(), `
		CREATE TABLE payments (
			id                 INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id            INTEGER NOT NULL,
			telegram_charge_id TEXT    NOT NULL UNIQUE,
			payload            TEXT    NOT NULL,
			stars              INTEGER NOT NULL
		)`)
	require.NoError(t, err)
	return db
}

func TestPaymentSuccess_InsertsAndReplies(t *testing.T) {
	db := newTestDB(t)
	ctx := &mockPaymentContext{
		senderID: 42,
		payment: &telebot.Payment{
			TelegramChargeID: "tg_charge_abc",
			Payload:          "prem_1m",
			Total:            100,
		},
	}

	handler := PaymentSuccess(db)
	require.NoError(t, handler(ctx))

	// Verify row was inserted.
	row := db.QueryRow(context.Background(),
		"SELECT user_id, telegram_charge_id, payload, stars FROM payments WHERE user_id = ?", 42)
	var uid int64
	var chargeID, payload string
	var stars int
	require.NoError(t, row.Scan(&uid, &chargeID, &payload, &stars))
	assert.Equal(t, int64(42), uid)
	assert.Equal(t, "tg_charge_abc", chargeID)
	assert.Equal(t, "prem_1m", payload)
	assert.Equal(t, 100, stars)

	// Verify confirmation message contains the star count.
	require.Len(t, ctx.sent, 1)
	assert.Contains(t, ctx.sent[0], "100")
}

func TestPaymentSuccess_NilPayment(t *testing.T) {
	ctx := &mockPaymentContext{senderID: 1, payment: nil}
	assert.NoError(t, PaymentSuccess(nil)(ctx))
	assert.Empty(t, ctx.sent)
}
