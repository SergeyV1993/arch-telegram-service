package models

type Update struct {
	ID      int
	Message *IncomingMessage
}

type IncomingMessage struct {
	Type MessageType

	Text     string
	Username string
	ChatID   int64

	Radius int

	Paginator *Paginator

	Location *Location
}

type Paginator struct {
	MaxPages    int
	CurrentPage int
	Offset      int
	Limit       int
}

type Location struct {
	Longitude string
	Latitude  string
}
