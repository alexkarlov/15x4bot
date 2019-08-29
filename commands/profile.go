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

func (c *profile) CleanUser() int {
	return c.u.ID
}

func (c *profile) IsAllow(u *store.User) bool {
	c.u = u
	return (u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR)
}

func (c *profile) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := profileMarkupButtons(c.u)
	return replyMarkup, nil
}

func profileMarkupButtons(u *store.User) *ReplyMarkup {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
		Text:    fmt.Sprintf(lang.PROFILE_ALL_INFO, u.Name, u.FB, u.VK, u.BDate.Format(Conf.DateLayout)),
	}
	if u.PictureID != "" {
		replyMarkup.FileID = u.PictureID
	}
	replyMarkup.Buttons = append(replyMarkup.Buttons,
		lang.MARKUP_BUTTON_PROFILE_NAME,
		lang.MARKUP_BUTTON_PROFILE_FB_ACCOUNT,
		lang.MARKUP_BUTTON_PROFILE_PICTURE,
		lang.MARKUP_BUTTON_PROFILE_VK_ACCOUNT,
		lang.MARKUP_BUTTON_PROFILE_BIRTHDAY)
	return replyMarkup
}

func (c *profile) secondStep(answer string) (*ReplyMarkup, error) {
	if answer == lang.MARKUP_BUTTON_PROFILE_ROLE {
		return &ReplyMarkup{
			Buttons: MainMarkup,
			Text:    lang.WRONG_PERMISSION,
		}, nil
	}
	replyMarkup, f := profileMarkupField(c.u, answer)
	c.field = f
	return replyMarkup, nil
}

func profileMarkupField(u *store.User, answer string) (*ReplyMarkup, string) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	field := ""
	switch answer {
	case lang.MARKUP_BUTTON_PROFILE_ROLE:
		field = "role"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + string(u.Role) + "\n" + lang.PROFILE_WHAT_IS_ROLE
		roles := MessageButtons{string(store.USER_ROLE_ADMIN), string(store.USER_ROLE_LECTOR), string(store.USER_ROLE_GUEST)}
		replyMarkup.Buttons = append(replyMarkup.Buttons, roles...)
	case lang.MARKUP_BUTTON_PROFILE_NAME:
		field = "name"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + u.Name + "\n" + lang.PROFILE_WHAT_IS_YOUR_NAME
	case lang.MARKUP_BUTTON_PROFILE_FB_ACCOUNT:
		field = "fb"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + u.FB + "\n" + lang.PROFILE_WHAT_IS_YOUR_FB
	case lang.MARKUP_BUTTON_PROFILE_PICTURE:
		field = "picture"
		replyMarkup.FileID = u.PictureID
		replyMarkup.Text = lang.PROFILE_WHAT_IS_YOUR_PICTURE
	case lang.MARKUP_BUTTON_PROFILE_VK_ACCOUNT:
		field = "vk"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + u.VK + "\n" + lang.PROFILE_WHAT_IS_YOUR_VK
	case lang.MARKUP_BUTTON_PROFILE_BIRTHDAY:
		field = "bdate"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + u.BDate.Format(Conf.DateLayout) + "\n" + lang.PROFILE_WHAT_IS_YOUR_BIRTHDAY
	}
	return replyMarkup, field
}

func (c *profile) thirdStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	resp, err := profileMarkupUpdateUser(c.field, c.u.ID, answer)
	if err != nil {
		return nil, err
	}
	replyMarkup.Text = resp
	return replyMarkup, nil
}

func profileMarkupUpdateUser(f string, uID int, answer string) (string, error) {
	resp := ""
	var err error
	switch f {
	case "role":
		resp = lang.PROFILE_ROLE_SUCCESSFULY_UPDATED
		err = store.UpdateRoleUser(uID, store.NewUserRole(answer))
	case "name":
		resp = lang.PROFILE_NAME_SUCCESSFULY_UPDATED
		err = store.UpdateNameUser(uID, answer)
	case "fb":
		resp = lang.PROFILE_FB_SUCCESSFULY_UPDATED
		err = store.UpdateFBUser(uID, answer)
	case "vk":
		resp = lang.PROFILE_VK_SUCCESSFULY_UPDATED
		err = store.UpdateVKUser(uID, answer)
	case "picture":
		resp = lang.PROFILE_PICTURE_SUCCESSFULY_UPDATED
		err = store.UpdatePictureUser(uID, answer)
	case "bdate":
		resp = lang.PROFILE_BIRTHDAY_SUCCESSFULY_UPDATED
		var bdate time.Time
		bdate, err = time.Parse(Conf.DateLayout, answer)
		if err != nil {
			resp = lang.ADD_EVENT_WRONG_DATE
			return resp, nil
		}
		err = store.UpdateBDateUser(uID, bdate)
	}
	return resp, err
}
