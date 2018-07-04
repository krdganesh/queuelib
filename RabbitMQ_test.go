package queuelib

import (
	"log"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

var config = &Config{
	ConString: "amqp://guest:guest@localhost:5672/",
}

var configIncorrect = &Config{
	ConString: "amqp://a:a@localhost:5672/",
}

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

var subscribeStruct = SubscribeStruct{
	queue:     "testQueue",
	consumer:  "",
	autoAck:   false,
	exclusive: false,
	noLocal:   false,
	noWait:    false,
	args:      nil,
}

var rmq = new(RabbitMQ)
var rmqIncorrect = new(RabbitMQ)

//It will test all scenarios of Connect()
func TestConnect(t *testing.T) {

	//Checking case where config is nil
	r, err := rmq.Connect(nil)
	assert.Equal(t, ErrCursor, err)
	assert.Equal(t, false, r)

	//Checking case where login id and password are wrongly passed
	r1, err := rmq.Connect(configIncorrect)
	assert.NotNil(t, err)
	assert.Equal(t, false, r1)

	//Checking case where entire connection string is wrongly passed
	configIncorrect.ConString = "abcdefg"
	r2, err := rmq.Connect(configIncorrect)
	assert.NotNil(t, err)
	assert.Equal(t, false, r2)

	result, err := rmq.Connect(config)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}

//It will test all scenarios of getChannel()
func TestGetChannel(t *testing.T) {

	//Checking case where *RabbitMQ and *amqp.Connection are nil
	r, err := getChannel(nil, nil)
	assert.Equal(t, false, r)
	assert.Equal(t, err, ErrCursor)

	r1, err := getChannel(rmqIncorrect, rmqIncorrect.Connection)
	assert.Equal(t, false, r1)
	assert.NotNil(t, err)

	result, err := getChannel(rmq, rmq.Connection)
	assert.Equal(t, true, result)
	assert.Nil(t, err)
}

//It will test all scenarios of Publish()
func TestPublish(t *testing.T) {
	r, err := rmqIncorrect.Publish(&pubStruct)
	assert.NotNil(t, err)
	assert.Equal(t, false, r)

	result, err := rmq.Publish(&pubStruct)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}

//It will test all scenarios of Subscribe()
func TestSubscribe(t *testing.T) {
	r, err := rmqIncorrect.Subscribe(&subscribeStruct)
	assert.NotNil(t, err)
	assert.Nil(t, r)

	chStop := make(chan bool)
	msgs, err := rmq.Subscribe(&subscribeStruct)
	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			chStop <- true
			break
		}
	}()
	result := <-chStop
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}
