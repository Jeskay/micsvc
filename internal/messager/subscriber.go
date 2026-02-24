package messager

type subscriber struct {
	msg       chan Message
	closeFunc func()
}

type Message struct {
	AuthorID string
	Binary bool
	Data []byte
}
