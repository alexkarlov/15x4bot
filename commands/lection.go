package commands

import (
	"errors"
	"github.com/alexkarlov/15x4bot/store"
	"strconv"
	"strings"
)

type addLection struct {
	step        int
	name        string
	description string
	user_id     int
}

func (c *addLection) IsAllow(u string) bool {
	//TODO: move it to db
	admins := []string{"zedman95", "alex_karlov"}
	for _, admin := range admins {
		if admin == u {
			return true
		}
	}
	return false
}

func (c *addLection) NextStep(answer string) (string, error) {
	replyMsg := ""
	switch c.step {
	case 0:
		users, err := store.GetUsers([]store.UserRole{store.USER_ROLE_ADMIN, store.USER_ROLE_LECTOR})
		if err != nil {
			return "", err
		}
		replyMsg = strings.Join([]string{"Хто лектор?", strings.Join(users, "\n")}, "\n")
	case 1:
		// TODO: validate it
		userID, err := strconv.Atoi(answer)
		if err != nil {
			return "", errors.New("failed string to int converting")
		}
		c.user_id = userID
		replyMsg = "Назва лекції"
	case 2:
		c.name = answer
		replyMsg = "Опис лекції"
	case 3:
		c.description = answer
		if err := store.AddLection(c.name, c.description, c.user_id); err != nil {
			return "", err
		}
		replyMsg = "Лекцію створено"
	}
	c.step++
	return replyMsg, nil
}

func (c *addLection) IsEnd() bool {
	return c.step == 4
}
