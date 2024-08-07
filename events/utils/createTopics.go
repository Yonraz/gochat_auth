package utils

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/yonraz/gochat_auth/constants"
)

func CreateTopic(brokerList []string, topic string) error {
	config := sarama.NewConfig()
    config.Version = sarama.V2_8_0_0 // Kafka version

    admin, err := sarama.NewClusterAdmin(brokerList, config)
    if err != nil {
        return fmt.Errorf("failed to create cluster admin: %w", err)
    }
    defer admin.Close()

    topicDetail := sarama.TopicDetail{
        NumPartitions:     1,
        ReplicationFactor: 1,
    }

    err = admin.CreateTopic(topic, &topicDetail, false)
    if err != nil {
        if err.(*sarama.TopicError).Err == sarama.ErrTopicAlreadyExists {
            log.Printf("Topic %s already exists", topic)
        } else {
            return fmt.Errorf("failed to create topic: %w", err)
        }
    } else {
        log.Printf("Topic %s created successfully", topic)
    }

    return nil
}

func CreateTopics(brokers []string, topics []constants.Topic) {
	for _, topic :=  range topics {
		CreateTopic(brokers, string(topic))
	}
}