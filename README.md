# GoMQ
[![Test Cases](https://github.com/RohanPoojary/gomq/actions/workflows/run-tests.yml/badge.svg)](https://github.com/RohanPoojary/gomq/actions/workflows/run-tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/RohanPoojary/gomq/badge.svg?branch=master)](https://coveralls.io/github/RohanPoojary/gomq?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/RohanPoojary/gomq)](https://goreportcard.com/report/github.com/RohanPoojary/gomq)
[![GoDoc](https://godoc.org/github.com/RohanPoojary/gomq?status.svg)](https://godoc.org/github.com/RohanPoojary/gomq)

An In-memory message broker in Golang.

# Introduction
Gomq is an in-memory message broker. The behaviour is similar to rabbitmq. 
Hence provides subscription to topics based on pattern match.

Gomq being a message broker is thread safe and underlyingly uses custom implementation of unbounded queue.
Hence no blocking should be observed while writing.

Supports Go 1.15+.

# Installation
```shell script
go get github.com/RohanPoojary/gomq
```

# Features
- Support for pattern based topic subscription.
- Fast variants of broker around Publish & Poll functionalities.
- Interface based functionality, for easier testing.

# Topics

### Creating a Broker
`NewBroker` returns a broker optimised around Polling, you can use `NewAsyncBroker` which is
optimised around Publish. View [Benchmarks](#benchmarks) for more details.

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

    count := broker.Publish("users.id", 100)
    // `count` contains the number of subscribers to whom the data has been published

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

# Benchmarks
The benchmark is done on Publish (publishes the message) & Poll (polls the message from queue).
In all benchmarks both publisher & subscriber are present at any point of time. 

Point to note is that the `NewBroker` returns a broker which is optimized around Poll,
where as `NewAsyncBroker` return a broker, which is optimized for Publish.

```powershell script
❯ go test -bench . -benchmem -count=10  | out-file -encoding utf8 "benchmark.txt"
❯ benchstat .\benchmark.txt
goarch: amd64
pkg: github.com/RohanPoojary/gomq
cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
                                 │ .\benchmark.txt │
                                 │     sec/op      │
Publish/SyncBroker/Pollers=1-8        1.007µ ±  6%
Publish/SyncBroker/Pollers=10-8       3.336µ ±  3%
Publish/SyncBroker/Pollers=30-8       9.540µ ±  4%
Publish/SyncBroker/Pollers=50-8       14.95µ ±  4%
Publish/AsyncBroker/Pollers=1-8       1.219µ ±  2%
Publish/AsyncBroker/Pollers=10-8      1.239µ ±  3%
Publish/AsyncBroker/Pollers=30-8      1.223µ ±  1%
Publish/AsyncBroker/Pollers=50-8      1.214µ ±  4%
Poll/SyncBroker/Publisher=1-8         1.677µ ±  8%
Poll/SyncBroker/Publisher=10-8        1.741µ ±  6%
Poll/SyncBroker/Publisher=30-8        2.097µ ±  5%
Poll/SyncBroker/Publisher=50-8        2.160µ ±  6%
Poll/AsyncBroker/Publisher=1-8        2.891µ ±  8%
Poll/AsyncBroker/Publisher=10-8       4.335µ ± 51%
Poll/AsyncBroker/Publisher=30-8       5.166µ ± 61%
Poll/AsyncBroker/Publisher=50-8       5.405µ ±  8%
geomean                               2.621µ

                                 │ .\benchmark.txt │
                                 │      B/op       │
Publish/SyncBroker/Pollers=1-8         103.0 ±  4%
Publish/SyncBroker/Pollers=10-8        268.0 ±  0%
Publish/SyncBroker/Pollers=30-8        705.0 ±  0%
Publish/SyncBroker/Pollers=50-8      1.102Ki ±  1%
Publish/AsyncBroker/Pollers=1-8        144.0 ±  5%
Publish/AsyncBroker/Pollers=10-8       152.5 ± 13%
Publish/AsyncBroker/Pollers=30-8       170.5 ±  2%
Publish/AsyncBroker/Pollers=50-8       163.0 ±  3%
Poll/SyncBroker/Publisher=1-8          38.00 ±  3%
Poll/SyncBroker/Publisher=10-8         46.50 ± 96%
Poll/SyncBroker/Publisher=30-8         137.0 ±  6%
Poll/SyncBroker/Publisher=50-8         152.0 ±  5%
Poll/AsyncBroker/Publisher=1-8         306.0 ±  5%
Poll/AsyncBroker/Publisher=10-8        509.5 ± 56%
Poll/AsyncBroker/Publisher=30-8        608.0 ± 74%
Poll/AsyncBroker/Publisher=50-8        642.0 ±  7%
geomean                                220.0

                                 │ .\benchmark.txt │
                                 │    allocs/op    │
Publish/SyncBroker/Pollers=1-8         2.000 ±  0%
Publish/SyncBroker/Pollers=10-8        10.00 ±  0%
Publish/SyncBroker/Pollers=30-8        29.00 ±  0%
Publish/SyncBroker/Pollers=50-8        48.00 ±  0%
Publish/AsyncBroker/Pollers=1-8        4.000 ±  0%
Publish/AsyncBroker/Pollers=10-8       4.000 ± 25%
Publish/AsyncBroker/Pollers=30-8       5.000 ±  0%
Publish/AsyncBroker/Pollers=50-8       5.000 ±  0%
Poll/SyncBroker/Publisher=1-8          2.000 ±  0%
Poll/SyncBroker/Publisher=10-8         1.000 ±  0%
Poll/SyncBroker/Publisher=30-8         2.000 ±  0%
Poll/SyncBroker/Publisher=50-8         2.000 ±  0%
Poll/AsyncBroker/Publisher=1-8         6.000 ± 17%
Poll/AsyncBroker/Publisher=10-8        11.00 ± 45%
Poll/AsyncBroker/Publisher=30-8        13.00 ± 62%
Poll/AsyncBroker/Publisher=50-8        13.00 ±  8%
geomean                                5.621
```

# Contribution
Feel free to contribute or request any feature. Reporting issues can also help.
