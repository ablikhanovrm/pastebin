package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string `env:"APP_ENV" required:"true"`
	Server HttpServer
	DB     DatabaseConfig
	S3     S3Config
	Redis  RedisConfig
}

type HttpServer struct {
	Host          string        `env:"HTTP_BIND_ADDRESS" required:"true"`
	Port          string        `env:"HTTP_PORT" required:"true"`
	ReadTimeout   time.Duration `env-default:"10s" env:"HTTP_READ_TIMEOUT" `
	WriteTimeout  time.Duration `env-default:"10s" env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout   time.Duration `env-default:"60s" env:"HTTP_IDLE_TIMEOUT"`
	JwtSecret     string        `env:"JWT_SECRET" required:"true"`
	SecureCookies bool          `env-default:"false" env:"SECURE_COOKIES" required:"true"`
}

type DatabaseConfig struct {
	Host             string `env-default:"localhost" env:"DB_HOST" required:"true"`
	Port             string `env-default:"5432" env:"DB_PORT" required:"true"`
	Username         string `env-default:"default" env:"DB_USER" required:"true"`
	Password         string `env-default:"default" env:"DB_PASSWORD" required:"true"`
	MigratorUser     string `env-default:"default" env:"DB_MIGRATOR_USER" required:"true"`
	MigratorPassword string `env-default:"default" env:"DB_MIGRATOR_PASSWORD" required:"true"`
	DbName           string `env-default:"pastebin" env:"DB_NAME" required:"true"`
	SslMode          string `env-default:"enable" env:"DB_SSLMODE" required:"true"`
}

type S3Config struct {
	Endpoint  string `env:"S3_ENDPOINT" required:"true"`
	Region    string `env:"S3_REGION" required:"true"`
	Bucket    string `env:"S3_BUCKET" required:"true"`
	AccessKey string `env:"S3_ACCESS_KEY" required:"true"`
	SecretKey string `env:"S3_SECRET_KEY" required:"true"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" required:"true"`
	Port     string `env:"REDIS_PORT" required:"true"`
	Password string `env:"REDIS_PASSWORD" required:"true"`
}

func GetConfig() *Config {
	var cfg Config

	if os.Getenv("APP_ENV") != "prod" {
		_ = godotenv.Load()
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to init config: %v", err.Error())
	}

	return &cfg
}
