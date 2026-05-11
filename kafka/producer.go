package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true, //IPV4, IPV6
	}
	config := kafka.WriterConfig{
		Brokers:  brokers,
		Dialer:   dialer,
		Topic:    fmt.Sprintf("%s", topic),
		Balancer: &kafka.LeastBytes{},
		//RequiredAcks: kafka.RequireAll,
		Async: true,
	}

	p := &Producer{
		writer: kafka.NewWriter(config),
	}

	logx.Info("kafka producer init success!")
	return p, nil
}

func (p *Producer) Produce(message []byte) (err error) {

	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err = p.writer.WriteMessages(ctx, kafka.Message{Value: message})
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}
		if err != nil {
			logx.Errorf("Failed to write messages. err:%s", err.Error())
			continue
		}

		logx.Info("kafka write messages successfully, req: " + fmt.Sprintf("%+v\n", message))
		break
	}

	return
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
