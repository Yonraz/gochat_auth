package producers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/yonraz/gochat_auth/constants"
)

type Producer struct {
	topic constants.Topic
	Producer sarama.SyncProducer
}

func (p *Producer) Produce(msg interface{}) error {
	messageAsBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshalling msg to json: %v", err)
	}

	message := &sarama.ProducerMessage{
		Topic: string(p.topic),
		Value: sarama.StringEncoder(messageAsBytes),
	}

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	log.Printf("message stored on topic(%s)/partition(%d)/offset(%d)\n", p.topic, partition, offset)

	log.Printf("produced message %v to topic %v\n", msg, p.topic)

	return nil
}

func NewUserRegisteredProducer(brokers []string) (*Producer, error) {
	p, err := CreateProducer(brokers)
	
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %v", err)
	}
	
	return &Producer{
		topic: constants.UserRegistered,
		Producer: p,
		}, err
}
func NewUserLoggedinProducer(brokers []string) (*Producer, error) {
	p, err := CreateProducer(brokers)
	
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %v", err)
	}

	return &Producer{
		topic: constants.UserLoggedIn,
		Producer: p,
	}, err
}
func NewUserSignedoutProducer(brokers []string) (*Producer, error) {
	p, err := CreateProducer(brokers)
	
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %v", err)
	}

	return &Producer{
		topic: constants.UserSignedout,
		Producer: p,
	}, err
}


