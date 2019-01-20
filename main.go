package main

import (
	"log"
	"reminder/chats"
	"reminder/config"
	"reminder/db"

	"github.com/antonmashko/envconf"
	_ "github.com/lib/pq"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	// citations := []string{`Теория — это когда все известно, но ничего не работает. Практика — это когда все работает, но никто не знает почему. Мы же объединяем теорию и практику: ничего не работает… и никто не знает почему! \nАльберт Эйнштейн`,
	// 	`Все мы гении. Но если вы будете судить рыбу по ее способности взбираться на дерево, она проживет всю жизнь, считая себя дурой. \nАльберт Эйнштейн`,
	// 	`Если вы что-то не можете объяснить шестилетнему ребенку, вы сами этого не понимаете. \nАльберт Эйнштейн`,
	// 	`Только дурак нуждается в порядке — гений господствует над хаосом. \nАльберт Эйнштейн`,
	// 	`Есть только два способа прожить жизнь. Первый — будто чудес не существует. Второй — будто кругом одни чудеса.\nАльберт Эйнштейн`,
	// 	`Единственное, что мешает мне учиться, — это полученное мной образование \nАльберт Эйнштейн`}

	var conf config.Config
	envconf.Parse(&conf)
	bot, err := tgbotapi.NewBotAPI(conf.TG.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	//TODO: refactor it
	db.Init(conf.DB.DSN)
	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, string(update.Message.Text))
		t := update.Message.Text
		//do we chat with this user now?
		chat := chats.GetChat(update.Message)

		//get next step for the current chat
		replyText := chat.Speak(t)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
		//msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
