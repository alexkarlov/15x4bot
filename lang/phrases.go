package lang

const (
	// ========================== EVENTS SECTION ================================================================
	EVENTS_CHOSE_EVENT = "Оберіть івент"

	// ========================== Command: addEvent (create a new event) ========================================
	ADD_EVENT_WHEN_START                   = "Коли починається? Дата та час в форматі 2018-12-31 19:00"
	ADD_EVENT_WRONG_DATETIME               = "Невірний формат дати та часу. Наприклад, якщо івент буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00. Спробуй ще!"
	ADD_EVENT_WRONG_DATE                   = "Невірний формат дати. Приклад: 2018-12-20"
	ADD_EVENT_WHEN_END                     = "Коли закінчується? Дата та час в форматі 2018-12-31 19:00"
	ADD_EVENT_END_PHRASE                   = "Кінець"
	ADD_EVENT_INTRO_LECTURES_LIST          = "Виберіть лекцію. Для закінчення натисніть " + ADD_EVENT_END_PHRASE
	ADD_EVENT_LECTURES_LIST                = "Лекція %d: %s.%s"
	ADD_EVENT_TEXT_EVENT                   = "Текст івенту"
	ADD_EVENT_SUCCESS_MSG                  = "Івент створено"
	ADD_EVENT_SEND_EVENT_TO_DESIGNERS_CHAT = `В нас новий івент)
	Коли: %s - %s
	Де: %s
	Лекції: 
	`
	ADD_EVENT_EMPTY_LECTURE_DESCRIPTION = "В лекції %d:%s відсутній опис, івент не може бути відправлений дизайнерам"
	ADD_EVENT_EMPTY_PICTURE_SPEAKER     = "В лектора %d:%s відсутнє фото, івент не може бути відправлений дизайнерам"
	ADD_EVENT_EMPTY_FB                  = "Поле \"FB івент\" пусте, івент не може бути відправлений в канал"
	ADD_EVENT_EMPTY_POSTER              = "Поле \"Постер\" пусте, івент не може бути відправлений в канал"
	ADD_EVENT_SEND_EVENT_TO_CHANNEL     = `Новий івент! 
	Коли: %s - %s
	Де: %s
	`
	ADD_EVENT_SEND_EVENT_TO_COMMON_CHAT = `Новий івент!
	Коли: %s
	Де: %s
	ФБ івент: %s
	Ставте вподобайки, поширюйте, запрошуйте друзів!
	Як і шо робити:
	`

	// ========================== Command: nextEvent (sends info about the next event) ===========================
	NEXT_EVENT           = "Де: %s, %s\nПочаток: %s\nКінець: %s"
	NEXT_EVENT_UNDEFINED = "Невідомо коли, спитай пізніше"

	// ========================== Command: eventsList (list of all events) =======================================
	EVENTS_LIST_EMPTY = "Поки івентів немає"
	EVENTS_LIST_ITEM  = "Івент %d. Початок о %s, кінець о %s, місце: %s, %s\n\n"

	// ========================== Command: deleteEvent (deleting events) =========================================
	DELETE_EVENT_COMPLETE = "Івент успішно видалено"
	DELETE_EVENT_ITEM     = "Івент %d, %s"

	// ========================== PLACES SECTION =================================================================
	PLACES_LIST_BUTTONS = "Місце %d: %s\n"
	PLACES_CHOSE_PLACE  = "Оберіть місце"

	// ========================== LECTURES SECTION ===============================================================
	LECTURES_ERROR_NOT_YOUR = "Це не твоя лекція!"

	// ========================== Command: upsertLecture (create or update lecture) ==============================
	UPSERT_LECTURE_STEP_SPEAKER_DETAILS     = "%d - %s, %s\n"
	UPSERT_LECTURE_STEP_SPEAKER             = "Хто лектор?\n%s"
	UPSERT_LECTURE_STEP_LECTURE_NAME        = "Назва лекції"
	UPSERT_LECTURE_STEP_LECTURE_DESCRIPTION = "Опис лекції"
	UPSERT_LECTURE_SUCCESS_CREATE_MSG       = "Лекцію створено"
	UPSERT_LECTURE_SUCCESS_UPDATE_MSG       = "Лекцію змінено"
	UPSERT_LECTURE_SEND_TO_GRAMMAR_NAZI     = "Відправити лекцію на перевірку нашим редакторам?"
	UPSERT_LECTURE_ITEM                     = "Лекція %d: %s"
	UPSERT_LECTURE_I_WILL_REMIND            = "Так як в лекції немає опису, я нагадаю про необхідність додати опис %s"

	// ========================== Command: addDescriptionLecture (add desccription for lecture) ===================
	ADD_LECTURE_DESCIRPTION_CHOSE_LECTURE       = "Оберіть лекцію"
	ADD_LECTURE_DESCRIPTION_COMPLETE            = "Опис лекції створено"
	ADD_LECTURE_DESCRIPTION_MSG_TO_GRAMMAR_NAZI = "Лекція: %s\nОпис: %s\n"
	ADD_LECTURE_DESCRIPTION_ERROR_REMINDER_MSG  = "Якась фігня скоїлась при створенні нагадувань в чат граммар-наці. Звернись пліз до @alex_karlov"
	ADD_LECTURE_DESCRIPTION_REMINDER_MSG_OK     = "Опис буде відправлено в чат граммар-наці"

	// ========================== Command: lecturesList (list of all lectures) ====================================
	LECTURE_LIST_ITEM  = "Лекція %d: %s\nОпис: %s\nЛектор: @%s,  %s"
	LECTURE_LIST_EMPTY = "Поки лекцій немає"

	// ========================== Command: deleteLecture ==========================================================
	DELETE_LECTURE_COMPLETE = "Лекцію успішно видалено"

	// ========================= MARKUP MESSAGES (MENU BUTTONS) ===================================================
	// ========================= ADMIN MARKUP =====================================================================
	MARKUP_BUTTON_LECTURES   = "Лекції"
	MARKUP_BUTTON_EVENTS     = "Івенти"
	MARKUP_BUTTON_USERS      = "Юзери"
	MARKUP_BUTTON_REHEARSALS = "Репетиції"

	// ========================= SPEAKER MARKUP ===================================================================
	MARKUP_BUTTON_NEXT_EVENT     = "Наступний івент"
	MARKUP_BUTTON_NEXT_REHEARSAL = "Наступна репетиція"

	// ========================= GUEST MARKUP =====================================================================
	MARKUP_BUTTON_I_WANT_TO_READ_LECTURES  = "Я хочу читати лекції!"
	MARKUP_BUTTON_I_WANT_TO_BE_A_VOLUNTEER = "Я хочу волонтерити!"

	// ========================= MAIN MARKUP ======================================================================
	MARKUP_BUTTON_MAIN_MENU     = "Головне меню"
	MARKUP_BUTTON_PROFILE       = "Профіль"
	MARKUP_BUTTON_WHO_WE_ARE    = "Хто ми"
	MARKUP_BUTTON_DOCUMENTATION = "Документація"

	// ========================= PROFILE MARKUP ===================================================================
	MARKUP_BUTTON_PROFILE_NAME       = "Змінити ім'я"
	MARKUP_BUTTON_PROFILE_BIRTHDAY   = "Змінити дату народження"
	MARKUP_BUTTON_PROFILE_VK_ACCOUNT = "Змінити VK акаунт"
	MARKUP_BUTTON_PROFILE_FB_ACCOUNT = "Змінити FB акаунт"
	MARKUP_BUTTON_PROFILE_PICTURE    = "Змінити фото"
	MARKUP_BUTTON_PROFILE_ROLE       = "Змінити роль"

	// ========================= LECTURE MARKUP ===================================================================
	MARKUP_BUTTON_CREATE_LECTURE                    = "Створити лекцію"
	MARKUP_BUTTON_UPDATE_LECTURE                    = "Змінити лекцію"
	MARKUP_BUTTON_ADD_DESCRIPTION                   = "Додати опис до лекції"
	MARKUP_BUTTON_LECTURES_LIST_ALL                 = "Список лекцій(всі)"
	MARKUP_BUTTON_LECTURES_LIST_WITHOUT_DESCRIPTION = "Список лекцій(без опису)"
	MARKUP_BUTTON_DELETE_LECTURE                    = "Видалити лекцію"

	// ========================= EVENT MARKUP =====================================================================
	MARKUP_BUTTON_CREATE_EVENT            = "Створити івент"
	MARKUP_BUTTON_LIST_EVENTS             = "Список івентів"
	MARKUP_BUTTON_DELETE_EVENT            = "Видалити івент"
	MARKUP_BUTTON_SEND_EVENT_TO_DESIGNERS = "Відправити івент в чат дизайнерів"
	MARKUP_BUTTON_SEND_EVENT_TO_CHANNEL   = "Відправити івент в канал"
	MARKUP_BUTTON_SEND_EVENT_TO_CHAT      = "Відправити івент в загальний чат"

	// ========================= USER MARKUP ======================================================================
	MARKUP_BUTTON_CREATE_USER = "Створити користувача"
	MARKUP_BUTTON_UPDATE_USER = "Змінити користувача"
	MARKUP_BUTTON_LIST_USERS  = "Список користувачів"
	MARKUP_BUTTON_DELETE_USER = "Видалити користувача"

	// ========================= REHEARSAL MARKUP =================================================================
	MARKUP_BUTTON_CREATE_REHEARSAL       = "Створити репетицію"
	MARKUP_BUTTON_DELETE_REHEARSAL       = "Видалити репетицію"
	MARKUP_BUTTON_NOTIF_REHEARSAL_NOW    = "Відправити зараз"
	MARKUP_BUTTON_NOTIF_BEFORE_REHEARSAL = "Відправити за день до репетиції"

	// ========================= MESSENGER SECTION ================================================================
	// ========================== Command: messenger ==============================================================
	MESSENGER_THANKS         = "Дякую! Я передав інформацію організаторам"
	MESSENGER_USERNAME_EMPTY = "Нажаль, в тебе не встановлений телеграм username. Напиши, будь ласка, @alex_karlov ! Він розповість що робити далі)"

	// ========================== QUIZ SECTION ====================================================================
	// ========================== Command: quiz====================================================================
	QUIZ_RESULT    = "Твій результат %d вірних відповідей та %d не вірних"
	QUIZ_15X4_GURU = "Ти отримуєш знання з комосу!"
	QUIZ_15X4_MID  = "Завдяки таким, як ти, в нас є 15x4!"
	QUIZ_15X4_LOW  = "Здається, ти щойно почав своє падіння в кролячу нору, тож усе попереду!"

	// ========================== REHEARSAL SECTION ===============================================================
	REHEARSAL_ITEM            = "Репетиція %d: коли: %s, місце: %s"
	REHEARSAL_CHOSE_REHEARSAL = "Оберіть репетицію"
	REHEARSAL_MSG_TO_CHANNEL  = "Привіт! Нова репетиція\nДе: %s\nКоли: %s, %d %s, %s\nАдреса: %s\nМапа: %s\n"

	// ==========================  Command: addRehearsal ==========================================================
	ADD_REHEARSAL_WHEN               = "Коли? Дата та час в форматі 2018-12-31 19:00"
	ADD_REHEARSAL_ERROR_DATE         = "Невірний формат дати та часу. Наприклад, якщо репетиція буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00. Спробуй ще!"
	ADD_REHEARSAL_SUCCESS_MSG        = "Репетиція створена. Відправити повідомлення в чат та канал зараз чи за день до репетиції?"
	ADD_REHEARSAL_ERROR_REMINDER_MSG = "Якась фігня скоїлась при створенні нагадувань. Звернись пліз до @alex_karlov"
	ADD_REHEARSAL_REMINDER_OK        = "Повідомлення про репетицію буде відправлене %s"

	// ==========================  Command: nextRep ===============================================================
	NEXT_REHEARSAL           = "Де: %s, %s\nКоли: %s\nМапа:%s"
	NEXT_REHEARSAL_UNDEFINED = "Невідомо коли, спитай пізніше"

	// ==========================  Command: deleteRehearsal =======================================================
	DELETE_REHEARSAL_COMPLETE = "Репетиція успішно видалена"

	// ========================== USER SECTION ====================================================================
	// ========================== Command: upsertUser =============================================================
	USER_UPSERT_LIST_ITEM           = "Юзер %d\nІм'я: %s\nРоль: %s\nTelegram: @%s\nFB: %s\nVK: %s\nДата народження: %s\n \n\n"
	USER_UPSERT_ITEM                = "Юзер %d: %s"
	USER_UPSERT_WHAT_IS_NAME        = "Як звуть лектора/лекторку?"
	USER_UPSERT_SUCCESSFULY_UPDATED = "Користувач успішно змінений"
	USER_UPSERT_SUCCESSFULY_CREATED = "Користувач успішно створений"
	USER_UPSERT_USER_ALREADY_EXISTS = "Користувач з таким телеграм аккаунтом вже існує! Якщо хочеш змінити дані юзера - вибери змінити юзера з меню Юзери"

	// ========================== Command: upsertUser =============================================================
	USER_DELETE_COMPLETE = "Юзер успішно видалений"

	// ========================= GENERAL MESSAGES =================================================================
	WRONG_PLACE_ID     = "Невідоме місце"
	WRONG_EVENT_ID     = "Невірно вибраний івент"
	WRONG_LECTURE_ID   = "Невірно вибрана лекція"
	WRONG_USER_ID      = "Невідомий користувач"
	WRONG_REHEARSAL_ID = "Невірно вибрана репетиція"
	CHOSE_MENU         = "Оберіть пункт"
	I_DONT_KNOW        = "Не знаю"
	CHOOSE_USER        = "Оберіть юзера"
	MARKUP_BUTTON_YES  = "Так"
	MARKUP_BUTTON_NO   = "Ні"
	WRONG_PERMISSION   = "Воу-воу, палєгче братиш"
	DONE               = "Зроблено!"

	// ========================= PROFILE SECTION ==================================================================
	// ========================= Command: profileName =============================================================
	PROFILE_ALL_INFO                     = "Профіль:\nІм'я: %s\nFB: %s\nVK: %s\nДата народження: %s\n"
	PROFILE_CURRENT_VALUE                = "Поточне значення:"
	PROFILE_WHAT_IS_YOUR_NAME            = "Якщо бажаєш змінити, введи ім'я та прізвище"
	PROFILE_WHAT_IS_ROLE                 = "Якщо бажаєш змінити, вибери роль"
	PROFILE_ROLE_SUCCESSFULY_UPDATED     = "Роль успішно змінено"
	PROFILE_NAME_SUCCESSFULY_UPDATED     = "Ім'я та прізвище успішно змінено"
	PROFILE_WHAT_IS_YOUR_FB              = "Якщо бажаєш змінити, введи fb акаунт"
	PROFILE_FB_SUCCESSFULY_UPDATED       = "fb акаунт успішно змінено"
	PROFILE_WHAT_IS_YOUR_VK              = "Якщо бажаєш змінити, введи vk акаунт"
	PROFILE_VK_SUCCESSFULY_UPDATED       = "vk успішно змінено"
	PROFILE_WHAT_IS_YOUR_BIRTHDAY        = "Якщо бажаєш змінити, введи дату народження в форматі 2018-12-31"
	PROFILE_BIRTHDAY_SUCCESSFULY_UPDATED = "Дату народження успішно змінено"
	PROFILE_WHAT_IS_YOUR_PICTURE         = "Якщо бажаєш змінити, завантаж фото"
	PROFILE_PICTURE_SUCCESSFULY_UPDATED  = "Фото успішно завантажено"
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
	UnknownMsgs = []string{"Вибач, я не розумію тебе", "Ніпанятна", "Шта?"}
)
