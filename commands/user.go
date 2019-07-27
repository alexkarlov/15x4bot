package commands

import (
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
	"time"
)

const (
	TEMPLATE_I_DONT_KNOW = "Не знаю"
)

type addUser struct {
	step     int
	name     string
	username string
	fb       string
	vk       string
	bdate    time.Time
	role     string
}

func (c *addUser) IsAllow(u string) bool {
	//TODO: impove filter instead of read all records
	admins, err := store.Users([]store.UserRole{store.USER_ROLE_ADMIN})
	if err != nil {
		log.Error("error while reading admins ", err)
		return false
	}
	for _, admin := range admins {
		if admin.Username == u {
			return true
		}
	}
	return false
}

func (c *addUser) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch c.step {
	case 0:
		replyMarkup.Text = "Як звуть лектора/лекторку?"
	case 1:
		c.name = answer
		replyMarkup.Text = "Аккаунт в телеграмі"
		replyMarkup.Buttons = append(replyMarkup.Buttons, TEMPLATE_I_DONT_KNOW)
	case 2:
		if answer != TEMPLATE_I_DONT_KNOW {
			c.username = answer
		}
		replyMarkup.Text = "Аккаунт в Фейсбуці"
		replyMarkup.Buttons = append(replyMarkup.Buttons, TEMPLATE_I_DONT_KNOW)
	case 3:
		if answer != TEMPLATE_I_DONT_KNOW {
			c.fb = answer
		}
		replyMarkup.Text = "Аккаунт в ВК"
		replyMarkup.Buttons = append(replyMarkup.Buttons, TEMPLATE_I_DONT_KNOW)
	case 4:
		if answer != TEMPLATE_I_DONT_KNOW {
			c.vk = answer
		}
		replyMarkup.Text = "Дата народження в форматі 2006-01-02"
		replyMarkup.Buttons = append(replyMarkup.Buttons, TEMPLATE_I_DONT_KNOW)
	case 5:
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
	case 6:
		role := store.NewUserRole(answer)
		if err := store.AddUser(c.username, role, c.name, c.fb, c.vk, c.bdate); err != nil {
			return nil, err
		}
		replyMarkup.Text = "Користувач успішно створений"
		replyMarkup.Buttons = StandardMarkup(u.Role)
	}
	c.step++
	return replyMarkup, nil
}

func (c *addUser) IsEnd() bool {
	return c.step == 7
}
