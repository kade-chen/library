package kafka_test

import (
	"context"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/kafka"
	cc "github.com/segmentio/kafka-go"
)

func TestKafkaProducer(t *testing.T) {
	m := kafka.Producer("test")
	m.WriteMessages(context.Background(),
		cc.Message{
			// 支持 Writing to multiple topics
			//  NOTE: Each Message has Topic defined, otherwise an error is returned.
			// Topic: "topic-A",
			Key:   []byte("Key-c"),
			Value: []byte("Hello World!"),
		})
	t.Log(m)
}

func TestKafkaConsumerGroup(t *testing.T) {
	m := kafka.ConsumerGroup("test", []string{"test"})
	t.Log(m)
}

func init() {
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
