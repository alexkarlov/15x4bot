package commands

import (
	"github.com/alexkarlov/15x4bot/store"
)

type msgProxy struct {
	cmd string
}

func (c *msgProxy) IsEnd() bool {
	return true
}

func (c *msgProxy) IsAllow(u *store.User) bool {
	return true
}

func (c *msgProxy) NextStep(answer string) (*ReplyMarkup, error) {
	reply := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	// admins, err := store.Users(store.USER_ROLE_ADMIN)
	// if err != nil {
	// 	return nil, err
	// }
	// if len(admins) == 0 {
	// 	return reply, nil
	// }
	// u := admins[0]
	// u.TGChat()
	// // send notification to the admins
	return reply, nil
}
