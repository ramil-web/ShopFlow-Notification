package worker

import (
	"log"
	"shopflow/notification/handlers"

	"github.com/streadway/amqp"
)

type Worker struct {
	handler      *handlers.HandlerContext
	connection   *amqp.Connection
	exchangeName string
	id           int
}

func NewWorker(handler *handlers.HandlerContext, conn *amqp.Connection, exchangeName string, id int) *Worker {
	return &Worker{
		handler:      handler,
		connection:   conn,
		exchangeName: exchangeName,
		id:           id,
	}
}

func (w *Worker) Start(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Printf("[INFO] Worker %d received message (routing key=%s): %s\n",
			w.id, d.RoutingKey, string(d.Body))

		log.Printf("[DEBUG] Worker %d received routing key=%s body=%s\n",
			w.id, d.RoutingKey, string(d.Body))

		switch d.RoutingKey {
		case "user.registered":
			w.handler.ProcessRegistration(d.Body)

		case "application_created":
			w.handler.ProcessApplication(d.Body)

		default:
			log.Printf("[WARN] Worker %d got unknown event: %s\n", w.id, d.RoutingKey)
		}
	}
}
