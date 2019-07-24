package commands

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
