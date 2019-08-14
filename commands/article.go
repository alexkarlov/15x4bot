package commands

import (
	"github.com/alexkarlov/15x4bot/store"
)

// article is the command for serving static pages (like documentation, about us, etc)
type article struct {
	name string
	u    *store.User
}

func (c *article) IsEnd() bool {
	return true
}

func (c *article) IsAllow(u *store.User) bool {
	c.u = u
	return true
}

func (c *article) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	article, err := store.LoadArticle(c.name)
	replyMarkup.Text = article.Text
	return replyMarkup, err
}
