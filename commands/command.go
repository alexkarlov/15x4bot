package commands

import (
	"regexp"
)

var commandPatterns = []struct {
	pattern     string
	compPattern *regexp.Regexp
	createCmd   func() Command
}{
	{
		pattern: `addrep`,
		createCmd: func() Command {
			return &addRepetition{}
		},
	},
	{
		pattern: `addevent`,
		createCmd: func() Command {
			return &addEvent{}
		},
	},
	{
		pattern: `adduser`,
		createCmd: func() Command {
			return &addUser{}
		},
	},
	{
		pattern: `addlection`,
		createCmd: func() Command {
			return &addLection{}
		},
	},
	{
		pattern: `nextevent|next event|наступний івент|следующий ивент|когда ивент|коли івент`,
		createCmd: func() Command {
			return &nextEvent{}
		},
	},
	{
		pattern: `nextrepetition|репетиці|репетици|repetition|коли рєпа|когда репа`,
		createCmd: func() Command {
			return &nextRep{}
		},
	},
	{
		pattern: `documentation|документац|где дока|де дока`,
		createCmd: func() Command {
			return &simple{
				action: "documentation",
			}
		},
	},
	{
		pattern: `about|хто ми|кто мы|про нас|15x4\?`,
		createCmd: func() Command {
			return &simple{
				action: "about",
			}
		},
	},
	{
		pattern: `101010|3\.14|advice|порада|що робити|что делать`,
		createCmd: func() Command {
			return &advice{}
		},
	},
}

func init() {
	for i, c := range commandPatterns {
		commandPatterns[i].compPattern = regexp.MustCompile(c.pattern)
	}
}

type Command interface {
	IsAllow(string) bool
	NextStep(answer string) (replyMsg string, err error)
	IsEnd() bool
}

func NewCommand(cmdName string, username string) Command {
	for _, cp := range commandPatterns {
		if cp.compPattern.MatchString(cmdName) {
			cmd := cp.createCmd()
			if cmd.IsAllow(username) {
				return cmd
			}
		}
	}
	c := &unknown{}
	return c
}
