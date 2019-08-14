package commands

import (
	"errors"
	"github.com/alexkarlov/15x4bot/config"
	"github.com/alexkarlov/15x4bot/store"
	"regexp"
	"strconv"
)

var (
	// Conf contains configuration for all chats
	Conf config.Chat
	// ErrWrongID happens when command can't parse id of any entity (lection, user, etc)
	ErrWrongID = errors.New("wrong id: failed to parse id")
	// regexpID is a pattern for parsing id of any entity (lection, user, etc)
	regexpID = regexp.MustCompile(`^[^\d]+?(\d+):`)
	// commandPatterns contains all available commands
	commandPatterns = []struct {
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
				return &upsertUser{}
			},
		},
		{
			pattern: `Змінити користувача`,
			createCmd: func(cmd string) Command {
				return &upsertUser{
					exists: true,
				}
			},
		},
		{
			pattern: `Створити лекцію`,
			createCmd: func(cmd string) Command {
				return &upsertLection{}
			},
		},
		{
			pattern: `Змінити лекцію`,
			createCmd: func(cmd string) Command {
				return &upsertLection{
					exists: true,
				}
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
			pattern: `Я хочу читати лекції!`,
			createCmd: func(cmd string) Command {
				return &messenger{
					role: "читати лекції",
				}
			},
		},
		{
			pattern: `Я хочу волонтерити!`,
			createCmd: func(cmd string) Command {
				return &messenger{
					role: "волонтерити",
				}
			},
		},
		// hidden menu
		{
			pattern: `(?i)share your knowledge`,
			createCmd: func(cmd string) Command {
				return &quiz{}
			},
		},
	}
)

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

// Command represents a command which can used by a user
type Command interface {
	// IsAllow receives a user and performs checks whether this command is allow by the user
	IsAllow(*store.User) bool
	// NextStep receives an user answer and replies tg markup (text + buttons) and error
	// If err is not nil, that means that command failed and need to send a general error text for the user
	NextStep(answer string) (reply *ReplyMarkup, err error)
	// IsEnd determines whether the command has been finished (last iteration)
	IsEnd() bool
}

// IsMainMenu returns if user wants to come back to main menu
func IsMainMenu(m string) bool {
	return m == "Головне меню"
}

// NewCommand creates a new command by user anwser
// It checks whether requested command is allow for the user
// If there is no appropriate command - it returns "unknown" command
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

// parseID parses id from standard user answers
func parseID(a string) (int, error) {
	matches := regexpID.FindStringSubmatch(a)
	if len(matches) < 2 {
		return 0, ErrWrongID
	}
	return strconv.Atoi(matches[1])
}
