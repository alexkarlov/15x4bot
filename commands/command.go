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
		pattern: `Створити репетицію`,
		createCmd: func(cmd string) Command {
			return &addRehearsal{}
		},
	},
	{
		pattern: `Створити івент`,
		createCmd: func(cmd string) Command {
			return &addEvent{}
		},
	},
	{
		pattern: `Створити користувача`,
		createCmd: func(cmd string) Command {
			return &addUser{}
		},
	},
	{
		pattern: `Створити лекцію`,
		createCmd: func(cmd string) Command {
			return &addLection{}
		},
	},
	{
		pattern: `Список лекцій\(всі\)`,
		createCmd: func(cmd string) Command {
			return &lectionsList{}
		},
	},
	{
		pattern: `Список лекцій\(без опису\)`,
		createCmd: func(cmd string) Command {
			return &lectionsList{
				withoutDescription: true,
			}
		},
	},
	{
		pattern: `Наступний івент`,
		createCmd: func(cmd string) Command {
			return &nextEvent{}
		},
	},
	{
		pattern: `Наступна репетиція`,
		createCmd: func(cmd string) Command {
			return &nextRep{}
		},
	},
	{
		pattern: `Документація`,
		createCmd: func(cmd string) Command {
			return &article{
				name: "documentation",
			}
		},
	},
	{
		pattern: `Хто ми`,
		createCmd: func(cmd string) Command {
			return &article{
				name: "about",
			}
		},
	},
	{
		pattern: `(?i)/start`,
		createCmd: func(cmd string) Command {
			return &article{
				name: "start",
			}
		},
	},
	{
		pattern: `(?i)/help`,
		createCmd: func(cmd string) Command {
			return &article{
				name: "help",
			}
		},
	},
	{
		pattern: `Головне меню`,
		createCmd: func(cmd string) Command {
			return &article{
				name: "main_menu",
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
		pattern: `Додати опис до лекції`,
		createCmd: func(cmd string) Command {
			return &addDescriptionLection{}
		},
	},
	{
		pattern: `Видалити лекцію`,
		createCmd: func(cmd string) Command {
			return &deleteLection{}
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
	{
		pattern: `Список івентів`,
		createCmd: func(cmd string) Command {
			return &eventsList{}
		},
	},
	{
		pattern: `Видалити івент`,
		createCmd: func(cmd string) Command {
			return &deleteEvent{}
		},
	},
	{
		pattern: `Список користувачів`,
		createCmd: func(cmd string) Command {
			return &usersList{}
		},
	},
	{
		pattern: `Видалити користувача`,
		createCmd: func(cmd string) Command {
			return &deleteUser{}
		},
	},
	{
		pattern: `Видалити репетицію`,
		createCmd: func(cmd string) Command {
			return &deleteRehearsal{}
		},
	},
	{
		pattern: `Я хочу (читати лекції|волонтерити)!`,
		createCmd: func(cmd string) Command {
			p := regexp.MustCompile(`Я хочу (читати лекції|волонтерити)`)
			m := p.FindStringSubmatch(cmd)
			role := cmd
			if len(m) > 1 {
				role = m[1]
			}
			return &messenger{
				role: role,
			}
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
	IsAllow(*store.User) bool
	NextStep(answer string) (reply *ReplyMarkup, err error)
	IsEnd() bool
}

// IsMainMenu returns if user wants to come back to main menu
func IsMainMenu(m string) bool {
	return m == "Головне меню"
}

func NewCommand(cmdName string, u *store.User) Command {
	for _, cp := range commandPatterns {
		if cp.compPattern.MatchString(cmdName) {
			cmd := cp.createCmd(cmdName)
			if cmd.IsAllow(u) {
				return cmd
			}
		}
	}
	c := &unknown{
		u: u,
	}
	return c
}
