package main

import (
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/scheduler"
	"github.com/alexkarlov/simplelog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexkarlov/15x4bot/config"
	"github.com/alexkarlov/15x4bot/store"
	"github.com/antonmashko/envconf"
	_ "github.com/lib/pq"
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
	log.SetLevel(log.LogLevel(conf.LogLevel))
	// setup DB config and establish connection
	store.Conf = conf.DB
	if err := store.Init(); err != nil {
		panic(err)
	}
	// setup bot config and run bot
	bot.Conf = conf.TG
	bot, err := bot.NewBot()
	if err != nil {
		panic(err)
	}
	// start listening updates
	go bot.ListenUpdates()
	// start background job manager
	go scheduler.Run(bot)
	// wating for signals (SIGTERM - correct exit from application)
	log.Info("app has been started. waiting for signals")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	// TODO: finalise all active chats
	log.Infof("got signal %s.exiting...", sig)
}
