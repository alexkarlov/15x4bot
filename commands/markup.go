package commands

import (
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
)

// MessageButtons needs for me nu buttons
type MessageButtons []string

var (
	// AdminMarkup contains buttons for admins
	AdminMarkup = MessageButtons{
		lang.MARKUP_BUTTON_PROFILE,
		lang.MARKUP_BUTTON_LECTURES,
		lang.MARKUP_BUTTON_EVENTS,
		lang.MARKUP_BUTTON_USERS,
		lang.MARKUP_BUTTON_REHEARSALS,
		lang.MARKUP_BUTTON_WHO_WE_ARE,
		lang.MARKUP_BUTTON_DOCUMENTATION,
	}

	// SpeakerMarkup contains buttons for speakers
	SpeakerMarkup = MessageButtons{
		lang.MARKUP_BUTTON_PROFILE,
		lang.MARKUP_BUTTON_LECTURES,
		lang.MARKUP_BUTTON_NEXT_EVENT,
		lang.MARKUP_BUTTON_NEXT_REHEARSAL,
		lang.MARKUP_BUTTON_DOCUMENTATION,
	}

	// GuestMarkup contains buttons for guests (users without any special permissions)
	GuestMarkup = MessageButtons{
		lang.MARKUP_BUTTON_I_WANT_TO_READ_LECTURES,
		lang.MARKUP_BUTTON_I_WANT_TO_BE_A_VOLUNTEER,
		lang.MARKUP_BUTTON_NEXT_EVENT,
		lang.MARKUP_BUTTON_NEXT_REHEARSAL,
		lang.MARKUP_BUTTON_WHO_WE_ARE,
	}

	// MainMarkup contains buttons for come back to main manu
	MainMarkup = MessageButtons{
		lang.MARKUP_BUTTON_MAIN_MENU,
	}

	// LectureMarkup contains buttons for admins (lectures submenu)
	LectureMarkup = MessageButtons{
		lang.MARKUP_BUTTON_CREATE_LECTURE,
		lang.MARKUP_BUTTON_UPDATE_LECTURE,
		lang.MARKUP_BUTTON_ADD_DESCRIPTION,
		lang.MARKUP_BUTTON_LECTURES_LIST_ALL,
		lang.MARKUP_BUTTON_LECTURES_LIST_WITHOUT_DESCRIPTION,
		lang.MARKUP_BUTTON_DELETE_LECTURE,
	}

	// EventMarkup contains buttons for admins (events submenu)
	EventMarkup = MessageButtons{
		lang.MARKUP_BUTTON_NEXT_EVENT,
		lang.MARKUP_BUTTON_CREATE_EVENT,
		lang.MARKUP_BUTTON_LIST_EVENTS,
		lang.MARKUP_BUTTON_DELETE_EVENT,
	}

	// UserMarkup contains buttons for admins (users submenu)
	UserMarkup = MessageButtons{
		lang.MARKUP_BUTTON_CREATE_USER,
		lang.MARKUP_BUTTON_UPDATE_USER,
		lang.MARKUP_BUTTON_LIST_USERS,
		lang.MARKUP_BUTTON_DELETE_USER,
	}

	// RehearsalMarkup contains buttons for admins (rehearsals submenu)
	RehearsalMarkup = MessageButtons{
		lang.MARKUP_BUTTON_NEXT_REHEARSAL,
		lang.MARKUP_BUTTON_CREATE_REHEARSAL,
		lang.MARKUP_BUTTON_DELETE_REHEARSAL,
	}

	// ProfileMarkup contains buttons for speakers and admins
	ProfileMarkup = MessageButtons{
		lang.MARKUP_BUTTON_PROFILE_NAME,
		lang.MARKUP_BUTTON_PROFILE_BIRTHDAY,
		lang.MARKUP_BUTTON_PROFILE_FB_ACCOUNT,
		lang.MARKUP_BUTTON_PROFILE_VK_ACCOUNT,
		lang.MARKUP_BUTTON_PROFILE_PICTURE,
	}
)

// StandardMarkup returns general markup depends on provided role
func StandardMarkup(role store.UserRole) MessageButtons {
	buttons := MessageButtons(GuestMarkup)
	if role == store.USER_ROLE_ADMIN {
		buttons = MessageButtons(AdminMarkup)
	} else if role == store.USER_ROLE_LECTOR {
		buttons = MessageButtons(SpeakerMarkup)
	}
	buttons = append(buttons, MainMarkup...)
	return buttons
}

// markup is a simple command for send menu buttons
type markup struct {
	buttons MessageButtons
}

func (c *markup) IsEnd() bool {
	return true
}

func (c *markup) IsAllow(u *store.User) bool {
	return true
}

func (c *markup) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: c.buttons,
		Text:    lang.CHOSE_MENU,
	}
	replyMarkup.Buttons = append(replyMarkup.Buttons, MainMarkup...)
	return replyMarkup, nil
}
