package queuelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//It will test all scenarios of Init()
func TestInit(t *testing.T) {
	q, err := Init("")
	assert.Equal(t, err, ErrorInvalidQueue)
	assert.Nil(t, q)

	queue, err := Init("rabbitmq")
	assert.Nil(t, err)
	assert.NotNil(t, queue)
}
