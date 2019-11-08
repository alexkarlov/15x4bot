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

type updateUser struct {
	u           *store.User
	updatedUser *store.User
	field       string
	stepConstructor
}

func (c *updateUser) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *updateUser) CleanUser() int {
	return c.updatedUser.ID
}

// newUpdateUser creates updateUser and registers all steps
func newUpdateUser() *updateUser {
	c := &updateUser{}
	c.RegisterSteps(c.firstStep, c.secondStep, c.thirdStep, c.fourthStep)
	return c
}

func (c *updateUser) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
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

func (c *updateUser) secondStep(answer string) (*ReplyMarkup, error) {
	ID, err := parseID(answer)
	if err != nil {
		return nil, err
	}
	u, err := store.LoadUserByID(ID)
	if err != nil {
		return nil, err
	}
	c.updatedUser = u
	replyMarkup := profileMarkupButtons(c.updatedUser)
	replyMarkup.Buttons = append(replyMarkup.Buttons, lang.MARKUP_BUTTON_PROFILE_ROLE)
	return replyMarkup, nil
}

func (c *updateUser) thirdStep(answer string) (*ReplyMarkup, error) {
	replyMarkup, f := profileMarkupField(c.updatedUser, answer)
	c.field = f
	return replyMarkup, nil
}

func (c *updateUser) fourthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	resp, err := profileMarkupUpdateUser(c.field, c.updatedUser.ID, answer)
	if err != nil {
		return nil, err
	}
	replyMarkup.Text = resp
	return replyMarkup, nil
}

type createUser struct {
	u        *store.User
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

// newCreateUser creates createUser and registers all steps
func newCreateUser() *createUser {
	c := &createUser{}
	c.RegisterSteps(c.firstStep, c.secondStep, c.thirdStep, c.fourthStep, c.fifthStep, c.sixthStep, c.seventhStep)
	return c
}

func (c *createUser) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

// firstStep asks speaker name
func (c *createUser) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	replyMarkup.Text = lang.USER_UPSERT_WHAT_IS_NAME
	return replyMarkup, nil
}

// secondStep saves users name and asks users's tg account
func (c *createUser) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	c.name = answer
	replyMarkup.Text = lang.USER_TG_ACCOUNT
	return replyMarkup, nil
}

// thirdStep saves username and asks user's facebook account
func (c *createUser) thirdStep(answer string) (*ReplyMarkup, error) {
	if answer != lang.I_DONT_KNOW {
		c.username = strings.Trim(answer, "@")
	}
	replyMarkup := &ReplyMarkup{
		Buttons: append(MainMarkup, lang.I_DONT_KNOW),
		Text:    lang.USER_FB_ACCOUNT,
	}
	return replyMarkup, nil
}

// fourthStep saves fb account and asks user's vk account
func (c *createUser) fourthStep(answer string) (*ReplyMarkup, error) {
	if answer != lang.I_DONT_KNOW {
		c.fb = answer
	}
	replyMarkup := &ReplyMarkup{
		Text:    lang.USER_VK_ACCOUNT,
		Buttons: append(MainMarkup, lang.I_DONT_KNOW),
	}
	return replyMarkup, nil
}

// fifthStep saves vk coount and asks birthday
func (c *createUser) fifthStep(answer string) (*ReplyMarkup, error) {
	if answer != lang.I_DONT_KNOW {
		c.vk = answer
	}
	replyMarkup := &ReplyMarkup{
		Text:    lang.USER_DATE_BIRTH + Conf.DateLayout,
		Buttons: append(MainMarkup, lang.I_DONT_KNOW),
	}
	return replyMarkup, nil
}

// sixthStep saves birthday and asks user's role
func (c *createUser) sixthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	if answer != lang.I_DONT_KNOW {
		t, err := time.Parse(Conf.DateLayout, answer)
		if err != nil {
			replyMarkup.Text = lang.WRONG_DATE_TIME
			return replyMarkup, nil
		}
		c.bdate = t
	}
	replyMarkup.Text = lang.USER_ROLE
	roles := MessageButtons{string(store.USER_ROLE_ADMIN), string(store.USER_ROLE_LECTOR), string(store.USER_ROLE_GUEST)}
	replyMarkup.Buttons = append(replyMarkup.Buttons, roles...)
	return replyMarkup, nil
}

// seventhStep saves user's and sends success message or fail message
func (c *createUser) seventhStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	role := store.NewUserRole(answer)
	err := store.AddUserByAdmin(c.username, role, c.name, c.fb, c.vk, c.bdate)
	if err == store.ErrNoUser {
		replyMarkup.Text = lang.USER_UPSERT_USER_ALREADY_EXISTS
		return replyMarkup, nil
	}
	replyMarkup.Text = lang.USER_UPSERT_SUCCESSFULY_CREATED
	return replyMarkup, err
}

func (c *createUser) CleanUser() int {
	return c.ID
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
		replyMarkup.Text += fmt.Sprintf(lang.USER_UPSERT_LIST_ITEM, l.ID, l.Name, l.Role, l.Username, l.FB, l.VK, l.BDate.Format(Conf.DateLayout))
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
