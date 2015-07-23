package lora

import (
	"encoding/json"
	"log"
	"net"
)

const (
	PUSH_DATA = iota
	PUSH_ACK  = iota
	PULL_DATA = iota
	PULL_ACK  = iota
	PULL_RESP = iota
)

var buf = make([]byte, 2048)

type Message struct {
	ProtocolVersion int
	Token           []byte
	Identifier      int
	Payload         *json.RawMessage
	GatewayEUI      net.HardwareAddr
	SourceAddr      *net.UDPAddr
}

type Conn struct {
	Raw *net.UDPConn
}

func NewConn(r *net.UDPConn) *Conn {
	return &Conn{r}
}

func (c *Conn) ReadMessage() (*Message, error) {
	n, addr, err := c.Raw.ReadFromUDP(buf)
	if err != nil {
		log.Print("Error: ", err)
		return nil, err
	}
	log.Print("Received raw bytes", buf[0:n], " from ", addr)
	log.Print("Received ", string(buf[0:n]), " from ", addr)
	msg := &Message{
		SourceAddr:      addr,
		ProtocolVersion: int(buf[0]),
		Token:           buf[1:3],
		Identifier:      int(buf[3]),
	}
	return msg, nil
}

func (m *Message) Ack() error {
	conn, err := net.DialUDP("udp", nil, m.SourceAddr)
	defer conn.Close()
	if err != nil {
		return err
	}
	p := make([]byte, 4)
	p[0] = byte(m.ProtocolVersion)
	p[1] = m.Token[0]
	p[2] = m.Token[1]
	p[3] = PUSH_ACK
	_, err = conn.Write(p)
	if err != nil {
		return err
	}
	return nil
}
