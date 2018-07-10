# Queuelib - Golang, RabbitMQ
>### This is comman queue library. Currently supports RabbitMQ.

## Instantiating the queue type
```
rmq, err := Init("rabbitmq")
```

## Connect to Queue
```
config := Config{
	ConString: "amqp://guest:guest@localhost:5672/",
}

result, err := rmq.Connect(config)
```

## Publishing Message
```
 var pubStruct = PublishStruct{
	exchange:  "testExchange",
	key:       "",
	mandatory: false,
	immediate: false,
	msg: amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Testing Publish()"),
	},
}

result, err := rmq.Publish(pubStruct)
```
## Publish with Delay
```
//To use the delay feature, an exchange with type 'x-delayed-message' must be there.

var pubDelayStruct = PublishStruct{
	exchange:  "amqp.delay",
	key:       "testKey",
	mandatory: false,
	immediate: false,
	msg: amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Testing Delayed Publish()"),
		Headers: amqp.Table{
			"x-delay": "15000", //15 sec. delay
		},
	},
}
```

## Subscribing Queue and Acknowledge
```
var subStruct = SubscribeStruct{
	queue:     "testQueue",
	consumer:  "",
	autoAck:   false,
	exclusive: false,
	noLocal:   false,
	noWait:    false,
	args:      nil,
}

chForever := make(chan bool)
msgs, err := rmq.Subscribe(subStruct)
go func() {
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
		result, err := rmq.Acknowledge(msg.DeliveryTag)
	}
}()
<-chForever
```

## Get a Message and Acknowledge
```
var getStruct = GetStruct{
	queue:   "testQueue",
	autoAck: false,
}

msg, ok, err := rmq.Get(getStruct)
log.Printf("Got a message: %s", msg.Body)

result, err := rmq.Acknowledge(msg.DeliveryTag)
```

# Project Details
## Author
```
Ganesh Karande
Email : krd.ganesh@gmail.com
```
## Version
```
1.0.0
```
## License
```
This project is licensed under the MIT License
```

