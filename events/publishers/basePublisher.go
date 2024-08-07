package publishers

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_auth/constants"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) *Publisher {
	return &Publisher{channel: channel}
}

func (p *Publisher) Publish(exchange constants.Exchange, routingKey constants.Topic, body interface{}) error {
	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return fmt.Errorf("failed to marshal body to json: %w", err)
	}

	err = p.channel.Publish(
		string(exchange),     // Exchange name
		string(routingKey),   // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	fmt.Printf("Published %v to exchange %v with key %v\n", body, exchange, routingKey)
	
	return nil
}