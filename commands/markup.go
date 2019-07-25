package commands

import (
	"github.com/alexkarlov/15x4bot/store"
)

type MessageButtons []string

var (
	AdminMarkup = MessageButtons{
		"Створити івент",
		"Створити репетицію",
		"Створити користувача",
		"Створити лекцію",
	}

	SpeakerMarkup = MessageButtons{
		"Додати опис до лекції",
	}

	GuestMarkup = MessageButtons{
		"Наступний івент",
		"Наступна репетиція",
		"Хто ми",
		"Документація",
	}
)

// StandardMarkup returns general markup depends on provided role
func StandardMarkup(role store.UserRole) MessageButtons {
	buttons := MessageButtons(GuestMarkup)
	if role == store.USER_ROLE_ADMIN {
		buttons = append(buttons, AdminMarkup...)
		buttons = append(buttons, SpeakerMarkup...)
	} else if role == store.USER_ROLE_LECTOR {
		buttons = append(buttons, SpeakerMarkup...)
	}
	return buttons
}
