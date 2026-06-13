// Package pagination builds inline keyboards for paginating DB result sets.
//
// Usage:
//
//	pg := pagination.New(totalUsers, 10) // 10 items per page
//
//	// In handler:
//	c.Send(renderPage(items), pg.Keyboard(page, "users"))
//
//	// Register navigation callbacks:
//	b.Handle("\fusers_prev", usersPageHandler)
//	b.Handle("\fusers_next", usersPageHandler)
//
//	func usersPageHandler(c telebot.Context) error {
//	    page, _ := pagination.Page(c)
//	    items := fetchPage(pg.Offset(page), pg.PageSize)
//	    return c.Edit(renderPage(items), pg.Keyboard(page, "users"))
//	}
package pagination

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

// Pager holds pagination config. Pages are 1-indexed.
type Pager struct {
	Total    int
	PageSize int
}

// New returns a Pager for total items split into pages of pageSize.
func New(total, pageSize int) *Pager {
	return &Pager{Total: total, PageSize: pageSize}
}

// Pages returns the total number of pages (0 when Total or PageSize is 0).
func (p *Pager) Pages() int {
	if p.Total == 0 || p.PageSize == 0 {
		return 0
	}
	return (p.Total + p.PageSize - 1) / p.PageSize
}

// Offset returns the SQL OFFSET value for the given 1-indexed page.
func (p *Pager) Offset(page int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * p.PageSize
}

// Keyboard returns an inline ← N/M → navigation row.
// prefix drives the callback unique names:
//   - ← button: unique = prefix+"_prev", data = target page
//   - → button: unique = prefix+"_next", data = target page
//
// Returns nil when the entire result set fits on one page.
func (p *Pager) Keyboard(current int, prefix string) *telebot.ReplyMarkup {
	pages := p.Pages()
	if pages <= 1 {
		return nil
	}

	m := &telebot.ReplyMarkup{}
	var row []telebot.Btn

	if current > 1 {
		row = append(row, telebot.Btn{
			Text:   "←",
			Unique: prefix + "_prev",
			Data:   strconv.Itoa(current - 1),
		})
	}

	row = append(row, telebot.Btn{
		Text:   fmt.Sprintf("%d / %d", current, pages),
		Unique: prefix + "_cur",
		Data:   strconv.Itoa(current),
	})

	if current < pages {
		row = append(row, telebot.Btn{
			Text:   "→",
			Unique: prefix + "_next",
			Data:   strconv.Itoa(current + 1),
		})
	}

	m.Inline(m.Row(row...))
	return m
}

// Page extracts the target page number from a navigation callback context.
// Telegram callback data has the form "\f{unique}|{page}" - this function
// parses the page from the segment after "|".
func Page(c telebot.Context) (int, error) {
	cb := c.Callback()
	if cb == nil {
		return 0, fmt.Errorf("pagination.Page: not a callback context")
	}
	idx := strings.Index(cb.Data, "|")
	if idx < 0 {
		return 0, fmt.Errorf("pagination.Page: no page in callback data %q", cb.Data)
	}
	return strconv.Atoi(cb.Data[idx+1:])
}
