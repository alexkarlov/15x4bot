package commands

import (
	"encoding/json"
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"strconv"
	"strings"
	"time"
)

const (
	UserRemindHour = 19
)

func nextDay(hour int) time.Time {
	curr := time.Now()
	y, m, d := curr.Date()
	loc, _ := time.LoadLocation(Conf.Location)
	rTime := time.Date(y, m, d, hour, 0, 0, 0, loc).AddDate(0, 0, 1)
	return rTime
}

type upsertLection struct {
	exists      bool
	ID          int
	u           *store.User
	step        int
	name        string
	description string
	userID      int
}

func (c *upsertLection) IsAllow(u *store.User) bool {
	c.u = u
	return (u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR)
}

func (c *upsertLection) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	switch c.step {
	case 0:
		if c.exists {
			replyMarkup, err = lections(c.u, false, false)
			break
		}
		// if the user is a lector = add him as an lection owner and skip the next step
		if c.u.Role == store.USER_ROLE_LECTOR {
			c.userID = c.u.ID
			replyMarkup.Text = lang.UPSERT_LECTURE_STEP_LECTURE_NAME
			c.step++
			break
		}
		users, err := store.Users([]store.UserRole{store.USER_ROLE_ADMIN, store.USER_ROLE_LECTOR})
		if err != nil {
			return nil, err
		}
		speakerText := ""
		for _, u := range users {
			speakerText += fmt.Sprintf(lang.UPSERT_LECTURE_STEP_SPEAKER_DETAILS, u.ID, u.Username, u.Name)
		}
		replyMarkup.Text = fmt.Sprintf(lang.UPSERT_LECTURE_STEP_SPEAKER, speakerText)
	case 1:
		replyMarkup.Text = lang.UPSERT_LECTURE_STEP_LECTURE_NAME
		if c.exists {
			lID, err := parseID(answer)
			if err != nil {
				return nil, err
			}
			c.ID = lID
			break
		}
		userID, err := strconv.Atoi(answer)
		if err != nil {
			return nil, ErrWrongID
		}
		ok, err := store.DoesUserExist(userID)
		if err != nil {
			return nil, err
		}
		if !ok {
			replyMarkup.Text = lang.WRONG_USER_ID
			return replyMarkup, nil
		}
		c.userID = userID
	case 2:
		c.name = answer
		replyMarkup.Text = lang.UPSERT_LECTURE_STEP_LECTURE_DESCRIPTION
		replyMarkup.Buttons = append(replyMarkup.Buttons, lang.I_DONT_KNOW)
	case 3:
		if answer != lang.I_DONT_KNOW {
			c.description = answer
		}
		var err error
		var lectionID int
		if c.exists {
			lectionID = c.ID
			err = store.UpdateLection(c.ID, c.name, c.description)
		} else {
			lectionID, err = store.AddLection(c.name, c.description, c.userID)
		}
		if err != nil {
			return nil, err
		}
		if c.exists {
			replyMarkup.Text = lang.UPSERT_LECTURE_SUCCESS_UPDATE_MSG
		} else {
			replyMarkup.Text = lang.UPSERT_LECTURE_SUCCESS_CREATE_MSG
		}
		if answer == lang.I_DONT_KNOW {
			execTime := nextDay(UserRemindHour)
			r := &store.RemindLection{
				ID: lectionID,
			}
			details, err := json.Marshal(r)
			if err != nil {
				return nil, err
			}
			store.AddTask(store.TASK_TYPE_REMINDER_LECTOR, execTime, string(details))
			replyMarkup.Text += "\n" + fmt.Sprintf(lang.UPSERT_LECTURE_I_WILL_REMIND, execTime.Format(Conf.TimeLayout))
		}
	}
	c.step++
	return replyMarkup, err
}

func (c *upsertLection) IsEnd() bool {
	return c.step == 4
}

type addDescriptionLection struct {
	u         *store.User
	lectionID int
	step      int
}

func (c *addDescriptionLection) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	switch c.step {
	case 0:
		replyMarkup, err = lections(c.u, true, true)
	case 1:
		lID, err := parseID(answer)
		if err != nil {
			return nil, err
		}
		l, err := store.LoadLection(lID)
		if err != nil {
			return nil, err
		}
		if c.u.Role != store.USER_ROLE_ADMIN && l.Lector.Username != c.u.Username {
			replyMarkup.Text = lang.LECTURES_ERROR_NOT_YOUR
			return replyMarkup, nil
		}
		c.lectionID = lID
		replyMarkup.Text = lang.UPSERT_LECTURE_STEP_LECTURE_DESCRIPTION
	case 2:
		err := store.AddLectionDescription(c.lectionID, answer)
		if err != nil {
			return nil, err
		}
		// send to the grammar-nazi chat
		err = sendTextToGrammarNazi(c.lectionID)
		if err != nil {
			replyMarkup.Text = lang.ADD_LECTURE_DESCRIPTION_ERROR_REMINDER_MSG
			break
		}
		replyMarkup.Text = lang.ADD_LECTURE_DESCRIPTION_COMPLETE
	}
	c.step++
	return replyMarkup, err
}

