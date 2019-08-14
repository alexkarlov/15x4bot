package commands

import (
	"encoding/json"
	"github.com/alexkarlov/15x4bot/store"
	"io/ioutil"
	"net/http"
)

type adviceResp struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Sound string `json:"sound"`
}

// advice provides fun advices from the site fucking-great-advice.ru
type advice struct {
	Resp adviceResp
	u    *store.User
}

func (c *advice) IsEnd() bool {
	return true
}

func (c *advice) IsAllow(u *store.User) bool {
	c.u = u
	return true
}

func (c *advice) NextStep(answer string) (*ReplyMarkup, error) {
	resp, err := http.Get("http://fucking-great-advice.ru/api/random")
	if err != nil {
		return nil, err
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(d, &c.Resp)
	if err != nil {
		return nil, err
	}
	replyMarkup := &ReplyMarkup{
		Text:    c.Resp.Text,
		Buttons: StandardMarkup(c.u.Role),
	}

	return replyMarkup, nil
}
