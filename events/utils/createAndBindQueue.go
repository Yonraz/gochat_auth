package utils

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_auth/constants"
)

func DeclareAndBindQueue(
	channel *amqp.Channel,
	queueName constants.Queues,
	routingKey constants.RoutingKey,
	exchangeName constants.Exchange,
	durable bool,
	autoDelete bool,
	exclusive bool,
	args amqp.Table,
) error {
	_, err := channel.QueueDeclare(
		string(queueName),
		durable,
		autoDelete,
		exclusive,
		false,
		args,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = channel.QueueBind(
		string(queueName),
		string(routingKey),
		string(exchangeName),
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

func DeclareQueues(channel *amqp.Channel) error {
	// create registration queue
	err := DeclareAndBindQueue(
		channel,
		constants.UserRegistrationQueue,
		constants.UserRegisteredKey,
		constants.UserEventsExchange,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// create login queue
	err = DeclareAndBindQueue(
		channel,
		constants.UserLoginQueue,
		constants.UserLoggedInKey,
		constants.UserEventsExchange,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	// create logout queue
	err = DeclareAndBindQueue(
		channel,
		constants.UserSignoutQueue,
		constants.UserSignedoutKey,
		constants.UserEventsExchange,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}