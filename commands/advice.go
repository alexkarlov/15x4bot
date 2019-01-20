package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/telegram-bot-api.v4"
)

type advice struct {
}

func (c *advice) IsEnd() bool {
	return true
}

func (c *advice) IsAllow(u *tgbotapi.User) bool {
	return true
}

func (c *advice) NextStep(answer string) (replyMsg string, err error) {
	resp, err := http.Get("http://fucking-great-advice.ru/api/random")
	if err != nil {
		return "", err
	}
	type advice struct {
		ID    int    `json:"id"`
		Text  string `json:"text"`
		Sound string `json:"sound"`
	}

	a := &advice{}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(d, &a)
	if err != nil {
		return "", err
	}
	replyMsg = a.Text

	return
}
