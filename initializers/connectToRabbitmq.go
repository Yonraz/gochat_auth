package initializers

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_auth/constants"
	"github.com/yonraz/gochat_auth/events/utils"
)

var RmqChannel *amqp.Channel
var RmqConn *amqp.Connection

func ConnectToRabbitmq() {
	var err error
	RmqConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Connected to Rabbitmq")

	RmqChannel, err = RmqConn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Declaring the topic exchange
	err = RmqChannel.ExchangeDeclare(
		string(constants.UserEventsExchange),
		"topic",             
		true,                
		false,               
		false, 
		false,              
		nil,                 
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// create registration, login, logout queues
	err = utils.DeclareQueues(RmqChannel)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}