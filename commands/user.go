package commands

import (
	"errors"
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"strings"
	"time"
)

// ErrEmptyUserTGAccount happens when admin doesn't provide tg account for a new/old user
var ErrEmptyUserTGAccount = errors.New("empty user tg account")

type upsertUser struct {
	exists   bool
	ID       int
	name     string
	username string
	fb       string
	vk       string
	bdate    time.Time
	role     string
	stepConstructor
}

// newUpsertUser creates upsertUser and registers all steps
// it receives argument whether user exists or no
func newUpsertUser(e bool) *upsertUser {
	c := &upsertUser{
		exists: e,
	}
	c.RegisterSteps(c.firstStep, c.secondStep, c.thirdStep, c.fourthStep, c.fifthStep, c.sixthStep, c.seventhStep, c.eighthStep)
	return c
}

func (c *upsertUser) IsAllow(u *store.User) bool {
	return u.Role == store.USER_ROLE_ADMIN
}

// firstStep asks speaker name (if we create a new user)
// or sends users list for further chosing and manipulation
func (c *upsertUser) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	// if we try to insert a user
	if !c.exists {
		return c.SkipStep(answer)
	}
	// if we try to update a user
	users, err := store.Users(nil)
	if err != nil {
		return nil, err
	}
	for _, l := range users {
		lText := fmt.Sprintf(lang.USER_UPSERT_ITEM, l.ID, l.Name)
		replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
	}
	replyMarkup.Text = lang.CHOOSE_USER
	return replyMarkup, nil
}

// secondStep asks speaker name or parses user id from user's answer
// or asks user's name
// TODO: refactor skiping steps
func (c *upsertUser) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	// if we try to update a user - catch user ID
	if c.exists {
		c.ID, err = parseID(answer)
		if err != nil {
			return nil, err
		}
	}
	replyMarkup.Text = lang.USER_UPSERT_WHAT_IS_NAME
	return replyMarkup, nil
}

// thirdStep saves users name and asks users's tg account
func (c *upsertUser) thirdStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	c.name = answer
	replyMarkup.Text = "Аккаунт в телеграмі"
	if c.exists {
		replyMarkup.Buttons = append(replyMarkup.Buttons, lang.I_DONT_KNOW)
	}
	return replyMarkup, nil
}

// fourthStep saves username and asks user's facebook account
func (c *upsertUser) fourthStep(answer string) (*ReplyMarkup, error) {
	if answer != lang.I_DONT_KNOW {
		c.username = strings.Trim(answer, "@")
	}
	replyMarkup := &ReplyMarkup{
		Buttons: append(MainMarkup, lang.I_DONT_KNOW),
		Text:    "Аккаунт в Фейсбуці",
	}
	return replyMarkup, nil
}

// fifthStep saves fb account and asks user's vk account
func (c *upsertUser) fifthStep(answer string) (*ReplyMarkup, error) {
	if answer != lang.I_DONT_KNOW {
		c.fb = answer
	}
	replyMarkup := &ReplyMarkup{
		Text:    "Аккаунт в ВК",
		Buttons: append(MainMarkup, lang.I_DONT_KNOW),
	}
	return replyMarkup, nil
}

// sixthStep saves vk coount and asks birthday
func (c *upsertUser) sixthStep(answer string) (*ReplyMarkup, error) {
	if answer != lang.I_DONT_KNOW {
		c.vk = answer
	}
	replyMarkup := &ReplyMarkup{
		Text:    "Дата народження в форматі 2006-01-02",
		Buttons: append(MainMarkup, lang.I_DONT_KNOW),
	}
	return replyMarkup, nil
}

// seventhStep saves birthday and asks user's role
func (c *upsertUser) seventhStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	if answer != lang.I_DONT_KNOW {
		t, err := time.Parse("2006-01-02", answer)
		if err != nil {
			replyMarkup.Text = "Невірний формат дати та часу. Спробуй ще!"
			return replyMarkup, nil
		}
		c.bdate = t
	}
	replyMarkup.Text = "Роль в проекті"
	roles := MessageButtons{string(store.USER_ROLE_ADMIN), string(store.USER_ROLE_LECTOR), string(store.USER_ROLE_GUEST)}
	replyMarkup.Buttons = append(replyMarkup.Buttons, roles...)
	return replyMarkup, nil
}

// eighthStep saves user's and sends success message or fail message
func (c *upsertUser) eighthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	role := store.NewUserRole(answer)
	var err error
	if c.exists {
		err = store.UpdateUser(c.ID, c.username, role, c.name, c.fb, c.vk, c.bdate)
		replyMarkup.Text = lang.USER_UPSERT_SUCCESSFULY_UPDATED
	} else {
		err = store.AddUserByAdmin(c.username, role, c.name, c.fb, c.vk, c.bdate)
		if err == store.ErrNoUser {
			replyMarkup.Text = lang.USER_UPSERT_USER_ALREADY_EXISTS
			return replyMarkup, nil
		}
		replyMarkup.Text = lang.USER_UPSERT_SUCCESSFULY_CREATED
	}
	return replyMarkup, err
}

type usersList struct {
	u *store.User
}

func (c *usersList) IsEnd() bool {
	return true
}

func (c *usersList) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *usersList) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	list, err := store.Users(nil)
	if err != nil {
		return nil, err
	}
	for _, l := range list {
		replyMarkup.Text += fmt.Sprintf(lang.USER_UPSERT_LIST_ITEM, l.ID, l.Name, l.Role, l.Username)
	}
	return replyMarkup, nil
}

type deleteUser struct {
	step   int
	userID int
	u      *store.User
	stepConstructor
}

// newDeleteUser creates deleteUser and registers all steps
func newDeleteUser() *deleteUser {
	c := &deleteUser{}
	c.RegisterSteps(c.firstStep, c.secondStep)
	return c
}

func (c *deleteUser) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

// firstStep sends users list for further deleting
func (c *deleteUser) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	users, err := store.Users(nil)
	if err != nil {
		return nil, err
	}
	for _, l := range users {
		lText := fmt.Sprintf(lang.USER_UPSERT_ITEM, l.ID, l.Name)
		replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
	}
	replyMarkup.Text = lang.CHOOSE_USER
	return replyMarkup, nil
}

// secondStep deletes user from db
func (c *deleteUser) secondStep(answer string) (*ReplyMarkup, error) {
	uID, err := parseID(answer)
	if err != nil {
		return nil, err
	}
	err = store.DeleteUser(uID)
	if err != nil {
		return nil, err
	}
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
		Text:    lang.USER_DELETE_COMPLETE,
	}
	return replyMarkup, nil
}
