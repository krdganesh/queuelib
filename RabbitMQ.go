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
//        pub  : struct PublishStruct
func (rabbitmq *RabbitMQ) Publish(pub PublishStruct) (result bool, err error) {
	if rabbitmq == nil || rabbitmq.Channel == nil {
		return false, ErrCursor
	}
	// ch, _ := rabbitmq.Connection.Channel()
	// defer ch.Close()
	// err = ch.Publish(pub.exchange, pub.key, pub.mandatory, pub.immediate, pub.msg)

	err = rabbitmq.Channel.Publish(pub.exchange, pub.key, pub.mandatory, pub.immediate, pub.msg)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Subscribe : Function consumes the messages using existing connection object.
// Input Parameters
//        sub  : struct SubscribeStruct
func (rabbitmq *RabbitMQ) Subscribe(sub SubscribeStruct) (delivery <-chan amqp.Delivery, err error) {
	if rabbitmq == nil || rabbitmq.Channel == nil {
		return nil, ErrCursor
	}
	// ch, _ := rabbitmq.Connection.Channel()
	// defer ch.Close()
	// ch.Qos(sub.prefetchCount, sub.prefetchSize, sub.global)

	rabbitmq.Channel.Qos(sub.prefetchCount, sub.prefetchSize, sub.global)

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

// Get : Function gets a message using existing connection object.
// Input Parameters
//        get  : struct GetStruct
func (rabbitmq *RabbitMQ) Get(get GetStruct) (msg amqp.Delivery, ok bool, err error) {
	if rabbitmq == nil || rabbitmq.Channel == nil {
		return amqp.Delivery{}, false, ErrCursor
	}
	// ch, _ := rabbitmq.Connection.Channel()
	// defer ch.Close()

	msg, ok, err = rabbitmq.Channel.Get(
		get.queue,
		get.autoAck,
	)

	if err != nil || ok == false {
		return amqp.Delivery{}, ok, err
	}
	return msg, ok, nil
}

// Acknowledge : Function acknowledges a message using existing connection object.
// Input Parameters
//        DeliveryTag  : uint64
func (rabbitmq *RabbitMQ) Acknowledge(DeliveryTag uint64) (result bool, err error) {
	// ch, err := rabbitmq.Connection.Channel()
	err = rabbitmq.Channel.Ack(DeliveryTag, true)
	if err != nil {
		return false, err
	}
	return true, nil
}
