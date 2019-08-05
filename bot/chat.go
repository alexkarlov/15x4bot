package bot

import (
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
	"sync"

	"github.com/alexkarlov/15x4bot/commands"
)

type chats struct {
	list map[int64]*chat
	l    *sync.RWMutex
}

var chatsManager = &chats{
	list: make(map[int64]*chat),
	l:    &sync.RWMutex{},
}

type chat struct {
	ID  int64
	cmd commands.Command
	l   *sync.RWMutex
	u   *store.User
}

func (c *chat) ReplyMarkup(m *Message) (*commands.ReplyMarkup, error) {
	c.l.Lock()
	defer c.l.Unlock()
	// Main menu hack
	if c.cmd == nil || commands.IsMainMenu(m.Text) {
		c.cmd = commands.NewCommand(m.Text, c.u)
	}

	log.Infof("current command %#v", c.cmd)
	answer, err := c.cmd.NextStep(m.Text)

	if err != nil {
		c.cmd = nil
		return nil, err
	}

	if c.cmd.IsEnd() {
		log.Infof("command %#v has been finished", c.cmd)
		c.cmd = nil
	}
	return answer, err
}

// lookupChat tries to find a chat by chatID in the internal list
// if it didn't find - create a new one and insert/update it
func lookupChat(msg *Message) (*chat, error) {
	chatsManager.l.Lock()
	defer chatsManager.l.Unlock()
	// check whether this is a chat message
	// chat message means that we need to consider it as guest, becuase we can't be sure that action will be safe
	if msg.Type != ChatPrivate {
		res := &chat{
			ID: msg.ChatID,
			l:  &sync.RWMutex{},
			u:  store.GuestUser(),
		}
		return res, nil
	}

	res, ok := chatsManager.list[msg.ChatID]
	if !ok {
		log.Infof("chat with user %s not found", msg.Username)
		u, err := store.LoadUser(msg.Username)
		if err != nil && err != store.ErrNoUser {
			return nil, err
		}
		//if we haven't chatted before with this user - create a new chat
		res = &chat{
			ID: msg.ChatID,
			l:  &sync.RWMutex{},
			u:  u,
		}
		// TODO: save in the DB
		if err := store.ChatUpsert(msg.ChatID, msg.Username); err != nil {
			return nil, err
		}
		if u == nil {
			res.u, err = store.LoadUser(msg.Username)
			// here we should get user; if no - something happened wrong
			if err != nil {
				return nil, err
			}
		}
		chatsManager.list[msg.ChatID] = res
	}
	return res, nil
}
