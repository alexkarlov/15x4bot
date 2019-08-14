package commands

import (
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"math"
	"regexp"
)

const ()

type quizItem struct {
	question    string
	answer      *regexp.Regexp
	correctResp string
	wrongResp   string
}

// quiz is a hidden command for users whom knows secrets of 15x4
type quiz struct {
	correctAnswers int
	wrongAnswers   int
	step           int
	u              *store.User
	items          []*quizItem
}

func (c *quiz) IsEnd() bool {
	// we need one more step for reply with quiz evaluation
	return c.step == len(c.items)+1
}

func (c *quiz) IsAllow(u *store.User) bool {
	c.u = u
	c.items = []*quizItem{
		{
			question:    "Хочеш дізнатись наскільки глибоко ти в кролячій норі? Тоді напиши quiz me",
			answer:      regexp.MustCompile(`quiz me`),
			correctResp: "Почнемо!",
			wrongResp:   "Ок, тоді іншим разом",
		},
		{
			question:    "Як називалась лекція про горизонтальний перенесення генів?",
			answer:      regexp.MustCompile(`(?i)прокаріот.*?секс`),
			correctResp: "Ти знаєш секрети 15x4!",
			wrongResp:   "Сорян, але ні",
		},
		{
			question:    "Какой самый известный самолет на тихоокеанском театре военных действий?!",
			answer:      regexp.MustCompile(`(?i)B-29|Zero`),
			correctResp: "Молодец, братишка!",
			wrongResp:   "Нет. Это классика, это знать надо",
		},
		{
			question:    "З яким видом секса переплутала дівчина лекцію Євгена Тарасова?",
			answer:      regexp.MustCompile(`(?i)тантр`),
			correctResp: "Вірно! Шляхи до науки несповідимі!",
			wrongResp:   "Ніт!",
		},
		{
			question:    "Коли 15x4 отримали оргазм (день та місяць, наприклад: 1 січня)?",
			answer:      regexp.MustCompile(`(?i)20 жовтня`),
			correctResp: "Та ти знаєш 15x4 з давніх-давен!",
			wrongResp:   "Ніт!",
		},
		{
			question:    "З якої країни до нас завітали друзі з іншого бранчу 15x4 в 2016 році?",
			answer:      regexp.MustCompile(`(?i)молдова`),
			correctResp: "Вірно! А потім наші лектори їздили до них в гості)",
			wrongResp:   "Не вірно",
		},
		{
			question:    "Яка лекція найбільш цитована в нашому бранчі?",
			answer:      regexp.MustCompile(`(?i)слайди`),
			correctResp: "Cаме так!",
			wrongResp:   "Ні. Раджу передивитись https://www.youtube.com/watch?v=JV5MEtqtrhU",
		},
		{
			question:    "Лекція-мохіто",
			answer:      regexp.MustCompile(`(?i)лайм`),
			correctResp: "Точно!",
			wrongResp:   "Це не мохіто, це якась бодяга",
		},
		{
			question:    "Які курці мають непогані шанси на виживання?",
			answer:      regexp.MustCompile(`(?i)чорні`),
			correctResp: "Правильно. Коля радіє від щастя!",
			wrongResp:   "Не ті курці",
		},
		{
			question:    "За яку лекцію на нас наїхав Петренко?",
			answer:      regexp.MustCompile(`(?i)гомеопат`),
			correctResp: "А ти знавець холіварів!",
			wrongResp:   "Думай. Таке прізвище досить популярне в Одесі",
		},
		{
			question:    "Бородатий джентельмен що курить люльку",
			answer:      regexp.MustCompile(`(?i)Гельфанд`),
			correctResp: "In GMO we trust!",
			wrongResp:   "Ще варіанти",
		},
	}
	return true
}

// TODO: COVER THIS WITH TESTS!
func (c *quiz) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	if c.step != 0 {
		prevItem := c.items[c.step-1]
		ok, resp := prevItem.evaluate(answer)
		replyMarkup.Text = resp + "\n"
		if ok && c.step != 1 {
			c.correctAnswers++
		}
		if !ok {
			// dirty hack. TODO: fix that
			if c.step == 1 {
				c.step = 12
				return replyMarkup, nil
			}
			c.wrongAnswers++
		}
	}
	// if this is the last step
	if c.step == len(c.items) {
		c.step++
		replyMarkup.Text += c.Evaluation()
		return replyMarkup, nil
	}
	item := c.items[c.step]
	replyMarkup.Text += item.question + "\n"
	c.step++
	return replyMarkup, nil
}

func (i *quizItem) evaluate(a string) (bool, string) {
	if i.answer.MatchString(a) {
		return true, i.correctResp
	}
	return false, i.wrongResp
}

func (c *quiz) Evaluation() string {
	w := c.wrongAnswers
	if w == 0 {
		w = 1
	}
	ind := math.Round(float64(c.correctAnswers) / float64(c.wrongAnswers))
	resp := fmt.Sprintf(lang.QUIZ_RESULT, c.correctAnswers, c.wrongAnswers) + ". "
	switch true {
	case ind > 2:
		resp += lang.QUIZ_15X4_GURU
	case ind == 1:
		resp += lang.QUIZ_15X4_MID
	default:
		resp += lang.QUIZ_15X4_LOW
	}
	return resp
}
