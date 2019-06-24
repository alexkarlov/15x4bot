package bot

import (
	"github.com/alexkarlov/15x4bot/chat"
	"github.com/alexkarlov/15x4bot/config"
	"github.com/alexkarlov/simplelog"
	"gopkg.in/telegram-bot-api.v4"
)

var Conf config.TG

type Bot struct {
	bot *tgbotapi.BotAPI
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
		t := update.Message.Text
		//do we chat with this user now?
		chat := chat.GetChat(update.Message)

		//get next step for the current chat
		replyText := chat.Speak(t)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
		//msg.ReplyToMessageID = update.Message.MessageID

		b.bot.Send(msg)
	}
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
