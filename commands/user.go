package commands

import (
	"github.com/alexkarlov/15x4bot/store"
	"time"
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
	//TODO: move it to db
	admins := []string{"zedman95", "alex_karlov"}
	for _, admin := range admins {
		if admin == u {
			return true
		}
	}
	return false
}

func (c *addUser) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{}
	switch c.step {
	case 0:
		replyMarkup.Text = "Як звуть лектора/лекторку?"
	case 1:
		c.name = answer
		replyMarkup.Text = "Аккаунт в телеграмі"
	case 2:
		c.username = answer
		replyMarkup.Text = "Аккаунт в Фейсбуці"
	case 3:
		c.fb = answer
		replyMarkup.Text = "Аккаунт в ВК"
	case 4:
		c.vk = answer
		replyMarkup.Text = "Дата народження в форматі 2006-01-02"
	case 5:
		t, err := time.Parse("2006-01-02", answer)
		if err != nil {
			replyMarkup.Text = "Невірний формат дати та часу. Спробуй ще!"
			return replyMarkup, nil
		}
		c.bdate = t
		replyMarkup.Text = "Роль в проекті"
	case 6:
		role := store.NewUserRole(answer)
		if err := store.AddUser(c.username, role, c.name, c.fb, c.vk, c.bdate); err != nil {
			return nil, err
		}
		replyMarkup.Text = "Користувач успішно створений"
	}
	c.step++
	return replyMarkup, nil
}

func (c *addUser) IsEnd() bool {
	return c.step == 7
}
