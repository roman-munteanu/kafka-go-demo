package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type App struct {
	ctx       context.Context
	host      string
	topic     string
	groupID   string
	partition int
	conn      *kafka.Conn
	reader    *kafka.Reader
	writer    *kafka.Writer
}

var app App

func init() {
	app = App{
		ctx:       context.Background(),
		host:      "localhost:9092",
		topic:     "my_topic",
		groupID:   "my_group_id",
		partition: 0,
	}

	// assuming auto.create.topics.enable='true'
	conn, err := kafka.DialLeader(app.ctx, "tcp", app.host, app.topic, app.partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	app.conn = conn

	app.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{app.host},
		GroupID:  app.groupID,
		Topic:    app.topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	// app.writer = kafka.NewWriter(kafka.WriterConfig{
	// 	Brokers:  []string{app.host},
	// 	Topic:    app.topic,
	// 	Balancer: &kafka.LeastBytes{},
	// })

	app.writer = &kafka.Writer{
		Addr:     kafka.TCP(app.host),
		Topic:    app.topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func main() {
	defer app.Shutdown()

	app.listTopics()

	// app.produce()
	// time.Sleep(3 * time.Second)
	// app.consume()

	app.writeMessages()
	time.Sleep(3 * time.Second)
	app.readMessages()
}

func (a *App) consume() {
	a.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	// fetch 10KB min, 1MB max
	batch := a.conn.ReadBatch(10e3, 10e6)

	// 10KB max per message
	buff := make([]byte, 10e3)
	for {
		n, err := batch.Read(buff)
		if err != nil {
			break
		}
		fmt.Println(string(buff[:n]))

	}

	err := batch.Close()
	if err != nil {
		log.Fatal("failed to close batch:", err)
	}
}

func (a *App) produce() {
	a.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	_, err := a.conn.WriteMessages(
		kafka.Message{Value: []byte("test1")},
		kafka.Message{Value: []byte("test2")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}

func (a *App) listTopics() {
	partitions, err := a.conn.ReadPartitions()
	if err != nil {
		log.Fatal("could not read partitions:", err)
		return
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}

	for k := range m {
		fmt.Println(k)
	}
}

func (a *App) readMessages() {
	for {
		m, err := a.reader.ReadMessage(a.ctx)
		if err != nil {
			break
		}
		fmt.Println("-------------- msg:")
		fmt.Println("Topic:", m.Topic)
		fmt.Println("Partition:", m.Partition)
		fmt.Println("Offset:", m.Offset)
		fmt.Println("Key:", string(m.Key))
		fmt.Println("Value:", string(m.Value))
	}
}

func (a *App) writeMessages() {
	err := a.writer.WriteMessages(a.ctx,
		kafka.Message{
			Key:   []byte("routingKey1"),
			Value: []byte("message 1"),
		},
		kafka.Message{
			Key:   []byte("routingKey2"),
			Value: []byte("message 2"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}

func (a *App) Shutdown() {
	err := a.conn.Close()
	if err != nil {
		log.Fatal("failed to close connection:", err)
	}

	err = a.reader.Close()
	if err != nil {
		log.Fatal("failed to close reader:", err)
	}

	err = a.writer.Close()
	if err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
