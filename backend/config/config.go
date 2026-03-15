package config

import "os"

type Config struct {
	OpenAIKey   string
	DatabaseURL string
}

func Load() Config {
	return Config{
		OpenAIKey:   os.Getenv("OPENAI_API_KEY"),
		DatabaseURL: "postgres://postgres:postgres@host.docker.internal:5434/gorag?sslmode=disable",
	}
}
