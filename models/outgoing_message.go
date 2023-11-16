package models

type OutgoingMessage struct {
	Type MessageType

	Text     string
	Username string
	ChatID   int64

	Radius int

	Paginator *Paginator

	Location *Location

	Buildings []Building
}

type Building struct {
	Name           string
	Address        string
	LinkMapAddress string // ссылка на карту

	Link        string // ссылка на описание здания
	Description string

	Distance float64
}
