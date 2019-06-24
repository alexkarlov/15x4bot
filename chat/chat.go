package chat

import (
	"log"
	"sync"

	"github.com/alexkarlov/15x4bot/commands"
	"gopkg.in/telegram-bot-api.v4"
)

type chats struct {
	list map[int64]*chat
	l    sync.RWMutex
}

var chatsManager = &chats{
	list: make(map[int64]*chat),
}

type chat struct {
	ID  int64
	cmd commands.Command
}

func (c *chat) Speak(m string) (a string) {
	if c.cmd == nil {
		log.Printf("Error while speaking: empty command")
		a = "Внутрішня помилка, сорян"
		return
	}

	a, err := c.cmd.NextStep(m)

	if err != nil {
		log.Printf("Error while speaking. Source message: %v , error: %v", m, err)
		a = "Внутрішня помилка, сорян"
	}

	if c.cmd.IsEnd() {
		log.Println("FINISH!")
		c.cmd = nil
	}
	return
}

// TODO: receive interface?
func GetChat(msg *tgbotapi.Message) *chat {
	chatsManager.l.RLock()
	res, ok := chatsManager.list[msg.Chat.ID]
	chatsManager.l.RUnlock()
	if !ok {
		log.Println("Chat with user %s not found")
		//if it is new chat - create chat
		res = &chat{
			ID:  msg.Chat.ID,
			cmd: commands.NewCommand(msg.Text, msg.From),
		}
		//TODO: it's a criminal
		chatsManager.l.Lock()
		chatsManager.list[msg.Chat.ID] = res
		chatsManager.l.Unlock()
		log.Println(res.cmd)
	} else if res.cmd == nil {
		log.Println("NOT FOUND CMD :(")
		//create command
		res.cmd = commands.NewCommand(msg.Text, msg.From)
		//TODO: it's a criminal
		chatsManager.l.Lock()
		chatsManager.list[msg.Chat.ID] = res
		chatsManager.l.Unlock()
	}

	return res
}
