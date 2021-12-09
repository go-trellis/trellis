package message

//// Message is the interface for publishing asynchronously
//type Message interface {
//	Service() *service.Service
//	Topic() string
//	SetTopic(string)
//	SetBody(v interface{}) error
//	GetPayload() *Payload
//	ToObject(v interface{}) error
//}
//
//type message struct {
//}

func (m *Response) SetPayload(payload *Payload) {
	m.Payload = payload
}
