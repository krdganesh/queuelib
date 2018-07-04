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

result, err := rmq.Publish(&pubStruct)
```

## Subscribing Queue
```
var subscribeStruct = SubscribeStruct{
	queue:     "testQueue",
	consumer:  "",
	autoAck:   false,
	exclusive: false,
	noLocal:   false,
	noWait:    false,
	args:      nil,
}

chForever := make(chan bool)
msgs, err := rmq.Subscribe(&subscribeStruct)
go func() {
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
	}
}()
<-chForever
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

