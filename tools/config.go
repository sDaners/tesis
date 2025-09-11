package tools

import (
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

var EnvPath = ".env"

var (
	once sync.Once
	conf config
)

type config struct {
	OpenAIAPIKey string `env:"OPENAI_API_KEY"`
}

func Get() config {
	once.Do(func() {
		godotenv.Load(EnvPath)

		if err := env.Parse(&conf); err != nil {
			panic(err)
		}
	})
	return conf
}
