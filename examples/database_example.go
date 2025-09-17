// Package examples shows how to use the database layer
package examples

import (
	"context"
	"fmt"
	"tgbase/internal/database"
	"tgbase/internal/i18n"

	"gopkg.in/telebot.v3"
)

// UserStatsHandler demonstrates how to interact with the database
// This handler stores and retrieves user interaction statistics
func UserStatsHandler(db database.Database, i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		ctx := context.Background()
		userID := c.Sender().ID

		// Example: Create a simple user stats table if it doesn't exist
		createTableQuery := `
			CREATE TABLE IF NOT EXISTS user_stats (
				user_id INTEGER PRIMARY KEY,
				interaction_count INTEGER DEFAULT 0,
				last_interaction TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`

		if _, err := db.Exec(ctx, createTableQuery); err != nil {
			return fmt.Errorf("failed to create user_stats table: %w", err)
		}

		// Update user interaction count
		updateQuery := `
			INSERT INTO user_stats (user_id, interaction_count, last_interaction)
			VALUES ($1, 1, CURRENT_TIMESTAMP)
			ON CONFLICT (user_id) DO UPDATE SET
				interaction_count = user_stats.interaction_count + 1,
				last_interaction = CURRENT_TIMESTAMP`

		if _, err := db.Exec(ctx, updateQuery, userID); err != nil {
			return fmt.Errorf("failed to update user stats: %w", err)
		}

		// Retrieve user stats
		var count int
		selectQuery := `SELECT interaction_count FROM user_stats WHERE user_id = $1`
		row := db.QueryRow(ctx, selectQuery, userID)
		if err := row.Scan(&count); err != nil {
			return fmt.Errorf("failed to get user stats: %w", err)
		}

		lang := c.Sender().LanguageCode
		if lang == "" {
			lang = "en"
		}

		message := i18n.Localize(lang, "user_stats", map[string]any{
			"Count": count,
			"Username": c.Sender().FirstName,
		})

		if message == "" {
			message = fmt.Sprintf("Hello %s! You've interacted with me %d times.",
				c.Sender().FirstName, count)
		}

		return c.Send(message)
	}
}

// UserListHandler demonstrates querying multiple records from the database
func UserListHandler(db database.Database, i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		ctx := context.Background()

		// Get top 5 most active users
		query := `
			SELECT user_id, interaction_count
			FROM user_stats
			ORDER BY interaction_count DESC
			LIMIT 5`

		rows, err := db.Query(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to query user stats: %w", err)
		}
		defer rows.Close()

		var stats []string
		for rows.Next() {
			var userID, count int
			if err := rows.Scan(&userID, &count); err != nil {
				continue
			}
			stats = append(stats, fmt.Sprintf("User %d: %d interactions", userID, count))
		}

		if len(stats) == 0 {
			return c.Send("No user statistics available yet.")
		}

		response := "Top 5 most active users:\n" +
			fmt.Sprintf("%v", stats)

		return c.Send(response)
	}
}