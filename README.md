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

# Installation
``` 
    go get github.com/RohanPoojary/gomq
```

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

# Benchmark
The benchmark is done on Publish (publishes the message) & Poll (polls the message from queue).
In all benchmarks both publisher & subscriber are present at any point of time. 

Point to note is that the broker is optimized around Poll, Hence publish slows with more number of
Consumers, while the Poll doesn't slow at the same rate.

```shell script
>  go test -cpu 4 -bench .
goos: darwin
goarch: amd64
pkg: github.com/RohanPoojary/gomq
BenchmarkPublish/Consumers=1-4         	                     1000000	      1178 ns/op
BenchmarkPublish/Consumers=10-4        	                      300000	      4504 ns/op
BenchmarkPublish/Consumers=30-4        	                      200000	     11523 ns/op
BenchmarkPublish/Consumers=50-4        	                      100000	     19384 ns/op
BenchmarkConsume/Publisher=1-4         	                     1000000	      1599 ns/op
BenchmarkConsume/Publisher=10-4        	                     1000000	      1788 ns/op
BenchmarkConsume/Publisher=30-4        	                     1000000	      2182 ns/op
BenchmarkConsume/Publisher=50-4        	                     1000000	      2330 ns/op
BenchmarkMultiConsume/Publisher=50/Consumer=10-4         	 1000000	      1134 ns/op
BenchmarkMultiConsume/Publisher=50/Consumer=30-4         	  500000	      3568 ns/op
BenchmarkMultiConsume/Publisher=50/Consumer=50-4         	  500000	      4416 ns/op
PASS
```

# Contribution
Feel free to contribute or request any feature. Reporting issues can also help.
