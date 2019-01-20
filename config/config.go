package config

type DB struct {
	DSN string `env:"DB_DSN" default:"host=postgresql port=5432 user=bot sslmode=disable"`
}

type TG struct {
	Token string `env:"TG_TOKEN" default:"566944285:AAFO3UOClwS4NkFFlOFEbkxdmmd-y7VVShg"`
}

type Config struct {
	DB DB
	TG TG
}
