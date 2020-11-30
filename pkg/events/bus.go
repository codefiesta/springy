package events

import (
	"log"
	"sync"
)

var (
	bus  *Bus
	once sync.Once
)

// See: https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997
type Bus struct {
	// Registered clients.
	subscribers map[Topic][]Channel
	mutex       sync.RWMutex
}

func init() {
	log.Println("Initializing Event Bus")
	once.Do(func() {
		bus = &Bus{
			subscribers: make(map[Topic][]Channel),
		}
	})
}

// Subscribes to events
func Subscribe(topic Topic, c Channel) {
	bus.mutex.Lock()
	if channels, found := bus.subscribers[topic]; found {
		bus.subscribers[topic] = append(channels, c)
	} else {
		bus.subscribers[topic] = append([]Channel{}, c)
	}
	bus.mutex.Unlock()
}

// Publishes events
func Publish(topic Topic, sender interface{}, data interface{}) {
	bus.mutex.RLock()
	if chs, found := bus.subscribers[topic]; found {
		// Create a new slice to preserve locking
		channels := append([]Channel{}, chs...)
		// Send an event to subscribed channels
		go func(e Event, channels []Channel) {
			for _, value := range channels {
				value <- e
			}
		}(Event{Topic: topic, Sender: sender, Data: data}, channels)
	}
	bus.mutex.RUnlock()
}
