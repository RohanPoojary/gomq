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

    # Subscribes to "users" topic.
    usersPoller := broker.Subscribe(gomq.ExactMatcher("users"))

```

Regex Match
```go script
    
    broker := gomq.NewBroker()

    # Subscribes to "users.*" topic.
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