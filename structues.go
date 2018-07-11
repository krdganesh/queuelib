package queuelib

import (
	"time"

	"github.com/streadway/amqp"
)

//Queue : This interface is used as return type of Init()
type Queue interface {
	Connect(*Config) (bool, error)
	Publish(PublishStruct) (bool, error)
	Subscribe(SubscribeStruct) (<-chan Delivery, error)
	Get(GetStruct) (Delivery, bool, error)
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

//PublishStruct : This struct is an input parameter for Publish()
type PublishStruct struct {
	Exchange    string
	Key         string
	Mandatory   bool
	Immediate   bool
	Message     []byte
	ContentType string
	DeliveryTag uint64
	Delay       uint64 //Delay in milliseconds
}

//SubscribeStruct : This struct is an input parameter for Subscribe()
type SubscribeStruct struct {
	Queue                 string
	Consumer              string
	AutoAck               bool
	Exclusive             bool
	NoLocal               bool
	NoWait                bool
	PrefetchCount         int
	PrefetchSize          int
	ApplyPrefetchGlobally bool //apply prefetch settings to all channels - across all consumers
}

//GetStruct : This struct is an input parameter for Get()
type GetStruct struct {
	Queue   string
	AutoAck bool
}

// Delivery captures the fields for a previously delivered message resident in
// a queue to be delivered by the server to a consumer from Channel.Consume or
// Channel.Get.
type Delivery struct {
	// Properties
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	DeliveryMode    uint8     // queue implementation use - non-persistent (1) or persistent (2)
	Priority        uint8     // queue implementation use - 0 to 9
	CorrelationID   string    // application use - correlation identifier
	ReplyTo         string    // application use - address to reply to (ex: RPC)
	Expiration      string    // implementation use - message expiration spec
	MessageID       string    // application use - message identifier
	Timestamp       time.Time // application use - message timestamp
	Type            string    // application use - message type name
	UserID          string    // application use - creating user - should be authenticated user
	AppID           string    // application use - creating application id

	// Valid only with Channel.Consume
	ConsumerTag string

	// Valid only with Channel.Get
	MessageCount uint32

	DeliveryTag uint64
	Redelivered bool
	Exchange    string // basic.publish exhange
	RoutingKey  string // basic.publish routing key

	Body []byte
}
