package mq

import (
	"github.com/bitly/go-nsq"
	"github.com/tylertreat/mq-benchmarking/benchmark"
)

type Nsq struct {
	handler *benchmark.MessageHandler
	pub     *nsq.Producer
	sub     *nsq.Consumer
	topic   string
	channel string
}

func NewNsq(numberOfMessages int) Nsq {
	topic := "test"
	channel := "test"
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	sub, _ := nsq.NewConsumer(topic, channel, nsq.NewConfig())

	return Nsq{
		handler: &benchmark.MessageHandler{NumberOfMessages: numberOfMessages},
		pub:     pub,
		sub:     sub,
		topic:   topic,
		channel: channel,
	}
}

func (n Nsq) Setup() {
	n.sub.SetHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		n.ReceiveMessage(message.Body)
		return nil
	}))
	n.sub.ConnectToNSQD("localhost:4150")
}

func (n Nsq) Teardown() {
	n.sub.Stop()
	n.pub.Stop()
}

func (n Nsq) Send(message []byte) {
	n.pub.Publish(n.topic, message)
}

func (n Nsq) ReceiveMessage(message []byte) bool {
	return n.handler.ReceiveMessage(message)
}

func (n Nsq) MessageHandler() *benchmark.MessageHandler {
	return n.handler
}