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
		log.Printf("[INFO] Worker %d received message: %s\n", w.id, string(d.Body))
		w.handler.ProcessRegistration(d.Body)
	}
}
