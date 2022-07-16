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
	topic     string
	groupID   string
	partition int
	conn      *kafka.Conn
}

var app App

func init() {
	app = App{
		ctx:       context.Background(),
		topic:     "my_topic",
		groupID:   "my_group_id",
		partition: 0,
	}

	// assuming auto.create.topics.enable='true'
	conn, err := kafka.DialLeader(app.ctx, "tcp", "localhost:9092", app.topic, app.partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	app.conn = conn
}

func main() {
	app.listTopics()

	// app.produce()
	// time.Sleep(3 * time.Second)
	app.consume()

	err := app.conn.Close()
	if err != nil {
		log.Fatal("failed to close connection:", err)
	}
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
