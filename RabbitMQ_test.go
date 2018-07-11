package queuelib

import (
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var messagesCount = 10000 //No. of messages to be published

var config = &Config{
	ConString: "amqp://guest:guest@localhost:5672/",
}

var configIncorrect = &Config{
	ConString: "amqp://a:a@localhost:5672/",
}

var pubStruct = PublishStruct{
	Exchange:    "testExchange",
	Key:         "",
	Mandatory:   false,
	Immediate:   false,
	ContentType: "text/plain",
	Message:     []byte("Testing Publish()"),
}

var pubDelayStruct = PublishStruct{
	Exchange:    "amqp.delay",
	Key:         "testKey",
	Mandatory:   false,
	Immediate:   false,
	ContentType: "text/plain",
	Message:     []byte("Testing Delayed Publish()"),
	Delay:       15000, //15 sec. delay
}

var subStruct = SubscribeStruct{
	Queue:     "testQueue",
	Consumer:  "",
	AutoAck:   false,
	Exclusive: false,
	NoLocal:   false,
	NoWait:    false,
}

var getStruct = GetStruct{
	Queue:   "testQueue",
	AutoAck: false,
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
	r, err := rmqIncorrect.Publish(pubStruct)
	assert.NotNil(t, err)
	assert.Equal(t, false, r)

	r1, err := rmq.Connect(config)
	assert.Nil(t, err)
	assert.Equal(t, true, r1)

	//Publishing delayed message
	result, err := rmq.Publish(pubDelayStruct)
	assert.Nil(t, err)
	assert.Equal(t, true, result)

	//Publishing normal messages
	chStop := make(chan bool)
	for i := 1; i <= messagesCount; i++ {
		go func(i int) {
			pubStruct.Message = []byte("Testing message - " + strconv.Itoa(i))
			result, err := rmq.Publish(pubStruct)
			assert.Nil(t, err)
			assert.Equal(t, true, result)
			if i == messagesCount {
				chStop <- true
			}
		}(i)
	}
	<-chStop
}

//It will test all scenarios of Get()
func TestGet(t *testing.T) {
	r, ok, err := rmqIncorrect.Get(getStruct)
	assert.NotNil(t, err)
	assert.Equal(t, false, ok)
	assert.Nil(t, r.Body)

	msg, ok, err := rmq.Get(getStruct)
	log.Printf("Got a message: %s", msg.Body)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.NotNil(t, msg.Body)

	result, err := rmq.Acknowledge(msg.DeliveryTag)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}

// It will test all scenarios of Subscribe()
func TestSubscribe(t *testing.T) {
	r, err := rmqIncorrect.Subscribe(subStruct)
	assert.NotNil(t, err)
	assert.Nil(t, r)

	chStop := make(chan bool)
	msgs, err := rmq.Subscribe(subStruct)
	subCounter := 1
	go func() {
		for msg := range msgs {
			log.Printf("Subscribed a message: %s", msg.Body)
			result, err := rmq.Acknowledge(msg.DeliveryTag)
			assert.Nil(t, err)
			assert.Equal(t, true, result)
			if subCounter == (messagesCount) {
				chStop <- true
			}
			subCounter++
		}
	}()
	result := <-chStop
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}
