package gomq

import (
	"fmt"
	"regexp"
	"time"
)

func Example_basic() {
	broker := NewBroker()
	defer broker.Close(0)

	usersPoller := broker.Subscribe(ExactMatcher("users"))

	type User struct {
		ID int
		Name string
	}

	go func () {
		// Publisher Routine.
		for i := 1; i < 10; i++ {
			broker.Publish("users", User{ID: i})
		}
	}()

	for i := 1; i <= 5; i++ {
		val, ok := usersPoller.Poll()
		fmt.Printf("User%+v, Ok: %v\n", val, ok)
	}

	// Output:
	// User{ID:1 Name:}, Ok: true
	// User{ID:2 Name:}, Ok: true
	// User{ID:3 Name:}, Ok: true
	// User{ID:4 Name:}, Ok: true
	// User{ID:5 Name:}, Ok: true
}

func Example_advanced() {
	broker := NewBroker()
	defer broker.Close(0)

	// Subscribe should be done before publishing.
	// Or else all the data before subscription will be lost.
	// Subscribe to users.* topics.
	usersPoller := broker.Subscribe(regexp.MustCompile(`users\..*`))

	type User struct {
		ID int
		Country string
	}

	go func () {
		// Publishing User Data from India.
		for i := 1; i < 10; i++ {
			broker.Publish("users.india", User{ID: i, Country: "India"})
		}
	}()

	go func () {
		// Publishing User Data from USA.
		for i := 1; i < 10; i++ {
			broker.Publish("users.usa", User{ID: i, Country: "USA"})
		}
	}()

	// Subscribe to log.* topics
	logPoller := broker.Subscribe(regexp.MustCompile(`log\..*`))

	go func() {
		for val, ok := usersPoller.Poll(); ok; val, ok = usersPoller.Poll() {
			// Processing on User Data.
			log := fmt.Sprintf("User%+v, Ok: %v\n", val, ok)
			broker.Publish("log.user", log)
		}
	}()

	for i := 1; i <= 5; i++ {
		val, ok := logPoller.Poll()
		fmt.Println(time.Now(), val, ok )
	}

}

func ExampleBroker_Subscribe_exact() {
	broker := NewBroker()

	// Subscribes to "all" topic.
	poller := broker.Subscribe(ExactMatcher("all"))
	for i := 1; i <= 5; i++ {
		broker.Publish("all", fmt.Sprintf("All-%d", i))
		broker.Publish("all2", fmt.Sprintf("All2-%d", i))
	}

	broker.Close(-1)

	for val, ok := poller.Poll(); ok; val, ok = poller.Poll() {
		fmt.Println(val)
	}

	// Output:
	// All-1
	// All-2
	// All-3
	// All-4
	// All-5

}

func ExampleBroker_Subscribe_regex() {
	broker := NewBroker()

	// Subscribes to "all" topic.
	poller := broker.Subscribe(regexp.MustCompile(`all-.*`))
	for i := 1; i <= 3; i++ {
		broker.Publish("all-1", fmt.Sprintf("All-%d", i))
		broker.Publish("all-2", fmt.Sprintf("All2-%d", i))
	}

	broker.Close(-1)

	for val, ok := poller.Poll(); ok; val, ok = poller.Poll() {
		fmt.Println(val)
	}

	// Output:
	// All-1
	// All2-1
	// All-2
	// All2-2
	// All-3
	// All2-3
}
