package commands

import (
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"time"
)

// profile is the command for updating name
type profile struct {
	field string
	u     *store.User
	stepConstructor
}

// newProfile creates profile and registers all steps
func newProfile() *profile {
	c := &profile{}
	c.RegisterSteps(c.firstStep, c.secondStep, c.thirdStep)
	return c
}

func (c *profile) IsAllow(u *store.User) bool {
	c.u = u
	return (u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR)
}

func (c *profile) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
		Text:    fmt.Sprintf(lang.PROFILE_ALL_INFO, c.u.Name, c.u.FB, c.u.VK, c.u.BDate.Format(Conf.DateLayout)),
	}
	if c.u.PictureID != "" {
		replyMarkup.FileID = c.u.PictureID
	}
	replyMarkup.Buttons = append(replyMarkup.Buttons,
		lang.MARKUP_BUTTON_PROFILE_NAME,
		lang.MARKUP_BUTTON_PROFILE_FB_ACCOUNT,
		lang.MARKUP_BUTTON_PROFILE_PICTURE,
		lang.MARKUP_BUTTON_PROFILE_VK_ACCOUNT,
		lang.MARKUP_BUTTON_PROFILE_BIRTHDAY)
	return replyMarkup, nil
}

func (c *profile) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch answer {
	case lang.MARKUP_BUTTON_PROFILE_NAME:
		c.field = "name"
		replyMarkup.Text = lang.PROFILE_CURRENT_VALUE + "\n" + c.u.Name + "\n" + lang.PROFILE_WHAT_IS_YOUR_NAME
	case lang.MARKUP_BUTTON_PROFILE_FB_ACCOUNT:
		c.field = "fb"
		replyMarkup.Text = lang.PROFILE_CURRENT_VALUE + "\n" + c.u.FB + "\n" + lang.PROFILE_WHAT_IS_YOUR_FB
	case lang.MARKUP_BUTTON_PROFILE_PICTURE:
		c.field = "picture"
		replyMarkup.FileID = c.u.PictureID
		replyMarkup.Text = lang.PROFILE_WHAT_IS_YOUR_PICTURE
	case lang.MARKUP_BUTTON_PROFILE_VK_ACCOUNT:
		c.field = "vk"
		replyMarkup.Text = lang.PROFILE_CURRENT_VALUE + "\n" + c.u.VK + "\n" + lang.PROFILE_WHAT_IS_YOUR_VK
	case lang.MARKUP_BUTTON_PROFILE_BIRTHDAY:
		c.field = "bdate"
		replyMarkup.Text = lang.PROFILE_CURRENT_VALUE + "\n" + c.u.BDate.Format(Conf.DateLayout) + "\n" + lang.PROFILE_WHAT_IS_YOUR_BIRTHDAY
	}
	return replyMarkup, nil
}

func (c *profile) thirdStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	switch c.field {
	case "name":
		replyMarkup.Text = lang.PROFILE_NAME_SUCCESSFULY_UPDATED
		err = store.UpdateNameUser(c.u.ID, answer)
	case "fb":
		replyMarkup.Text = lang.PROFILE_FB_SUCCESSFULY_UPDATED
		err = store.UpdateFBUser(c.u.ID, answer)
	case "vk":
		replyMarkup.Text = lang.PROFILE_VK_SUCCESSFULY_UPDATED
		err = store.UpdateVKUser(c.u.ID, answer)
	case "picture":
		replyMarkup.Text = lang.PROFILE_PICTURE_SUCCESSFULY_UPDATED
		err = store.UpdatePictureUser(c.u.ID, answer)
	case "bdate":
		replyMarkup.Text = lang.PROFILE_BIRTHDAY_SUCCESSFULY_UPDATED
		var bdate time.Time
		bdate, err = time.Parse(Conf.DateLayout, answer)
		if err != nil {
			replyMarkup.Text = lang.ADD_EVENT_WRONG_DATE
			return replyMarkup, nil
		}
		err = store.UpdateBDateUser(c.u.ID, bdate)

	}
	if err != nil {
		return nil, err
	}
	return replyMarkup, nil
}
