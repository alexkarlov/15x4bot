package commands

import (
	"github.com/alexkarlov/15x4bot/store"
	"regexp"
	"strconv"
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
		pattern: `(?i)101010|3\.14|advice|порада|що робити|что делать`,
		createCmd: func(cmd string) Command {
			return &advice{}
		},
	},
	{
		pattern: `task_\d{1,3}:`,
		createCmd: func(cmd string) Command {
			r := regexp.MustCompile(`task_(\d{1,3})?:(.*)`)
			match := r.FindStringSubmatch(cmd)
			if len(match) > 3 {
				return &unknown{}
			}
			taskID, err := strconv.Atoi(match[1])
			if err != nil {
				return &unknown{}
			}
			descr := match[2]
			return &addDescriptionLection{
				taskID:      taskID,
				description: descr,
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
	IsAllow(string) bool
	NextStep(u *store.User, answer string) (reply *ReplyMarkup, err error)
	IsEnd() bool
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
