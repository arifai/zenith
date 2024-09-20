package config

import (
	"aidanwoods.dev/go-paseto"
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"log"
)

type (
	// Env is an interface for `env` file
	Env interface {
		// LoadDefault loads the default configuration settings from the provided environment files.
		// These settings include database connection parameters such as DB_HOST, DB_USER, DB_PASSWORD, etc.
		// `filenames` is a variadic parameter allowing multiple filenames to be specified.
		LoadDefault(filenames ...string) (config Config)

		// LoadSMTP loads SMTP configuration from the specified environment files.
		LoadSMTP(filenames ...string) (config SMTPConfig)

		// LoadRedis loads Redis configuration from the specified environment file(s). Returns a RedisConfig struct.
		LoadRedis(filenames ...string) (config RedisConfig)
	}

	// EnvImpl provides functionality for load `env` file
	EnvImpl struct {
		defaultConfig Config
		smtpConfig    SMTPConfig
		redisConfig   RedisConfig
	}

	// Config contains configuration settings loaded from environment variables.
	Config struct {
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
	SMTPConfig struct {
		Host     string `env:"SMTP_HOST"`
		Port     int    `env:"SMTP_PORT"`
		Username string `env:"SMTP_USERNAME"`
		Password string `env:"SMTP_PASSWORD"`
	}

	// RedisConfig holds the configuration details required to connect to a Redis client.
	RedisConfig struct {
		Host     string `env:"REDIS_HOST"`
		Port     int    `env:"REDIS_PORT"`
		Database int    `env:"REDIS_DB"`
		Username string `env:"REDIS_USERNAME"`
		Password string `env:"REDIS_PASSWORD"`
	}
)

var (
	SecretKey = paseto.NewV4AsymmetricSecretKey()
	PublicKey = SecretKey.Public()
)

// NewEnv creates a new EnvImpl instance with the provided Config, SMTPConfig, RedisConfig.
func NewEnv(defaultConfig Config, smtp SMTPConfig, redis RedisConfig) *EnvImpl {
	return &EnvImpl{defaultConfig: defaultConfig, smtpConfig: smtp, redisConfig: redis}
}

func (e *EnvImpl) LoadDefault(filenames ...string) Config {
	cfg := e.defaultConfig
	cfg = loadEnvFile[Config](filenames...)

	return cfg
}

func (e *EnvImpl) LoadSMTP(filenames ...string) SMTPConfig {
	cfg := e.smtpConfig
	cfg = loadEnvFile[SMTPConfig](filenames...)

	return cfg
}

func (e *EnvImpl) LoadRedis(filenames ...string) RedisConfig {
	cfg := e.redisConfig
	cfg = loadEnvFile[RedisConfig](filenames...)

	return cfg
}

// loadEnvFile loads the configuration from the provided `.env` files and environment variables.
func loadEnvFile[T interface{}](filenames ...string) (config T) {
	if err := godotenv.Load(filenames...); err != nil {
		log.Fatalf("Error loading `.env` file: %v", err)
	}

	if _, err := env.UnmarshalFromEnviron(&config); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	return config
}
