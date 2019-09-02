package commands

import (
	"errors"
	"github.com/alexkarlov/15x4bot/config"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"regexp"
	"strconv"
)

var (
	// Conf contains configuration for all chats
	Conf config.Chat
	// ErrWrongID happens when command can't parse id of any entity (lecture, user, etc)
	ErrWrongID = errors.New("wrong id: failed to parse id")
	// regexpID is a pattern for parsing id of any entity (lecture, user, etc)
	regexpID = regexp.MustCompile(`^[^\d]+?(\d+):`)
	// commandPatterns contains all available commands
	commandPatterns = []struct {
		pattern     string
		compPattern *regexp.Regexp
		createCmd   func(cmd string) Command
	}{
		{
			pattern: lang.MARKUP_BUTTON_CREATE_REHEARSAL,
			createCmd: func(cmd string) Command {
				return newAddRehearsal()
			},
		},
		{
			pattern: `Створити івент`,
			createCmd: func(cmd string) Command {
				return newAddEvent()
			},
		},
		{
			pattern: `Створити користувача`,
			createCmd: func(cmd string) Command {
				return newCreateUser()
			},
		},
		{
			pattern: `Змінити користувача`,
			createCmd: func(cmd string) Command {
				return newUpdateUser()
			},
		},
		{
			pattern: `Створити лекцію`,
			createCmd: func(cmd string) Command {
				return newUpsertLecture(false)
			},
		},
		{
			pattern: `Змінити лекцію`,
			createCmd: func(cmd string) Command {
				return newUpsertLecture(true)
			},
		},
		{
			pattern: `Список лекцій\(всі\)`,
			createCmd: func(cmd string) Command {
				return &lecturesList{}
			},
		},
		{
			pattern: `Список лекцій\(без опису\)`,
			createCmd: func(cmd string) Command {
				return &lecturesList{
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
				return newAddDescriptionLecture()
			},
		},
		{
			pattern: `Видалити лекцію`,
			createCmd: func(cmd string) Command {
				return newDeleteLecture()
			},
		},
		{
			pattern: `^Лекції|Івенти|Юзери|Репетиції$`,
			createCmd: func(cmd string) Command {
				reply := &markup{}
				switch cmd {
				case "Лекції":
					reply.buttons = LectureMarkup
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
			pattern: `Профіль`,
			createCmd: func(cmd string) Command {
				return newProfile()
			},
		},
		{
			pattern: `Список івентів`,
			createCmd: func(cmd string) Command {
				return &eventsList{}
			},
		},
		{
			pattern: `Змінити івент`,
			createCmd: func(cmd string) Command {
				return newUpdateEvent()
			},
		},
		{
			pattern: lang.MARKUP_BUTTON_SEND_EVENT_TO_DESIGNERS,
			createCmd: func(cmd string) Command {
				return newSendEvent(SEND_EVENT_TO_DESIGNERS)
			},
		},
		{
			pattern: lang.MARKUP_BUTTON_SEND_EVENT_TO_CHAT,
			createCmd: func(cmd string) Command {
				return newSendEvent(SEND_EVENT_TO_CHAT)
			},
		},
		{
			pattern: lang.MARKUP_BUTTON_SEND_EVENT_TO_CHANNEL,
			createCmd: func(cmd string) Command {
				return newSendEvent(SEND_EVENT_TO_CHANNEL)
			},
		},
		{
			pattern: lang.MARKUP_BUTTON_CHANGE_PHOTO_MANUAL,
			createCmd: func(cmd string) Command {
				return newUpdateArticle("event_manual_photo_id")
			},
		},
		{
			pattern: `Видалити івент`,
			createCmd: func(cmd string) Command {
				return newDeleteEvent()
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
				return newDeleteUser()
			},
		},
		{
			pattern: `Видалити репетицію`,
			createCmd: func(cmd string) Command {
				return newDeleteRehearsal()
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

// stepCallback callback which will be run when NextStep got run
type stepCallback func(string) (*ReplyMarkup, error)

// stepConstructor is a base structure for real commands
type stepConstructor struct {
	CurrentStep int
	steps       []stepCallback
}

// RegisterSteps saves steps for further call in NextStep
func (s *stepConstructor) RegisterSteps(sc ...stepCallback) {
	s.steps = sc
}

// RepeatStep decrease current step (for repeating this step in the next iteration)
func (s *stepConstructor) RepeatStep() {
	s.CurrentStep--
}

// SkipStep call next step of current command
func (s *stepConstructor) InterruptCommand() {
	s.CurrentStep = len(s.steps)
}

// SkipStep call next step of current command
func (s *stepConstructor) SkipStep(answer string) (*ReplyMarkup, error) {
	s.CurrentStep++
	defer func() {
		s.CurrentStep--
	}()
	return s.NextStep(answer)
}

// NextStep call next callback from steps according to the current step
func (s *stepConstructor) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup, err := s.steps[s.CurrentStep](answer)
	if err != nil {
		return nil, err
	}
	s.CurrentStep++
	return replyMarkup, err
}

// IsEnd determines whether the command has been finished
func (s *stepConstructor) IsEnd() bool {
	return s.CurrentStep >= len(s.steps)
}

// ReplyMarkup contains text answer of the bot and (optional) special command buttons
type ReplyMarkup struct {
	Text    string
	FileID  string
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
