package config

type DB struct {
	DSN string `env:"DB_DSN"`
}

type TG struct {
	Token string `env:"TG_TOKEN"`
}

type Config struct {
	DB DB
	TG TG
}
