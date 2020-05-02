# GoMQ
[![Build Status](https://travis-ci.org/RohanPoojary/gomq.svg?branch=master)](https://travis-ci.org/RohanPoojary/gomq)
[![Coverage Status](https://coveralls.io/repos/github/RohanPoojary/gomq/badge.svg?branch=master)](https://coveralls.io/github/RohanPoojary/gomq?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/RohanPoojary/gomq)](https://goreportcard.com/report/github.com/RohanPoojary/gomq)
[![GoDoc](https://godoc.org/github.com/RohanPoojary/gomq?status.svg)](https://godoc.org/github.com/RohanPoojary/gomq)

An In-memory message broker in Golang.

# Introduction
Gomq is an in-memory message broker. The behaviour is similar to rabbitmq. 
Hence provides subscription to topics based on pattern match.

Gomq being a message broker it is thread safe and is Interface based. 
Thereby allowing easy testing.

# TODO
- [ ] Optimise Publish.
- [ ] Benchmarking of different functionalities.

# Topics

### Creating a Broker
```go script

import (
	"github.com/RohanPoojary/gomq"
)

func main() {
	broker := gomq.NewBroker()
        ...
}

```

### Subscribing to a Topic
You can subscribe to a topic both on exact & regex matcher.

Exact Match
```go script
    
    broker := gomq.NewBroker()

    // Subscribes to "users" topic.
    usersPoller := broker.Subscribe(gomq.ExactMatcher("users"))

```

Regex Match
```go script
    
    broker := gomq.NewBroker()

    // Subscribes to any topic that matches "users.*" .
    usersPoller := broker.Subscribe(regexp.MustCompile(`users\.\w*`))

```

### Publishing to a Topic

```go script
    
    broker := gomq.NewBroker()

    broker.Publish("users.id", 100)

```

### Reading from a Subscriber

```go script
    
    broker := gomq.NewBroker()
    defer broker.Close(0)

    usersPoller := broker.Subscribe(regexp.MustCompile(`users\.\w*`))
    
    go func() {
        for {
            // Poll is a blocking call. Hence wrapped into a different routine.
            value, ok := usersPoller.Poll()
            if !ok {
                return
            }
    
            userID := value.(int)
            // Do Something
        }
    }()

```