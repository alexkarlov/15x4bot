package bot

import (
	"github.com/alexkarlov/15x4bot/config"
	"github.com/alexkarlov/simplelog"
	"gopkg.in/telegram-bot-api.v4"
)

var Conf config.TG

const (
	InternalErrorText = "Внутрішня помилка, сорян"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot() (*Bot, error) {
	tgbot, err := tgbotapi.NewBotAPI(Conf.Token)
	if err != nil {
		return nil, err
	}
	log.Infof("Authorized on account %s", tgbot.Self.UserName)
	tgbot.Debug = Conf.Debug
	bot := &Bot{
		bot: tgbot,
	}
	return bot, nil
}

type Message struct {
	Text     string
	Username string
	ChatID   int64
}

func (b *Bot) ListenUpdates() {
	u := tgbotapi.NewUpdate(Conf.UpdatesOffset)
	u.Timeout = Conf.ChatTimeout

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Error("error while starting listening updates:", err)
		return
	}

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}
		log.Infof("got new msg from [%s]: %s", update.Message.From.UserName, string(update.Message.Text))
		msg := &Message{
			Text:     update.Message.Text,
			Username: update.Message.From.UserName,
			ChatID:   update.Message.Chat.ID,
		}
		go b.Reply(msg)
	}
}

func (b *Bot) SendText(chatID int64, msg string) {
	replyMsg := tgbotapi.NewMessage(chatID, msg)
	b.bot.Send(replyMsg)
}

func (b *Bot) Reply(msg *Message) {
	c := LookupChat(msg)
	replyText, err := c.ReplayText(msg)
	if err != nil {
		log.Error("Error while getting reply text", err)
		replyText = InternalErrorText
	}
	replyMsg := tgbotapi.NewMessage(msg.ChatID, replyText)
	// button1 := tgbotapi.NewKeyboardButton("sss")
	// button2 := tgbotapi.NewKeyboardButton("aaa")
	// keyboardRow := tgbotapi.NewKeyboardButtonRow(button1, button2)
	// replyMsg.BaseChat.ReplyMarkup = tgbotapi.NewReplyKeyboard(keyboardRow)
	// replyMsg.BaseChat.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	b.bot.Send(replyMsg)
}
