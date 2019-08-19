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
		cc, ok := (c.cmd).(CacheCleaner)
		if ok {
			uID := cc.UserID()
			log.Infof("cleaned user cache for %d", uID)
			err = clearUserCache(uID)
			if err != nil {
				log.Infof("unsuccessful clearing cache for %d: %s", uID, err)
			}
		}
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
		ch := &chat{
			ID: msg.ChatID,
			l:  &sync.RWMutex{},
			u:  store.GuestUser(),
		}
		return ch, nil
	}

	ch, ok := chatsManager.list[msg.ChatID]
	var err error
	if !ok {
		log.Infof("chat with user %d:%s not found", msg.UserID, msg.Username)
		//if we haven't chatted before with this user - create a new chat
		ch, err = loadChat(msg)
		if err != nil {
			return nil, err
		}
		chatsManager.list[msg.ChatID] = ch
	}
	return ch, nil
}

func loadChat(msg *Message) (*chat, error) {
	ch := &chat{
		ID: msg.ChatID,
		l:  &sync.RWMutex{},
	}

	// search by tg id
	u, err := store.LoadUserByTGID(msg.UserID)
	// if there is no error - we found user, need to just return it
	if err == nil {
		ch.u = u
		return ch, nil
	}

	// unexpected error
	if err != nil && err != store.ErrNoUser {
		return nil, err
	}

	// if user doesn't have username, we need to create new record
	// otherwise we can try to find user by username (if it was created before by admin)
	// admin can create a user and set username
	// in that case we need to match created user and chatted user
	if msg.Username != "" {
		// try to search user by username
		u, err = store.LoadUserByUsername(msg.Username)
		// if we found user by tg username, we need to update tg_id in the users table for further search
		// we don't need to create a new user so we can just update users table and return a chat
		if err == nil {
			err = store.UpdateTGIDUser(u.ID, msg.UserID)
			ch.u = u
			return ch, nil
		}
		// unexpected error
		if err != nil && err != store.ErrNoUser {
			return nil, err
		}
	}

	// here we need to create a new user because we didn't find it by tg id nor tg username
	u, err = store.AddGuestUser(msg.Username, msg.UserID, msg.Name)
	if err != nil {
		return nil, err
	}
	ch.u = u
	return ch, nil
}

// CacheCleaner is an interface for cleaning user data in the cache
type CacheCleaner interface {
	UserID() int
}

// clearUserCache deletes from chatsManager all data related to the user
func clearUserCache(ID int) error {
	chatsManager.l.Lock()
	u, err := store.LoadUserByID(ID)
	if err != nil {
		return err
	}
	delete(chatsManager.list, int64(u.TGUserID))
	defer chatsManager.l.Unlock()
	return nil
}
