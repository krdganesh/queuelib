package queuelib

import (
	"github.com/streadway/amqp"
)

// Connect : Function connects to RabbitMQ Server using connection string passed in Config.
// Input Parameters
//        config  : struct *Config
func (rabbitmq *RabbitMQ) Connect(config *Config) (result bool, err error) {
	if config == nil {
		return false, ErrCursor
	}

	conn, err := amqp.Dial(config.ConString)
	// defer conn.Close()
	if err != nil {
		return false, ErrConnection
	}
	rabbitmq.Connection = conn
	getChannel(rabbitmq, conn)
	return true, nil
}

func getChannel(rabbitmq *RabbitMQ, conn *amqp.Connection) (result bool, err error) {
	if rabbitmq == nil || conn == nil {
		return false, ErrCursor
	}
	ch, err := conn.Channel()
	if err != nil {
		return false, ErrConnection
	}
	rabbitmq.Channel = ch
	return true, nil
}

// Publish : Function publishes the message using existing connection object.
// Input Parameters
//        pub  : struct *PublishStruct
func (rabbitmq *RabbitMQ) Publish(pub *PublishStruct) (result bool, err error) {
	if rabbitmq == nil || rabbitmq.Channel == nil {
		return false, ErrCursor
	}
	err = rabbitmq.Channel.Publish(pub.exchange, pub.key, pub.mandatory, pub.immediate, pub.msg)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Subscribe : Function consumes the message using existing connection object.
// Input Parameters
//        sub  : struct *SubscribeStruct
func (rabbitmq *RabbitMQ) Subscribe(sub *SubscribeStruct) (delivery <-chan amqp.Delivery, err error) {
	if rabbitmq == nil || rabbitmq.Channel == nil {
		return nil, ErrCursor
	}

	msgs, err := rabbitmq.Channel.Consume(
		sub.queue,
		sub.consumer,
		sub.autoAck,
		sub.exclusive,
		sub.noLocal,
		sub.noWait,
		sub.args,
	)

	if err != nil {
		return nil, err
	}
	return msgs, nil
}
