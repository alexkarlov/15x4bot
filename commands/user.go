package commands

import (
	"errors"
	"fmt"
	"github.com/alexkarlov/15x4bot/store"
	"regexp"
	"strconv"
	"time"
)

var ErrWrongUser = errors.New("wrong user id: failed to convert from string to int")

const (
	TEMPLATE_I_DONT_KNOW          = "Не знаю"
	TEMPLATE_USERS_LIST_ITEM      = "Юзер %d: %s, role: %s, telegram: @%s\n\n"
	TEMPLATE_USER_BUTTON          = "%d: %s"
	TEMPLATE_USER_ERROR_WRONG_ID  = "Невірно вибраний юзер"
	TEMPLATE_DELETE_USER_COMPLETE = "Юзер успішно видалений"
	TEMPLATE_CHOOSE_USER          = "Оберіть юзера"

	TEMPLATE_USER_WHAT_IS_NAME        = "Як звуть лектора/лекторку?"
	TEMPLATE_USER_SUCCESSFULY_UPDATED = "Користувач успішно змінений"
	TEMPLATE_USER_SUCCESSFULY_CREATED = "Користувач успішно створений"
	TEMPLATE_USER_IS_ALREADY_EXIST    = "Користувач з таким телеграм аккаунтом вже існує!"
)

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
	switch c.step {
	case 0:
		// if we try to insert a user
		if !c.exists {
			replyMarkup.Text = TEMPLATE_USER_WHAT_IS_NAME
			c.step++
			break
		}
		// if we try to update a user
		users, err := store.Users(nil)
		if err != nil {
			return nil, err
		}
		for _, l := range users {
			lText := fmt.Sprintf(TEMPLATE_USER_BUTTON, l.ID, l.Name)
			replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
		}
		replyMarkup.Text = TEMPLATE_CHOOSE_USER
	case 1:
		// if we try to update a user - catch user ID
		if c.exists {
			regexpLectionID := regexp.MustCompile(`^(\d+)?\:`)
			matches := regexpLectionID.FindStringSubmatch(answer)
			if len(matches) < 2 {
				replyMarkup.Text = TEMPLATE_USER_ERROR_WRONG_ID
				return replyMarkup, nil
			}
			uID, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, ErrWrongUser
			}
			c.ID = uID
		}
		replyMarkup.Text = TEMPLATE_USER_WHAT_IS_NAME
	case 2:
		c.name = answer
		replyMarkup.Text = "Аккаунт в телеграмі"
		replyMarkup.Buttons = append(replyMarkup.Buttons, TEMPLATE_I_DONT_KNOW)
	case 3:
		if answer != TEMPLATE_I_DONT_KNOW {
			c.username = answer
		}
		replyMarkup.Text = "Аккаунт в Фейсбуці"
		replyMarkup.Buttons = append(replyMarkup.Buttons, TEMPLATE_I_DONT_KNOW)
	case 4:
		if answer != TEMPLATE_I_DONT_KNOW {
			c.fb = answer
		}
		replyMarkup.Text = "Аккаунт в ВК"
		replyMarkup.Buttons = append(replyMarkup.Buttons, TEMPLATE_I_DONT_KNOW)
	case 5:
		if answer != TEMPLATE_I_DONT_KNOW {
			c.vk = answer
		}
		replyMarkup.Text = "Дата народження в форматі 2006-01-02"
		replyMarkup.Buttons = append(replyMarkup.Buttons, TEMPLATE_I_DONT_KNOW)
	case 6:
		if answer != TEMPLATE_I_DONT_KNOW {
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
		if err := store.UpsertUser(c.ID, c.username, role, c.name, c.fb, c.vk, c.bdate); err != nil {
			if err == store.ErrUserExists {
				replyMarkup.Text = TEMPLATE_USER_IS_ALREADY_EXIST
				return replyMarkup, nil
			}
			return nil, err
		}

		if c.exists {
			replyMarkup.Text = TEMPLATE_USER_SUCCESSFULY_UPDATED
		} else {
			replyMarkup.Text = TEMPLATE_USER_SUCCESSFULY_CREATED
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
		replyMarkup.Text += fmt.Sprintf(TEMPLATE_USERS_LIST_ITEM, l.ID, l.Name, l.Role, l.Username)
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
			lText := fmt.Sprintf(TEMPLATE_USER_BUTTON, l.ID, l.Name)
			replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
		}
		replyMarkup.Text = TEMPLATE_CHOOSE_USER
	case 1:
		regexpLectionID := regexp.MustCompile(`^(\d+)?\:`)
		matches := regexpLectionID.FindStringSubmatch(answer)
		if len(matches) < 2 {
			replyMarkup.Text = TEMPLATE_USER_ERROR_WRONG_ID
			return replyMarkup, nil
		}
		uID, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, ErrWrongUser
		}
		err = store.DeleteUser(uID)
		if err != nil {
			return nil, err
		}
		replyMarkup.Buttons = StandardMarkup(c.u.Role)
		replyMarkup.Text = TEMPLATE_DELETE_USER_COMPLETE
	}
	c.step++
	return replyMarkup, nil
}
