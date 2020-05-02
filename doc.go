// package gomq provides methods for subscription & publishing to different topics.
//
// GoMQ is thread safe. It provides Poll based functionality to read data from topic.
// A subscriber can subscribe to any numbers of topics through pattern matching.
//
// GoMQ works very similar to RabbitMQ. Every subscription creates a queue, from which data can be polled.
// The broker makes sure that the data pushed a topic will be sent to appropriate queues.
package gomq
