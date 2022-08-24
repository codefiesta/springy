package events

type Topic int

const (

	// Mongo Events
	Mongo Topic = iota

	// Websocket Events
	Websocket
)
