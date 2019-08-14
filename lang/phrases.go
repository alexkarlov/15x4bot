package lang

const (
	// ========================== EVENTS SECTION ================================================================
	EVENTS_CHOSE_EVENT = "Оберіть івент"

	// ========================== Command: addEvent (create a new event) ========================================
	ADD_EVENT_WHEN_START          = "Коли починається? Дата та час в форматі 2018-12-31 19:00"
	ADD_EVENT_WRONG_DATE          = "Невірний формат дати та часу. Наприклад, якщо івент буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00. Спробуй ще!"
	ADD_EVENT_WHEN_END            = "Коли закінчується? Дата та час в форматі 2018-12-31 19:00"
	ADD_EVENT_END_PHRASE          = "Кінець"
	ADD_EVENT_INTRO_LECTIONS_LIST = "Виберіть лекцію. Для закінчення натисніть " + ADD_EVENT_END_PHRASE
	ADD_EVENT_LECTIONS_LIST       = "Лекція %d: %s.%s"

	// ========================== Command: nextEvent (sends info about the next event) ==========================
	NEXT_EVENT           = "Де: %s, %s\nПочаток: %s\nКінець: %s"
	NEXT_EVENT_UNDEFINED = "Невідомо коли, спитай пізніше"

	// ========================== Command: eventsList (list of all events) ==========================
	EVENTS_LIST_EMPTY = "Поки івентів немає"
	EVENTS_LIST_ITEM  = "Івент %d. Початок о %s, кінець о %s, місце: %s, %s\n\n"

	// ========================== Command: deleteEvent (deleting events) ==========================
	DELETE_EVENT_COMPLETE = "Івент успішно видалено"
	DELETE_EVENT_ITEM     = "Івент %d, %s"

	// ========================== PLACES SECTION ==========================
	PLACES_LIST_BUTTONS = "Місце %d: %s\n"
	PLACES_CHOSE_PLACE  = "Оберіть місце"

	// ========================== LECTURES SECTION ==========================
	LECTURES_ERROR_NOT_YOUR = "Це не твоя лекція!"

	// ========================== Command: upsertLection (create or update lecture) ==========================
	UPSERT_LECTURE_STEP_SPEAKER_DETAILS     = "%d - %s, %s\n"
	UPSERT_LECTURE_STEP_SPEAKER             = "Хто лектор?\n%s"
	UPSERT_LECTURE_STEP_LECTURE_NAME        = "Назва лекції"
	UPSERT_LECTURE_STEP_LECTURE_DESCRIPTION = "Опис лекції. Цей опис буде відправлений в чат редакторам, тому бажано відправляти остаточний варіант"
	UPSERT_LECTURE_SUCCESS_CREATE_MSG       = "Лекцію створено"
	UPSERT_LECTURE_SUCCESS_UPDATE_MSG       = "Лекцію змінено"
	UPSERT_LECTURE_ITEM                     = "Лекція %d: %s"
	UPSERT_LECTURE_I_WILL_REMIND            = "Так як в лекції немає опису, я нагадаю про необхідність додати опис %s"

	// ========================== Command: addDescriptionLecture (add desccription for lecture) ==========================
	ADD_LECTURE_DESCIRPTION_CHOSE_LECTURE       = "Оберіть лекцію"
	ADD_LECTURE_DESCRIPTION_COMPLETE            = "Опис лекції створено"
	ADD_LECTURE_DESCRIPTION_MSG_TO_GRAMMAR_NAZI = "Лекція: %s\nОпис: %s\n"
	ADD_LECTURE_DESCRIPTION_ERROR_REMINDER_MSG  = "Опис лекції створено, але якась фігня скоїлась при створенні нагадувань в чат граммар-наці. Звернись пліз до @alex_karlov"

	// ========================== Command: lecturesList (list of all lectures) ==========================
	LECTURE_LIST_ITEM  = "Лекція %d: %s\nОпис: %s\nЛектор: @%s,  %s"
	LECTURE_LIST_EMPTY = "Поки лекцій немає"

	// ========================== Command: deleteLecture ==========================
	DELETE_LECTURE_COMPLETE = "Лекцію успішно видалено"

	// ========================= MARKUP MESSAGES (MENU BUTTONS) =========================
	// ========================= ADMIN MARKUP =========================
	MARKUP_BUTTON_LECTURES      = "Лекції"
	MARKUP_BUTTON_EVENTS        = "Івенти"
	MARKUP_BUTTON_USERS         = "Юзери"
	MARKUP_BUTTON_REHEARSALS    = "Репетиції"
	MARKUP_BUTTON_WHO_WE_ARE    = "Хто ми"
	MARKUP_BUTTON_DOCUMENTATION = "Документація"

	// ========================= SPEAKER MARKUP =========================
	MARKUP_BUTTON_NEXT_EVENT     = "Наступний івент"
	MARKUP_BUTTON_NEXT_REHEARSAL = "Наступна репетиція"

	// ========================= GUEST MARKUP =========================
	MARKUP_BUTTON_I_WANT_TO_READ_LECTURES  = "Я хочу читати лекції!"
	MARKUP_BUTTON_I_WANT_TO_BE_A_VOLUNTEER = "Я хочу волонтерити!"

	// ========================= MAIN MARKUP =========================
	MARKUP_BUTTON_MAIN_MENU = "Головне меню"

	// ========================= LECTURE MARKUP =========================
	MARKUP_BUTTON_CREATE_LECTURE                    = "Створити лекцію"
	MARKUP_BUTTON_UPDATE_LECTURE                    = "Змінити лекцію"
	MARKUP_BUTTON_ADD_DESCRIPTION                   = "Додати опис до лекції"
	MARKUP_BUTTON_LECTURES_LIST_ALL                 = "Список лекцій(всі)"
	MARKUP_BUTTON_LECTURES_LIST_WITHOUT_DESCRIPTION = "Список лекцій(без опису)"
	MARKUP_BUTTON_DELETE_LECTURE                    = "Видалити лекцію"

	// ========================= EVENT MARKUP =========================
	MARKUP_BUTTON_CREATE_EVENT = "Створити івент"
	MARKUP_BUTTON_LIST_EVENTS  = "Список івентів"
	MARKUP_BUTTON_DELETE_EVENT = "Видалити івент"

	// ========================= USER MARKUP =========================
	MARKUP_BUTTON_CREATE_USER = "Створити користувача"
	MARKUP_BUTTON_UPDATE_USER = "Змінити користувача"
	MARKUP_BUTTON_LIST_USERS  = "Список користувачів"
	MARKUP_BUTTON_DELETE_USER = "Видалити користувача"

	// ========================= REHEARSAL MARKUP =========================
	MARKUP_BUTTON_CREATE_REHEARSAL = "Створити репетицію"
	MARKUP_BUTTON_DELETE_REHEARSAL = "Видалити репетицію"

	// ========================= MESSENGER SECTION ========================================
	// ========================== Command: messenger ==========================
	MESSENGER_THANKS         = "Дякую! Я передав інформацію організаторам"
	MESSENGER_USERNAME_EMPTY = "Напиши, будь ласка, @alex_karlov ! Він розповість що робити далі)"

	// ========================== QUIZ SECTION ==========================
	// ========================== Command: quiz==========================
	QUIZ_RESULT    = "Твій результат %d вірних відповідей та %d не вірних"
	QUIZ_15X4_GURU = "Ти отримуєш знання з комосу!"
	QUIZ_15X4_MID  = "Завдяки таким, як ти, в нас є 15x4!"
	QUIZ_15X4_LOW  = "Здається, ти щойно почав своє падіння в кролячу нору, тож усе попереду!"

	// ========================== REHEARSAL SECTION ==========================
	REHEARSAL_ITEM            = "Репетиція %d: коли: %s, місце: %s"
	REHEARSAL_CHOSE_REHEARSAL = "Оберіть репетицію"
	REHEARSAL_MSG_TO_CHANNEL  = "Привіт! Нова репетиція\nДе: %s\nКоли: %s, %d %s, %s\nАдреса: %s\nМапа: %s\n"

	// ==========================  Command: addRehearsal ==========================
	ADD_REHEARSAL_WHEN               = "Коли? Дата та час в форматі 2018-12-31 19:00"
	ADD_REHEARSAL_ERROR_DATE         = "Невірний формат дати та часу. Наприклад, якщо репетиція буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00. Спробуй ще!"
	ADD_REHEARSAL_SUCCESS_MSG        = "Репетиція створена"
	ADD_REHEARSAL_ERROR_REMINDER_MSG = "Репетиція створена, але якась фігня скоїлась при створенні нагадувань. Звернись пліз до @alex_karlov"

	// ==========================  Command: nextRep ==========================
	NEXT_REHEARSAL           = "Де: %s, %s\nКоли: %s\nМапа:%s"
	NEXT_REHEARSAL_UNDEFINED = "Невідомо коли, спитай пізніше"

	// ==========================  Command: deleteRehearsal ==========================
	DELETE_REHEARSAL_COMPLETE = "Репетиція успішно видалена"

	// ========================= GENERAL MESSAGES ========================================
	WRONG_PLACE_ID     = "Невідоме місце"
	WRONG_EVENT_ID     = "Невірно вибраний івент"
	WRONG_LECTURE_ID   = "Невірно вибрана лекція"
	WRONG_USER_ID      = "Невідомий користувач"
	WRONG_REHEARSAL_ID = "Невірно вибрана репетиція"
	CHOSE_MENU         = "Оберіть пункт"
)

var (
	Weekdays = map[string]string{
		"Monday":    "Понеділок",
		"Tuesday":   "Вівторок",
		"Wednesday": "Середа",
		"Thursday":  "Четвер",
		"Friday":    "П'ятниця",
		"Saturday":  "Субота",
		"Sunday":    "Неділя",
	}
	Months = map[string]string{
		"January":   "Січень",
		"February":  "Лютий",
		"March":     "Березень",
		"April":     "Квітень",
		"May":       "Травень",
		"June":      "Червень",
		"July":      "Липень",
		"August":    "Серпень",
		"September": "Вересень",
		"October":   "Жовтень",
		"November":  "Листопад",
		"December":  "Грудень",
	}
)
