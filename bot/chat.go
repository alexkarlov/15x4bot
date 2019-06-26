package bot

import (
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
}

func (c *chat) ReplayText(m *Message) (string, error) {
	c.l.Lock()
	defer c.l.Unlock()
	if c.cmd == nil {
		c.cmd = commands.NewCommand(m.Text, m.Username)
	}

	answer, err := c.cmd.NextStep(m.Text)

	if err != nil {
		c.cmd = nil
		return "", err
	}

	if c.cmd.IsEnd() {
		log.Infof("command %#v has been finished", c.cmd)
		c.cmd = nil
	}
	return answer, err
}

func LookupChat(msg *Message) *chat {
	chatsManager.l.Lock()
	defer chatsManager.l.Unlock()
	res, ok := chatsManager.list[msg.ChatID]
	if !ok {
		log.Infof("chat with user %s not found", msg.Username)
		// TODO: lookup in the DB
		//if we haven't chatted before with this user - create a new chat
		res = &chat{
			ID:  msg.ChatID,
			cmd: commands.NewCommand(msg.Text, msg.Username),
			l:   &sync.RWMutex{},
		}
		// TODO: save in the DB
		chatsManager.list[msg.ChatID] = res
	}
	return res
}
