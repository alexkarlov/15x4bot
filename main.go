package main

import (
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/commands"
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
	var conf config.Config
	err := envconf.Parse(&conf)
	if err != nil {
		panic(err)
	}
	log.SetLevel(log.LogLevel(conf.LogLevel))
	log.Infof("Configs: %#v", conf)
	// setup DB config and establish connection
	store.Conf = conf.DB
	if err = store.Init(); err != nil {
		panic(err)
	}
	// setup bot config and run bot
	bot.Conf = conf.TG
	bot, err := bot.NewBot()
	if err != nil {
		panic(err)
	}
	// setup commands config
	commands.Conf = conf.Chat
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
