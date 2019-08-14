package config

type DB struct {
	DSN string `env:"DB_DSN"`
}

type TG struct {
	Token         string `env:"TG_TOKEN"`
	ChatTimeout   int    `env:"TG_CHAT_TIMEOUT"`
	UpdatesOffset int    `env:"TG_UPDATES_OFFSET"`
	Debug         bool   `env:"TG_DEBUG"`
}

// Chat contains all configs for chatting with a user
type Chat struct {
	Location                string `env:"TIMEZONE" default:"Europe/Kiev"`
	TimeLayout              string `env:"TIMELAYOUT" default:"2006-01-02 15:04"`
	InternalChannelUsername string `env:"INTERNAL_CHANNEL_USERNAME"`
	InternalChatID          string `env:"INTERNAL_CHAT_ID"`
	GrammarNaziChatID       string `env:"GRAMMER_NAZI_CHAT_ID"`
	DesignerChatID          string `env:"DESIGNER_CHAT_ID"`
	RemindHourStart         int    `env:"REMIND_HOUR_START" default:"10"`
	RemindHourEnd           int    `env:"REMIND_HOUR_END" default:"21"`
}

type Config struct {
	LogLevel int `env:"LOG_LEVEL"`
	DB       DB
	TG       TG
	Chat     Chat
}
