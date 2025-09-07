package main

import (
	"log"
	"os"
	"shopflow/notification/handlers"
	"shopflow/notification/services"
	"shopflow/notification/worker"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARN] No .env file found, using system environment")
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatalln("[FATAL] RABBITMQ_URL is not set")
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalln("[FATAL] Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("[FATAL] Failed to open channel:", err)
	}
	defer ch.Close()

	exchangeName := os.Getenv("EXCHANGE_NAME")
	if exchangeName == "" {
		exchangeName = "shopflow.events"
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic", // важно: совпадает с уже существующим exchange
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("[FATAL] Failed to declare exchange:", err)
	}

	q, err := ch.QueueDeclare(
		"shopflow.notifications", // отдельная очередь для нотификейшн сервиса
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("[FATAL] Failed to declare queue:", err)
	}

	err = ch.QueueBind(
		q.Name,
		"#", // ловим все события
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("[FATAL] Failed to bind queue:", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("[FATAL] Failed to register consumer:", err)
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	emailService := services.NewEmailService(smtpUser, smtpPass, smtpHost, smtpPort)
	handler := handlers.NewHandlerContext(emailService)

	workerCount, _ := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	if workerCount == 0 {
		workerCount = 5
	}

	for i := 0; i < workerCount; i++ {
		w := worker.NewWorker(handler, conn, exchangeName, i)
		go w.Start(msgs)
	}

	log.Println("Notification service is running")
	select {}
}
