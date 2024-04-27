package config

import "os"

type (
	Config struct {
		DB          DB
		Server      Server
		Credentials Credentials
	}

	DB struct {
		Url string
	}

	Server struct {
		Port string
	}

	Credentials struct {
		Username string
		Password string
	}
)

func New() *Config {
	return &Config{
		DB: DB{
			Url: os.Getenv("DATABASE_URL"),
		},
		Server: Server{
			Port: os.Getenv("PORT"),
		},
		Credentials: Credentials{
			Username: os.Getenv("ADMIN_USERNAME"),
			Password: os.Getenv("ADMIN_PASSWORD"),
		},
	}
}
