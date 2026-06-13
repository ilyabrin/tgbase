package handlers

import (
	"context"
	"fmt"

	"tgbase/internal/database"

	"gopkg.in/telebot.v3"
)

// StarProduct describes a purchasable item priced in Telegram Stars.
type StarProduct struct {
	Title       string
	Description string
	// Payload is your backend identifier returned verbatim in OnPayment.
	// Use it to identify what was bought (e.g. "premium_1month", "coins_100").
	Payload string
	Stars   int // price in Telegram Stars
}

// Validate returns an error if the product is misconfigured.
func (p StarProduct) Validate() error {
	switch {
	case p.Title == "":
		return fmt.Errorf("StarProduct.Title is required")
	case p.Description == "":
		return fmt.Errorf("StarProduct.Description is required")
	case p.Payload == "":
		return fmt.Errorf("StarProduct.Payload is required")
	case p.Stars <= 0:
		return fmt.Errorf("StarProduct.Stars must be > 0")
	}
	return nil
}

// invoice builds the telebot.Invoice for this product.
func (p StarProduct) invoice() telebot.Invoice {
	return telebot.Invoice{
		Title:       p.Title,
		Description: p.Description,
		Payload:     p.Payload,
		Currency:    "XTR", // Telegram Stars currency code
		Token:       "",    // empty = Telegram Stars (no payment provider)
		Prices:      []telebot.Price{{Label: p.Title, Amount: p.Stars}},
	}
}

// SendInvoice sends a Stars invoice to the user.
//
// Register with:
//
//	b.Handle("/buy", handlers.SendInvoice(product))
func SendInvoice(product StarProduct) telebot.HandlerFunc {
	inv := product.invoice()
	return func(c telebot.Context) error {
		_, err := c.Bot().Send(c.Sender(), &inv)
		return err
	}
}

// PreCheckout approves pre-checkout queries automatically.
// Add inventory/validation logic before c.Accept() if needed.
//
// Register with:
//
//	b.Handle(telebot.OnCheckout, handlers.PreCheckout())
func PreCheckout() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		return c.Accept()
	}
}

// PaymentSuccess handles a confirmed Stars payment.
// Records the transaction in the "payments" table and thanks the user.
//
// Required table (run once):
//
//	CREATE TABLE IF NOT EXISTS payments (
//	    id                 SERIAL PRIMARY KEY,  -- or INTEGER for SQLite
//	    user_id            BIGINT NOT NULL,
//	    telegram_charge_id TEXT   NOT NULL UNIQUE,
//	    payload            TEXT   NOT NULL,
//	    stars              INT    NOT NULL,
//	    created_at         TIMESTAMP DEFAULT NOW()  -- or CURRENT_TIMESTAMP
//	);
//
// Register with:
//
//	b.Handle(telebot.OnPayment, handlers.PaymentSuccess(db))
func PaymentSuccess(db database.Database) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		payment := c.Message().Payment
		if payment == nil {
			return nil
		}

		_, err := db.Insert(context.Background(), "payments", map[string]any{
			"user_id":            c.Sender().ID,
			"telegram_charge_id": payment.TelegramChargeID,
			"payload":            payment.Payload,
			"stars":              payment.Total,
		})
		if err != nil {
			return fmt.Errorf("record payment: %w", err)
		}

		return c.Send(fmt.Sprintf("✅ Спасибо! Получено %d ⭐", payment.Total))
	}
}
