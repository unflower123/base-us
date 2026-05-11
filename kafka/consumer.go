package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topics string, groupId string) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topics,
		MaxBytes:       10e6,
		MinBytes:       10e3,
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
		GroupID:        groupId,
	})

	logx.Info("kafka reader init  success")
	c := &Consumer{
		reader: reader,
	}

	return c, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()

}

func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) error {
	for {
		/*readMsg, err := reader.FetchMessage(ctx)*/
		readMsg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			logx.Errorf(" kafka reader failed to read msg, err:%v", err)
			continue
		}
		success := handler(&readMsg)
		if success {
			if err := c.reader.CommitMessages(ctx, readMsg); err != nil {
				logx.Errorf("submit kafka message error:%s", err)
				return err
			}
		} else {
			//todo retry
			//todo add dead queue
		}
	}

	if err := c.Close(); err != nil {
		logx.Errorf("kafka reader close error: ", err)
	}

	logx.Info("kafka reader close success")

	return nil
}

type MessageHandler func(message *kafka.Message) bool
