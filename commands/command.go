package commands

import (
	"regexp"
	"strconv"
)

var commandPatterns = []struct {
	pattern     string
	compPattern *regexp.Regexp
	createCmd   func(cmd string) Command
}{
	{
		pattern: `addrep`,
		createCmd: func(cmd string) Command {
			return &addRepetition{}
		},
	},
	{
		pattern: `addevent`,
		createCmd: func(cmd string) Command {
			return &addEvent{}
		},
	},
	{
		pattern: `adduser`,
		createCmd: func(cmd string) Command {
			return &addUser{}
		},
	},
	{
		pattern: `addlection`,
		createCmd: func(cmd string) Command {
			return &addLection{}
		},
	},
	{
		pattern: `nextevent|next event|наступний івент|следующий ивент|когда ивент|коли івент`,
		createCmd: func(cmd string) Command {
			return &nextEvent{}
		},
	},
	{
		pattern: `nextrepetition|репетиці|репетици|repetition|коли рєпа|когда репа`,
		createCmd: func(cmd string) Command {
			return &nextRep{}
		},
	},
	{
		pattern: `documentation|документац|где дока|де дока`,
		createCmd: func(cmd string) Command {
			return &simple{
				action: "documentation",
			}
		},
	},
	{
		pattern: `about|хто ми|кто мы|про нас|15x4\?`,
		createCmd: func(cmd string) Command {
			return &simple{
				action: "about",
			}
		},
	},
	{
		pattern: `101010|3\.14|advice|порада|що робити|что делать`,
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

type Command interface {
	IsAllow(string) bool
	NextStep(answer string) (replyMsg string, err error)
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
