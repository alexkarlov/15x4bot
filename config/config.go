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

type Chat struct {
	Location   string `env:"TIMEZONE" default:"Europe/Kiev"`
	TimeLayout string `env:"TIMELAYOUT" default:"2006-01-02 15:04"`
}

type Config struct {
	LogLevel int `env:"LOG_LEVEL"`
	DB       DB
	TG       TG
	Chat     Chat
}
