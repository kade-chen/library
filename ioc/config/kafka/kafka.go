package kafka

import (
	"fmt"
	"time"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func init() {
	ioc.Config().Registry(&Kafka{
		Brokers:        []string{"127.0.0.1:9092"},
		ScramAlgorithm: SHA512,
	})
}

type Kafka struct {
	Brokers        []string       `toml:"brokers" json:"brokers" yaml:"brokers"  env:"KAFKA_BROKERS"`
	ScramAlgorithm ScramAlgorithm `toml:"scram_algorithm" json:"scram_algorithm" yaml:"scram_algorithm"  env:"KAFKA_SCRAM_ALGORITHM"`
	UserName       string         `toml:"username" json:"username" yaml:"username"  env:"KAFKA_USERNAME"`
	Password       string         `toml:"password" json:"password" yaml:"password"  env:"KAFKA_PASSWORD"`
	Debug          bool           `toml:"debug" json:"debug" yaml:"debug"  env:"KAFKA_DEBUG"`

	mechanism sasl.Mechanism
	ioc.ObjectImpl
}

func (m *Kafka) Name() string {
	return AppName
}

func (k *Kafka) Init() error {
	if k.UserName != "" {
		mechanism, err := scram.Mechanism(k.scramAlgorithm(), k.UserName, k.Password)
		if err != nil {
			return err
		}
		k.mechanism = mechanism
	}

	return nil
}

func (k *Kafka) scramAlgorithm() scram.Algorithm {
	switch k.ScramAlgorithm {
	case SHA256:
		return scram.SHA256
	default:
		return scram.SHA512
	}
}

func (k *Kafka) Producer(topic string) *kafka.Writer {
	l := log.Sub(fmt.Sprintf("producer_%s", topic))

	w := &kafka.Writer{
		Addr:                   kafka.TCP(k.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		ErrorLogger:            kafka.LoggerFunc(l.Error().Msgf),
		Transport:              &kafka.Transport{SASL: k.mechanism},
	}
	if k.Debug {
		w.Logger = kafka.LoggerFunc(l.Debug().Msgf)
	}
	return w
}

func (k *Kafka) ConsumerGroup(groupId string, topics []string) *kafka.Reader {
	l := log.Sub(fmt.Sprintf("consumer_group_%s", groupId))
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: k.mechanism,
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     k.Brokers,
		Dialer:      dialer,
		GroupID:     groupId,
		GroupTopics: topics,
		Logger:      kafka.LoggerFunc(l.Debug().Msgf),
		ErrorLogger: kafka.LoggerFunc(l.Error().Msgf),
	})
	
}