func sendTextToGrammarNazi(ID int) error {
	l, err := store.LoadLection(ID)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf(lang.ADD_LECTURE_DESCRIPTION_MSG_TO_GRAMMAR_NAZI, l.Name, l.Description)
	rh := &store.RemindChannel{
		Msg:             msg,
		ChannelUsername: Conf.GrammarNaziChatID,
	}
	details, err := json.Marshal(rh)
	if err != nil {
		return err
	}
	execTime := asSoonAsPossible()
	return store.AddTask(store.TASK_TYPE_REMINDER_TG_CHANNEL, execTime, string(details))
}

func lections(u *store.User, onlyNew bool, withoutDescription bool) (*ReplyMarkup, error) {
	reply := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	lections, err := store.Lections(onlyNew)
	if err != nil {
		return nil, err
	}
	var l []string
	for _, lection := range lections {
		if (u.Role != store.USER_ROLE_ADMIN && u.ID != lection.Lector.ID) || (withoutDescription && lection.Description != "") {
			// skip lections which doesn't belong to user (if he isn't admin) or it has description
			continue
		}
		l = append(l, fmt.Sprintf(lang.UPSERT_LECTURE_ITEM, lection.ID, lection.Name))
	}
	// if there are no appropriate lections - send special response
	if len(l) == 0 {
		reply.Text = lang.LECTURE_LIST_EMPTY
		return reply, nil
	}
	reply.Buttons = append(reply.Buttons, l...)
	reply.Text = lang.ADD_LECTURE_DESCIRPTION_CHOSE_LECTURE
	return reply, nil
}

func (c *addDescriptionLection) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR
}

func (c *addDescriptionLection) IsEnd() bool {
	return c.step == 3
}

type lectionsList struct {
	u                  *store.User
	withoutDescription bool
}

func (c *lectionsList) IsEnd() bool {
	return true
}

func (c *lectionsList) IsAllow(u *store.User) bool {
	c.u = u
	return (u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR)
}

func (c *lectionsList) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	list, err := store.Lections(true)
	if err != nil {
		return nil, err
	}
	var l []string
	for _, lection := range list {
		if c.withoutDescription && lection.Description != "" {
			// if we want to see only lections without descriptions and the current lection does have description - skip it
			continue
		}
		if c.u.Role != store.USER_ROLE_ADMIN && c.u.ID != lection.Lector.ID {
			// skip lections which doesn't belong to user (if he isn't admin)
			continue
		}
		l = append(l, fmt.Sprintf(lang.LECTURE_LIST_ITEM, lection.ID, lection.Name, lection.Description, lection.Lector.Username, lection.Lector.Name))
	}
	// if there are no appropriate lections - send special response
	if len(l) == 0 {
		replyMarkup.Text = lang.LECTURE_LIST_EMPTY
		return replyMarkup, nil
	}
	replyMarkup.Text = strings.Join(l, "\n\n")
	return replyMarkup, nil
}

type deleteLection struct {
	step      int
	lectionID int
	u         *store.User
}

func (c *deleteLection) IsEnd() bool {
	return c.step == 2
}

func (c *deleteLection) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR
}

func (c *deleteLection) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	switch c.step {
	case 0:
		replyMarkup, err = lections(c.u, false, false)
	case 1:
		lID, err := parseID(answer)
		if err != nil {
			return nil, err
		}
		l, err := store.LoadLection(lID)
		if err != nil {
			return nil, err
		}
		if c.u.Role != store.USER_ROLE_ADMIN && c.u.ID != l.Lector.ID {
			replyMarkup.Text = lang.LECTURES_ERROR_NOT_YOUR
			return replyMarkup, nil
		}
		err = store.DeleteLection(lID)
		if err != nil {
			return nil, err
		}
		replyMarkup.Buttons = StandardMarkup(c.u.Role)
		replyMarkup.Text = lang.DELETE_LECTURE_COMPLETE
	}
	c.step++
	return replyMarkup, err
}
