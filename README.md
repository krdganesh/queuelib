# Queuelib - Golang, RabbitMQ
>### This is common queue library. It currently supports RabbitMQ.

## Getting queuelib
```
go get -u github.com/krdganesh/queuelib
```

## Instantiating the queue type
```
rmq, err := queuelib.Init("rabbitmq")
```

## Connect to Queue
```
config := &queuelib.Config{
	ConString: "amqp://guest:guest@localhost:5672/",
}

result, err := rmq.Connect(config)
```

## Publishing Message
```
var pubStruct = queuelib.PublishStruct{
	Exchange:    "testExchange",
	Key:         "",
	Mandatory:   false,
	Immediate:   false,
	ContentType: "text/plain",
	Message:     []byte("Testing Publish()"),
}

result, err := rmq.Publish(pubStruct)
```
## Publish with Delay
```
//To use the delay feature, an exchange with type 'x-delayed-message' must be there.

var pubDelayStruct = queuelib.PublishStruct{
	Exchange:    "amqp.delay",
	Key:         "testKey",
	Mandatory:   false,
	Immediate:   false,
	ContentType: "text/plain",
	Message:     []byte("Testing Delayed Publish()"),
	Delay:       15000, //15 sec. delay
}

result, err := rmq.Publish(pubDelayStruct)
```

## Subscribing Queue and Acknowledge
```
var subStruct = queuelib.SubscribeStruct{
	Queue:     "testQueue",
	Consumer:  "",
	AutoAck:   false,
	Exclusive: false,
	NoLocal:   false,
	NoWait:    false,
	PrefetchCount: 10, //Allows batching of messages
}

chForever := make(chan bool)
msgs, err := rmq.Subscribe(subStruct)
go func() {
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
		result, err := rmq.Acknowledge(msg)
	}
}()
<-chForever
```

## Get a Message and Acknowledge
```
var getStruct = queuelib.GetStruct{
	Queue:   "testQueue",
	AutoAck: false,
}

msg, ok, err := rmq.Get(getStruct)
log.Printf("Got a message: %s", msg.Body)

result, err := rmq.Acknowledge(msg)
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

