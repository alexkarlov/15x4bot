package commands

import (
	"errors"
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"strings"
	"time"
)

// ErrEmptyUserTGAccount happens when admin doesn't provide tg account for a new/old user
var ErrEmptyUserTGAccount = errors.New("empty user tg account")

type upsertUser struct {
	exists   bool
	ID       int
	step     int
	name     string
	username string
	fb       string
	vk       string
	bdate    time.Time
	role     string
}

func (c *upsertUser) IsAllow(u *store.User) bool {
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *upsertUser) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	switch c.step {
	case 0:
		// if we try to insert a user
		if !c.exists {
			replyMarkup.Text = lang.USER_UPSERT_WHAT_IS_NAME
			c.step++
			break
		}
		// if we try to update a user
		users, err := store.Users(nil)
		if err != nil {
			return nil, err
		}
		for _, l := range users {
			lText := fmt.Sprintf(lang.USER_UPSERT_ITEM, l.ID, l.Name)
			replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
		}
		replyMarkup.Text = lang.CHOOSE_USER
	case 1:
		// if we try to update a user - catch user ID
		if c.exists {
			c.ID, err = parseID(answer)
			if err != nil {
				return nil, err
			}
		}
		replyMarkup.Text = lang.USER_UPSERT_WHAT_IS_NAME
	case 2:
		c.name = answer
		replyMarkup.Text = "Аккаунт в телеграмі"
		if c.exists {
			replyMarkup.Buttons = append(replyMarkup.Buttons, lang.I_DONT_KNOW)
		}
	case 3:
		if answer != lang.I_DONT_KNOW {
			c.username = strings.Trim(answer, "@")
		}
		replyMarkup.Text = "Аккаунт в Фейсбуці"
		replyMarkup.Buttons = append(replyMarkup.Buttons, lang.I_DONT_KNOW)
	case 4:
		if answer != lang.I_DONT_KNOW {
			c.fb = answer
		}
		replyMarkup.Text = "Аккаунт в ВК"
		replyMarkup.Buttons = append(replyMarkup.Buttons, lang.I_DONT_KNOW)
	case 5:
		if answer != lang.I_DONT_KNOW {
			c.vk = answer
		}
		replyMarkup.Text = "Дата народження в форматі 2006-01-02"
		replyMarkup.Buttons = append(replyMarkup.Buttons, lang.I_DONT_KNOW)
	case 6:
		if answer != lang.I_DONT_KNOW {
			t, err := time.Parse("2006-01-02", answer)
			if err != nil {
				replyMarkup.Text = "Невірний формат дати та часу. Спробуй ще!"
				return replyMarkup, nil
			}
			c.bdate = t
		}
		replyMarkup.Text = "Роль в проекті"
		roles := MessageButtons{string(store.USER_ROLE_ADMIN), string(store.USER_ROLE_LECTOR), string(store.USER_ROLE_GUEST)}
		replyMarkup.Buttons = append(replyMarkup.Buttons, roles...)
	case 7:
		role := store.NewUserRole(answer)
		var err error
		if c.exists {
			err = store.UpdateUser(c.ID, c.username, role, c.name, c.fb, c.vk, c.bdate)
			replyMarkup.Text = lang.USER_UPSERT_SUCCESSFULY_UPDATED
		} else {
			err = store.AddUserByAdmin(c.username, role, c.name, c.fb, c.vk, c.bdate)
			if err == store.ErrNoUser {
				replyMarkup.Text = lang.USER_UPSERT_USER_ALREADY_EXISTS
				return replyMarkup, nil
			}
			replyMarkup.Text = lang.USER_UPSERT_SUCCESSFULY_CREATED
		}

		if err != nil {
			return nil, err
		}
	}
	c.step++
	return replyMarkup, nil
}

func (c *upsertUser) IsEnd() bool {
	return c.step == 8
}

type usersList struct {
	u *store.User
}

func (c *usersList) IsEnd() bool {
	return true
}

func (c *usersList) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *usersList) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	list, err := store.Users(nil)
	if err != nil {
		return nil, err
	}
	for _, l := range list {
		replyMarkup.Text += fmt.Sprintf(lang.USER_UPSERT_LIST_ITEM, l.ID, l.Name, l.Role, l.Username)
	}
	return replyMarkup, nil
}

type deleteUser struct {
	step   int
	userID int
	u      *store.User
}

func (c *deleteUser) IsEnd() bool {
	return c.step == 2
}

func (c *deleteUser) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *deleteUser) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch c.step {
	case 0:
		users, err := store.Users(nil)
		if err != nil {
			return nil, err
		}
		for _, l := range users {
			lText := fmt.Sprintf(lang.USER_UPSERT_ITEM, l.ID, l.Name)
			replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
		}
		replyMarkup.Text = lang.CHOOSE_USER
	case 1:
		uID, err := parseID(answer)
		if err != nil {
			return nil, err
		}
		err = store.DeleteUser(uID)
		if err != nil {
			return nil, err
		}
		replyMarkup.Buttons = StandardMarkup(c.u.Role)
		replyMarkup.Text = lang.USER_DELETE_COMPLETE
	}
	c.step++
	return replyMarkup, nil
}
