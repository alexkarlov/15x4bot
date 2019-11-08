package commands

import (
	"errors"
	"github.com/alexkarlov/15x4bot/lang"
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
	var err error
	switch c.name {
	case "about":
		replyMarkup.Text = lang.ARTICLE_ABOUT
	case "documentation":
		replyMarkup.Text = lang.ARTICLE_DOCUMENTATION
	case "start":
		replyMarkup.Text = lang.ARTICLE_START
	case "help":
		replyMarkup.Text = lang.ARTICLE_HELP
	case "main_menu":
		replyMarkup.Text = lang.ARTICLE_MAIN_MENU
	default:
		err = errors.New("wrong article name")
	}
	return replyMarkup, err
}

type updateArticle struct {
	name string
	stepConstructor
}

func (c *updateArticle) IsAllow(u *store.User) bool {
	return true
}

func newUpdateArticle(n string) *updateArticle {
	c := &updateArticle{
		name: n,
	}
	c.RegisterSteps(c.firstStep, c.secondStep)
	return c
}

func (c *updateArticle) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	switch c.name {
	case "event_manual_photo_id":
		replyMarkup.Text = lang.MARKUP_BUTTON_ADD_PHOTO_MANUAL
	default:
		err = errors.New("unknown article name")
	}
	return replyMarkup, err
}

func (c *updateArticle) secondStep(answer string) (*ReplyMarkup, error) {
	err := store.UpdateArticle(c.name, answer)
	if err != nil {
		return nil, err
	}
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch c.name {
	case "event_manual_photo_id":
		replyMarkup.Text = lang.MARKUP_BUTTON_PHOTO_MANUAL_SUCCESSFULY_UPDATED
	default:
		err = errors.New("unknown article name")
	}
	return replyMarkup, err
}
