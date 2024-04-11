package consumer

import (
	"encoding/json"
	"fmt"
)

type Header struct {
	Key   string
	Value []byte
}

type Message struct {
	MsgBody string
	Headers []Header
}

func (m *Message) String() string {
	res, err := json.Marshal(m)
	if err != nil {
		newErr := fmt.Errorf("message json.Marshal error: %s", err.Error())
		panic(newErr.Error())
	}
	return string(res)
}

type InStream <-chan interface{}

type Consumer interface {
	Run()
	Stop()
	ConsumerStream() InStream
}
