package commands

import (
	"regexp"

	"gopkg.in/telegram-bot-api.v4"
)

var commandPatterns = []struct {
	pattern     string
	compPattern *regexp.Regexp
	cmd         Command
}{
	{
		pattern: `addrep`,
		cmd:     &addRepetition{},
	},
	{
		pattern: `addevent`,
		cmd:     &addRepetition{},
	},
	{
		pattern: `nextevent|next event|наступний івент|следующий ивент|когда ивент|коли івент`,
		cmd:     &nextEvent{},
	},
	{
		pattern: `nextrepetition|репетиці|репетици|repetition|коли рєпа|когда репа`,
		cmd:     &nextRep{},
	},
	{
		pattern: `documentation|документац|где дока|де дока`,
		cmd: &simple{
			action: "documentation",
		},
	},
	{
		pattern: `about|хто ми|кто мы|про нас|15x4\?`,
		cmd: &simple{
			action: "about",
		},
	},
	{
		pattern: `101010|3\.14|advice|порада|що робити|что делать`,
		cmd:     &advice{},
	},
}

func init() {
	for i, c := range commandPatterns {
		commandPatterns[i].compPattern = regexp.MustCompile(c.pattern)
	}
}

type Command interface {
	IsAllow(*tgbotapi.User) bool
	NextStep(answer string) (replyMsg string, err error)
	IsEnd() bool
}

func NewCommand(cmdName string, user *tgbotapi.User) Command {
	for _, cp := range commandPatterns {
		if cp.compPattern.MatchString(cmdName) && cp.cmd.IsAllow(user) {
			return cp.cmd
		}
	}
	c := &unknown{}
	return c
}
