package messaging

type Message struct {
	Body []byte
}

type Handler func(msg *Message) error
