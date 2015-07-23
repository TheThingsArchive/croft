package lora

import (
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
	Payload         []byte
	GatewayEUI      net.HardwareAddr
	SourceAddr      *net.UDPAddr
}

type Stat struct {
	Time string  `json:time`
	Lati float64 `json:lati`
	Long float64 `json:long`
	Alti float64 `json:alti`
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

	if msg.Identifier == PUSH_DATA {
		msg.Payload = buf[12:n]
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
