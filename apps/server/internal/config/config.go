package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppEnv string         `yaml:"env" env-default:"local" env:"APP_ENV"`
	Server HttpServer     `yaml:"http_server"`
	DB     DatabaseConfig `yaml:"db"`
	S3     S3Config       `yaml:"s3"`
}

type HttpServer struct {
	Host          string        `yaml:"host" env-default:"localhost" env:"API_HOST"`
	Port          string        `yaml:"port" env-default:"8000" env:"API_PORT"`
	ReadTimeout   time.Duration `yaml:"read_timeout" env-default:"10s" env:"HTTP_READ_TIMEOUT" `
	WriteTimeout  time.Duration `yaml:"write_timeout" env-default:"10s" env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout   time.Duration `yaml:"idle_timeout" env-default:"60s" env:"HTTP_IDLE_TIMEOUT"`
	JwtSecret     string        `yaml:"jwt_secret" env-default:"super_secret_key" env:"JWT_SECRET"`
	SecureCookies bool          `yaml:"secure_cookies" env-default:"false" env:"SECURE_COOKIES"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env-default:"localhost" env:"DB_HOST"`
	Port     string `yaml:"port" env-default:"3306" env:"DB_PORT"`
	Username string `yaml:"username" env-default:"default" env:"DB_USER_NAME"`
	Password string `yaml:"password" env-default:"Qwerty12345" env:"DB_PASSWORD"`
	DbName   string `yaml:"dbname" env-default:"pastebin" env:"DB_NAME"`
	SslMode  string `yaml:"sslmode" env-default:"disable" env:"SSL_MODE"`
}

type S3Config struct {
	Endpoint  string `yaml:"endpoint" env-default:"https://pupuha-dev.object.pscloud.io" env:"S3_ENDPOINT"`
	Region    string `yaml:"region" env-default:"us-east-1" env:"S3_REGION"`
	Bucket    string `yaml:"bucket" env-default:"pupuha-dev" env:"S3_BUCKET"`
	AccessKey string `yaml:"access_key" env-default:"V9M58P398E715KZHYNPS" env:"S3_ACCESS_KEY"`
	SecretKey string `yaml:"secret_key" env-default:"IO24QaIQ1i0xROcnKF2fwO77CyI7E0TXmFs88ffq" env:"S3_SECRET_KEY"`
}

func GetConfig(configPath string) *Config {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to init config", err.Error())
	}

	return &cfg
}
