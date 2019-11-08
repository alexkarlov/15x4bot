package config

type DB struct {
	DSN string `env:"DB_DSN" required:"true"`
}

type TG struct {
	Token         string `env:"TG_TOKEN" required:"true"`
	ChatTimeout   int    `env:"TG_CHAT_TIMEOUT" default:"60" required:"true"`
	UpdatesOffset int    `env:"TG_UPDATES_OFFSET"`
	Debug         bool   `env:"TG_DEBUG"`
	AdminAccount  string `env:"ADMIN_ACCOUNT"`
}

// Chat contains all configs for chatting with a user
type Chat struct {
	Location            string `env:"TIMEZONE" default:"Europe/Kiev" required:"true"`
	TimeLayout          string `env:"TIMELAYOUT" default:"2006-01-02 15:04" required:"true"`
	DateLayout          string `env:"DATELAYOUT" default:"2006-01-02" required:"true"`
	MainChannelUsername string `env:"MAIN_CHANNEL_USERNAME" required:"true"`
	OrgChannelUsername  string `env:"ORG_CHANNEL_USERNAME" required:"true"`
	OrgChatID           string `env:"ORG_CHAT_ID"`
	GrammarNaziChatID   string `env:"GRAMMER_NAZI_CHAT_ID"`
	DesignerChatID      string `env:"DESIGNER_CHAT_ID"`
	RemindHourStart     int    `env:"REMIND_HOUR_START" default:"10" required:"true"`
	RemindHourEnd       int    `env:"REMIND_HOUR_END" default:"21" required:"true"`
}

type Config struct {
	LogLevel int `env:"LOG_LEVEL" default:"2" required:"true"`
	DB       DB
	TG       TG
	Chat     Chat
}
