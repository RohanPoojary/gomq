// package queue provides methods for queue handling.
//
// The queue is thread safe and unbounded. Hence data can be pushed without any reader.
// Reading from a queue is Poll based, and thus it's a blocking call until the queue is either non-empty or closed.
//
// Internally it creates 2 unbuffered channel and an array to co-ordinate between the two.
// The implementation is based on implementation by rgooch in https://github.com/golang/go/issues/20352#issue-228477118 .
package queue
