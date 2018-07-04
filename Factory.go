package queuelib

//Init : Queue factory - Returns a Queue for general purpose
func Init(QueueType string) (Queue, error) {
	switch QueueType {
	case RABBITMQ:
		return new(RabbitMQ), nil
	default:
		return nil, ErrorInvalidQueue
	}
}
