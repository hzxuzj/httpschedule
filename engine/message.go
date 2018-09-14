package engine

import (
	"encoding/json"
)

type message struct {
	Msg string `json:"msg,omitempty"`
}

func newMessage(msg string) *message {
	return &message{
		Msg: msg,
	}
}

func (m *message) toJsonBytes() []byte {

	bytes, _ := json.Marshal(m)

	return bytes

}
