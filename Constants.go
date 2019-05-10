package queuelib

import "errors"

const (
	//RABBITMQ : This is used in Factory
	RABBITMQ = "rabbitmq"
	KAFKA    = "kafka"
)

var (
	//ErrCursor : This is used when cursor is pointing to nil
	ErrCursor = errors.New("Invalid cursor")

	//ErrConnection : This is used when error occurs while creating connection
	ErrConnection = errors.New("Failed to connect to RabbitMQ, Kindly check RMQ Server is running / reachable and config data is correct")

	//ErrorInvalidQueue : This is used in Init(), where user passes unsupported queue type
	ErrorInvalidQueue = errors.New("Unsupported queue type")
)
