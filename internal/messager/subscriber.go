package messager

type subscriber struct {
	msg       chan Message
	closeSlow func() error
}

type Message struct {
	AuthorID string
	Binary bool
	Data []byte
}
