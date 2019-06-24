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

type Config struct {
	LogLevel int `env:"LOG_LEVEL"`
	DB       DB
	TG       TG
}
