package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AuthRPC struct {
	conn *amqp.Connection
}

func NewAuthRPC(conn *amqp.Connection) *AuthRPC {
	return &AuthRPC{conn: conn}
}

type EmailResponse struct {
	Email string `json:"email"`
}

func (r *AuthRPC) GetUserEmail(ctx context.Context, userID string) (*EmailResponse, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	replyQueue, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return nil, err
	}

	corrID := fmt.Sprintf("%d", userID)

	err = ch.Publish(
		"",                    // exchange
		"auth.user.email.rpc", // routing key
		false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          []byte(fmt.Sprintf(`{"user_id":%s}`, userID)),
			CorrelationId: corrID,
			ReplyTo:       replyQueue.Name,
		},
	)
	if err != nil {
		return nil, err
	}

	msgs, _ := ch.Consume(replyQueue.Name, "", true, false, false, false, nil)
	for d := range msgs {
		if d.CorrelationId == corrID {
			var resp EmailResponse
			if err := json.Unmarshal(d.Body, &resp); err != nil {
				return nil, err
			}
			return &resp, nil
		}
	}
	return nil, fmt.Errorf("no response received")
}
