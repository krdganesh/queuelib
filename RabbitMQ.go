package queuelib

import (
	"strconv"

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

	msg := amqp.Publishing{
		ContentType: pub.ContentType,
		Body:        pub.Message,
	}
	if pub.Delay > 0 {
		msg.Headers = amqp.Table{
			"x-delay": strconv.FormatUint(pub.Delay, 10), //Delay in milliseconds
		}
	}

	err = rabbitmq.Channel.Publish(pub.Exchange, pub.Key, pub.Mandatory, pub.Immediate, msg)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Subscribe : Function consumes the messages using existing connection object.
// Input Parameters
//        sub  : struct SubscribeStruct
func (rabbitmq *RabbitMQ) Subscribe(sub SubscribeStruct) (delivery <-chan Delivery, err error) {
	if rabbitmq == nil || rabbitmq.Channel == nil {
		return nil, ErrCursor
	}
	// ch, _ := rabbitmq.Connection.Channel()
	// defer ch.Close()
	// ch.Qos(sub.prefetchCount, sub.prefetchSize, sub.global)

	rabbitmq.Channel.Qos(sub.PrefetchCount, sub.PrefetchSize, sub.ApplyPrefetchGlobally)

	msgs, err := rabbitmq.Channel.Consume(
		sub.Queue,
		sub.Consumer,
		sub.AutoAck,
		sub.Exclusive,
		sub.NoLocal,
		sub.NoLocal,
		nil,
	)

	if err != nil {
		return nil, err
	}
	return castDeliveryCh(msgs), nil
}

// Get : Function gets a message using existing connection object.
// Input Parameters
//        get  : struct GetStruct
func (rabbitmq *RabbitMQ) Get(get GetStruct) (msg Delivery, ok bool, err error) {
	if rabbitmq == nil || rabbitmq.Channel == nil {
		return Delivery{}, false, ErrCursor
	}
	// ch, _ := rabbitmq.Connection.Channel()
	// defer ch.Close()

	delivery, ok, err := rabbitmq.Channel.Get(
		get.Queue,
		get.AutoAck,
	)

	if err != nil || ok == false {
		return Delivery{}, ok, err
	}
	return castDelivery(delivery), ok, nil
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

func castDelivery(delivery amqp.Delivery) Delivery {
	return Delivery{
		delivery.ContentType,
		delivery.ContentEncoding,
		delivery.DeliveryMode,
		delivery.Priority,
		delivery.CorrelationId,
		delivery.ReplyTo,
		delivery.Expiration,
		delivery.MessageId,
		delivery.Timestamp,
		delivery.Type,
		delivery.UserId,
		delivery.AppId,
		delivery.ConsumerTag,
		delivery.MessageCount,
		delivery.DeliveryTag,
		delivery.Redelivered,
		delivery.Exchange,
		delivery.RoutingKey,
		delivery.Body,
	}
}

func castDeliveryCh(delivery <-chan amqp.Delivery) <-chan Delivery {
	amqpDel := <-delivery
	del := Delivery{
		amqpDel.ContentType,
		amqpDel.ContentEncoding,
		amqpDel.DeliveryMode,
		amqpDel.Priority,
		amqpDel.CorrelationId,
		amqpDel.ReplyTo,
		amqpDel.Expiration,
		amqpDel.MessageId,
		amqpDel.Timestamp,
		amqpDel.Type,
		amqpDel.UserId,
		amqpDel.AppId,
		amqpDel.ConsumerTag,
		amqpDel.MessageCount,
		amqpDel.DeliveryTag,
		amqpDel.Redelivered,
		amqpDel.Exchange,
		amqpDel.RoutingKey,
		amqpDel.Body,
	}
	chDel := make(chan Delivery)
	chDel <- del
	return (<-chan Delivery)(chDel)
}
