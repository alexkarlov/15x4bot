package commands

import (
	"github.com/alexkarlov/15x4bot/store"
)

type MessageButtons []string

const (
	TEMPLATE_CHOSE_MENU = "Оберіть пункт"
)

var (
	AdminMarkup = MessageButtons{
		"Лекції",
		"Івенти",
		"Юзери",
		"Репетиції",
		"Хто ми",
		"Документація",
	}

	SpeakerMarkup = MessageButtons{
		"Лекції",
		"Наступний івент",
		"Наступна репетиція",
		"Хто ми",
	}

	GuestMarkup = MessageButtons{
		"Я хочу читати лекції!",
		"Я хочу волонтерити!",
		"Наступний івент",
		"Наступна репетиція",
		"Хто ми",
	}

	MainMarkup = MessageButtons{
		"Головне меню",
	}

	LectionMarkup = MessageButtons{
		"Створити лекцію",
		"Додати опис до лекції",
		"Список лекцій(всі)",
		"Список лекцій(без опису)",
		"Видалити лекцію",
	}
	EventMarkup = MessageButtons{
		"Наступний івент",
		"Створити івент",
		"Список івентів",
		"Видалити івент",
	}
	UserMarkup = MessageButtons{
		"Створити користувача",
		"Список користувачів",
		"Видалити користувача",
	}
	RehearsalMarkup = MessageButtons{
		"Наступна репетиція",
		"Створити репетицію",
		"Видалити репетицію",
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
		Text:    TEMPLATE_CHOSE_MENU,
	}
	replyMarkup.Buttons = append(replyMarkup.Buttons, MainMarkup...)
	return replyMarkup, nil
}
