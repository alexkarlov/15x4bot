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

func (c *addUser) NextStep(answer string) (string, error) {
	replyMsg := ""
	switch c.step {
	case 0:
		replyMsg = "Як звуть лектора/лекторку?"
	case 1:
		c.name = answer
		replyMsg = "Аккаунт в телеграмі"
	case 2:
		c.username = answer
		replyMsg = "Аккаунт в Фейсбуці"
	case 3:
		c.fb = answer
		replyMsg = "Аккаунт в ВК"
	case 4:
		c.vk = answer
		replyMsg = "Дата народження в форматі 2006-01-02"
	case 5:
		t, err := time.Parse("2006-01-02", answer)
		if err != nil {
			replyMsg = "Невірний формат дати та часу. Спробуй ще!"
			return replyMsg, nil
		}
		c.bdate = t
		replyMsg = "Роль в проекті"
	case 6:
		role := store.NewUserRole(answer)
		if err := store.AddUser(c.username, role, c.name, c.fb, c.vk, c.bdate); err != nil {
			return "", err
		}
		replyMsg = "Користувач успішно створений"
	default:
		return "", ErrWrongCall
	}
	c.step++
	return replyMsg, nil
}

func (c *addUser) IsEnd() bool {
	return c.step == 7
}
