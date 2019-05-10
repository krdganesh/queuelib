package queuelib

import (
	"context"
	"errors"
	"strconv"
	"time"

	k "github.com/segmentio/kafka-go"
)

// Connect : Function connects to Kafka Server using connection string passed in Config.
// Input Parameters
//        config  : struct *Config
func (kafka *Kafka) Connect(config *Config) (result bool, err error) {
	if config == nil {
		return false, ErrCursor
	}

	conn, err := k.DialLeader(context.Background(), "tcp", config.Host+":"+strconv.FormatInt(int64(config.Port), 10), config.KafkaTopic, config.KafkaPartition)
	if err != nil {
		return false, ErrConnection
	}

	kafka.Connection = conn
	return true, nil
}

// Publish : Function publishes the message using existing connection object.
// Input Parameters
//        pub  : struct PublishStruct
func (kafka *Kafka) Publish(pub PublishStruct) (result bool, err error) {
	if kafka == nil {
		return false, ErrCursor
	}

	kafka.Connection.SetWriteDeadline(time.Now().Add(pub.KafkaTimeoutInSeconds * time.Second))
	_, err = kafka.Connection.Write(pub.Message)

	if err != nil {
		return false, err
	}
	return true, nil
}

// Subscribe : Function consumes the messages using existing connection object.
// Input Parameters
//        sub  : struct SubscribeStruct
func (kafka *Kafka) Subscribe(sub SubscribeStruct) (delivery <-chan Delivery, err error) {
	if kafka == nil {
		return nil, ErrCursor
	}

	r := k.NewReader(k.ReaderConfig{
		Brokers:  sub.KafkaBrokers,
		GroupID:  sub.KafkaConsumerGroupID,
		Topic:    sub.KafkaTopic,
		MinBytes: 1,    // 10Byte
		MaxBytes: 10e6, // 10MB
	})

	chDel := make(chan Delivery)

	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				break
			}
			del := castDeliveryKafka(m)
			chDel <- del
		}
	}()

	return (<-chan Delivery)(chDel), nil
}

func castDeliveryKafka(delivery k.Message) Delivery {
	return Delivery{
		KafkaTopic:     delivery.Topic,
		KafkaPartition: delivery.Partition,
		KafkaKey:       delivery.Key,
		Body:           delivery.Value,
	}
}

// Get : Function gets a message using existing connection object.
// Input Parameters
//        get  : struct GetStruct
func (kafka *Kafka) Get(get GetStruct) (msg Delivery, ok bool, err error) {
	return Delivery{}, false, errors.New("Not implemented Get")
}

// Acknowledge : Function acknowledges a message using existing connection object.
// Input Parameters
//        DeliveryTag  : uint64
func (kafka *Kafka) Acknowledge(delivery Delivery) (result bool, err error) {
	return false, errors.New("Not implemented Get")
}
