package bot

import (
	"github.com/alexkarlov/15x4bot/commands"
	"github.com/alexkarlov/15x4bot/config"
	"github.com/alexkarlov/simplelog"
	"gopkg.in/telegram-bot-api.v4"
)

var Conf config.TG

// ChatType represents a tg types of chat
type ChatType string

const (
	InternalErrorText = "Внутрішня помилка, сорян"
	ButtonsCountInRow = 2

	ChatGroup      ChatType = "group"
	ChatPrivate    ChatType = "private"
	ChatChannel    ChatType = "channel"
	ChatSupergroup ChatType = "supergroup"
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
	Type     ChatType
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
		// TODO: add new record to history table

		msg := &Message{
			Type:     ChatType(update.Message.Chat.Type),
			Text:     update.Message.Text,
			Username: update.Message.From.UserName,
			ChatID:   update.Message.Chat.ID,
		}
		go b.Reply(msg)
	}
}

// SendError sends message with general error
func (b *Bot) SendError(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, InternalErrorText)
	msg.BaseChat.ReplyMarkup = markup(commands.MainMarkup)
	// TODO: process error
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Error("error while sending error message: ", err)
	}
}

// SendText sends a message to particular chat
func (b *Bot) SendText(chatID int64, msg string) error {
	replyMsg := tgbotapi.NewMessage(chatID, msg)
	_, err := b.bot.Send(replyMsg)
	return err
}

// SendTextToChannel sends a message to particular chat
func (b *Bot) SendTextToChannel(channel string, msg string) error {
	replyMsg := tgbotapi.NewMessageToChannel(channel, msg)
	_, err := b.bot.Send(replyMsg)
	return err
}

// Reply sends response (text or markup)
func (b *Bot) Reply(msg *Message) {
	c, err := lookupChat(msg)
	if err != nil {
		log.Error("error while lookup chat", err)
		b.SendError(msg.ChatID)
		return
	}
	replyMarkup, err := c.ReplyMarkup(msg)
	if err != nil {
		log.Error("error while getting reply text", err)
		b.SendError(msg.ChatID)
		return
	}
	replyMsg := tgbotapi.NewMessage(msg.ChatID, replyMarkup.Text)
	if len(replyMarkup.Buttons) > 0 {
		replyMsg.BaseChat.ReplyMarkup = markup(replyMarkup.Buttons)
	} else {
		replyMsg.BaseChat.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	b.bot.Send(replyMsg)
}

func markup(b []string) tgbotapi.ReplyKeyboardMarkup {
	rows := [][]tgbotapi.KeyboardButton{}
	c := 0
	buttons := []tgbotapi.KeyboardButton{}
	for _, bText := range b {
		if c == ButtonsCountInRow {
			rows = append(rows, buttons)
			buttons = []tgbotapi.KeyboardButton{}
			c = 0
		}
		buttons = append(buttons, tgbotapi.NewKeyboardButton(bText))
		c++
	}
	rows = append(rows, buttons)
	return tgbotapi.NewReplyKeyboard(rows...)
}
