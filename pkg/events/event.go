package events

type Event struct {
	Topic  Topic
	Sender interface{}
	Data   interface{}
}

// Channel is a channel which can accept an Event
type Channel chan Event

