package config

import (
	"aidanwoods.dev/go-paseto"
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"log"
)

// Config struct to store environment variables
type Config struct {
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabaseName     string `env:"DB_NAME"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`
	SslMode          string `env:"SSL_MODE"`
	Timezone         string `env:"TIMEZONE"`
}

var (
	SecretKey = paseto.NewV4AsymmetricSecretKey()
	PublicKey = SecretKey.Public()
)

// Load function to load environment variables
func Load(filenames ...string) (config Config) {
	if err := godotenv.Load(filenames...); err != nil {
		log.Fatalf("Error loading `.env` file: %v", err)
	}

	_, envErr := env.UnmarshalFromEnviron(&config)
	if envErr != nil {
		log.Fatalf("Failed to load environment variables: %v", envErr)
	}

	return config
}
