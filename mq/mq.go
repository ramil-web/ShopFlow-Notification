package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type MQConnection struct {
	Conn *amqp.Connection
}

func NewMQConnection(url string) (*MQConnection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MQConnection{Conn: conn}, nil
}

func (m *MQConnection) Channel() (*amqp.Channel, error) {
	return m.Conn.Channel()
}

func (m *MQConnection) Close() error {
	return m.Conn.Close()
}
