package main

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type App struct {
	ctx      context.Context
	topic    string
	groupID  string
	consumer *kafka.Consumer
	producer *kafka.Producer
}

var app App

func init() {
	app = App{
		ctx:     context.Background(),
		topic:   "my_topic",
		groupID: "my_group_id",
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          app.groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	app.consumer = consumer

	err = app.consumer.Subscribe(app.topic, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	app.producer = producer
}

func (a *App) Shutdown() {
	a.consumer.Close()
	a.producer.Close()
}

func main() {
	defer app.Shutdown()

	app.produceMessage("test message")

	app.consumeMessages()
}

func (a *App) consumeMessages() {
	for {
		msg, err := app.consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func (a *App) produceMessage(msg string) {
	a.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &app.topic, Partition: kafka.PartitionAny},
		Value: []byte(msg),
	}, nil)

	a.producer.Flush(1 * 1000)
}