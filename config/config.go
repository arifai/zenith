package config

import (
	"aidanwoods.dev/go-paseto"
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"log"
)

// Config contains configuration settings loaded from environment variables.
type Config struct {
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabaseName     string `env:"DB_NAME"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`
	SslMode          string `env:"SSL_MODE"`
	Timezone         string `env:"TIMEZONE"`
	PasswordSalt     string `env:"PASSWORD_SALT"`
}

// SMTPConfig holds the configuration details required to connect to an SMTP server.
type SMTPConfig struct {
	Host     string `env:"SMTP_HOST"`
	Port     int    `env:"SMTP_PORT"`
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
}

var (
	SecretKey = paseto.NewV4AsymmetricSecretKey()
	PublicKey = SecretKey.Public()
)

// Load loads the configuration from the provided `.env` files and environment variables.
func Load(filenames ...string) (config Config) {
	if err := godotenv.Load(filenames...); err != nil {
		log.Fatalf("Error loading `.env` file: %v", err)
	}

	if _, err := env.UnmarshalFromEnviron(&config); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	return config
}
