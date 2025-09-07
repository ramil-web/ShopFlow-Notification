package config

import (
	"os"
	"strconv"
)

type Config struct {
	RabbitURL    string
	ExchangeName string
	WorkerCount  int
	Prefetch     int
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
}

func LoadConfig() *Config {
	workerCount, _ := strconv.Atoi(getEnv("WORKER_COUNT", "5"))
	prefetch, _ := strconv.Atoi(getEnv("PREFETCH_COUNT", "10"))

	return &Config{
		RabbitURL:    getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ExchangeName: getEnv("EXCHANGE_NAME", "shopflow.events"),
		WorkerCount:  workerCount,
		Prefetch:     prefetch,
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USERNAME", ""),
		SMTPPass:     getEnv("SMTP_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
