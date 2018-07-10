package queuelib

import (
	"github.com/streadway/amqp"
)

//Queue : This interface is used as return type of Init()
type Queue interface {
	Connect(*Config) (bool, error)
}

//RabbitMQ : Pointer to this struct is retured in Init() if input QueueType is "rabbitmq"
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

//Sample Queue Type
//type Kafka struct{
//
//}

//Config : This struct is used to define all neccessary parameters required by Supported Queue Client i.e. RabbitMQ (As of now)
type Config struct {
	ConString string
	Scheme    string
	Host      string
	Port      int
	Username  string
	Password  string
	Vhost     string
}

//PublishStruct : This struct is a input parameter for Publish()
type PublishStruct struct {
	exchange  string
	key       string
	mandatory bool
	immediate bool
	msg       amqp.Publishing
	delay     uint64
}

//SubscribeStruct : This struct is a input parameter for Subscribe()
type SubscribeStruct struct {
	queue         string
	consumer      string
	autoAck       bool
	exclusive     bool
	noLocal       bool
	noWait        bool
	args          amqp.Table
	prefetchCount int
	prefetchSize  int
	global        bool //apply prefetch settings to all channels - across all consumers
}

//GetStruct : This struct is a input parameter for Get()
type GetStruct struct {
	queue   string
	autoAck bool
}
