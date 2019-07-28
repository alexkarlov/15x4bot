package commands

import (
	"github.com/alexkarlov/15x4bot/store"
	"regexp"
)

var commandPatterns = []struct {
	pattern     string
	compPattern *regexp.Regexp
	createCmd   func(cmd string) Command
}{
	{
		pattern: `(?i)створити репетицію`,
		createCmd: func(cmd string) Command {
			return &addRehearsal{}
		},
	},
	{
		pattern: `(?i)створити івент`,
		createCmd: func(cmd string) Command {
			return &addEvent{}
		},
	},
	{
		pattern: `(?i)створити користувача`,
		createCmd: func(cmd string) Command {
			return &addUser{}
		},
	},
	{
		pattern: `(?i)створити лекцію`,
		createCmd: func(cmd string) Command {
			return &addLection{}
		},
	},
	{
		pattern: `(?i)наступний івент`,
		createCmd: func(cmd string) Command {
			return &nextEvent{}
		},
	},
	{
		pattern: `(?i)наступна репетиція`,
		createCmd: func(cmd string) Command {
			return &nextRep{}
		},
	},
	{
		pattern: `(?i)документація`,
		createCmd: func(cmd string) Command {
			return &simple{
				action: "documentation",
			}
		},
	},
	{
		pattern: `(?i)хто ми`,
		createCmd: func(cmd string) Command {
			return &simple{
				action: "about",
			}
		},
	},
	{
		pattern: `(?i)/start`,
		createCmd: func(cmd string) Command {
			return &simple{
				action: "start",
			}
		},
	},
	{
		pattern: `(?i)/help`,
		createCmd: func(cmd string) Command {
			return &simple{
				action: "help",
			}
		},
	},
	{
		pattern: `Головне меню`,
		createCmd: func(cmd string) Command {
			return &simple{
				action: "main_menu",
			}
		},
	},
	{
		pattern: `(?i)101010|3\.14|advice|порада|що робити|что делать`,
		createCmd: func(cmd string) Command {
			return &advice{}
		},
	},
	{
		pattern: `(?i)Додати опис до лекції`,
		createCmd: func(cmd string) Command {
			return &addDescriptionLection{}
		},
	},
	{
		pattern: `^Лекції|Івенти|Юзери|Репетиції$`,
		createCmd: func(cmd string) Command {
			reply := &markup{}
			switch cmd {
			case "Лекції":
				reply.buttons = LectionMarkup
			case "Івенти":
				reply.buttons = EventMarkup
			case "Юзери":
				reply.buttons = UserMarkup
			case "Репетиції":
				reply.buttons = RehearsalMarkup
			}
			return reply
		},
	},
}

func init() {
	for i, c := range commandPatterns {
		commandPatterns[i].compPattern = regexp.MustCompile(c.pattern)
	}
}

// ReplyMarkup contains text answer of the bot and (optional) special command buttons
type ReplyMarkup struct {
	Text    string
	Buttons []string
}

type Command interface {
	IsAllow(string) bool
	NextStep(u *store.User, answer string) (reply *ReplyMarkup, err error)
	IsEnd() bool
}

// IsMainMenu returns if user wants to come back to main menu
func IsMainMenu(m string) bool {
	return m == "Головне меню"
}

func NewCommand(cmdName string, username string) Command {
	for _, cp := range commandPatterns {
		if cp.compPattern.MatchString(cmdName) {
			cmd := cp.createCmd(cmdName)
			if cmd.IsAllow(username) {
				return cmd
			}
		}
	}
	c := &unknown{}
	return c
}
