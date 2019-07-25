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
		c.cmd = commands.NewCommand(m.Text, m.Username)
	}

	log.Infof("current command %#v", c.cmd)
	answer, err := c.cmd.NextStep(c.u, m.Text)

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
	res, ok := chatsManager.list[msg.ChatID]
	if !ok {
		log.Infof("chat with user %s not found", msg.Username)
		//if we haven't chatted before with this user - create a new chat
		res = &chat{
			ID:  msg.ChatID,
			cmd: commands.NewCommand(msg.Text, msg.Username),
			l:   &sync.RWMutex{},
		}
		// TODO: save in the DB
		if err := store.ChatUpsert(msg.ChatID, msg.Username); err != nil {
			log.Error("error while chat upserting: ", err)
		}
		u, err := store.LoadUser(msg.Username)
		if err != nil {
			return nil, err
		}
		res.u = u
		chatsManager.list[msg.ChatID] = res
	}
	return res, nil
}
