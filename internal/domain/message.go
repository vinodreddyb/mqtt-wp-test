package domain

type Message struct {
	MsgId   string
	Topic   string
	Payload []byte
}
