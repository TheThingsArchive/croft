package lora

import (
	"encoding/json"
	"net"
)

const (
	PUSH_DATA = iota
	PUSH_ACK  = iota
	PULL_DATA = iota
	PULL_ACK  = iota
	PULL_RESP = iota
)

type Message struct {
	ProtocolVersion int
	Token           []byte
	Identifier      int
	Payload         *json.RawMessage
	GatewayEUI      int64
}

type Conn struct {
	Raw *net.UDPConn
}

func NewConn(r *net.UDPConn) *Conn {
	return &Conn{r}
}
