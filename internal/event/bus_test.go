package event_test

import (
	"github.com/stretchr/testify/assert"
	"go.springy.io/pkg/event"
	"sync"
	"testing"
)

type TestSender struct {
	index int
}

type TestMessage struct {
	name  string
	count int
}

var messages = []TestMessage{

	{
		name:  "Foo",
		count: 2,
	},
	{
		name:  "Bar",
		count: 4,
	},
}

func TestBus(t *testing.T) {

	var wg sync.WaitGroup

	s := make(chan event.Event)
	event.Subscribe(event.Mongo, s)
	go subscribe(t, s, &wg)

	for i, m := range messages {
		// Set up our sender
		sender := TestSender{
			index: i,
		}
		// Block until Done is called
		wg.Add(1)
		event.Publish(event.Mongo, sender, m)
	}
	wg.Wait()
}

func subscribe(t *testing.T, s chan event.Event, wg *sync.WaitGroup) {
	for {
		select {
		case e := <-s:
			if sender, ok := e.Sender.(TestSender); ok {
				if data, ok := e.Data.(TestMessage); ok {
					t.Log("Received event: ", e)
					msg := messages[sender.index]
					assert.Equal(t, data.count, msg.count)
					assert.Equal(t, data.name, msg.name)
					wg.Done()
				}
			}
		}
	}

}
