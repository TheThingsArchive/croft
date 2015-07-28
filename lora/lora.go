package lora

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
)

const (
	PUSH_DATA = iota
	PUSH_ACK  = iota
	PULL_DATA = iota
	PULL_ACK  = iota
	PULL_RESP = iota
)

type Conn struct {
	Raw *net.UDPConn
}

type Message struct {
	SourceAddr *net.UDPAddr
	Conn       *Conn
	Header     *MessageHeader
	Payload    *Payload
}

type MessageHeader struct {
	ProtocolVersion byte
	Token           uint16
	Identifier      byte
}

type Stat struct {
	Time string  `json:"time"`
	Lati float64 `json:"lati"`
	Long float64 `json:"long"`
	Alti float64 `json:"alti"`
	Rxnb int     `json:"rxnb"`
	Rxok int     `json:"rxok"`
	Rxfw int     `json:"rxfw"`
	Ackr float64 `json:"ackr"`
	Dwnb int     `json:"dwnb"`
	Txnb int     `json:"txnb"`
}

type RXPK struct {
	Time time.Time `json:"time"`
	Tmst int       `json:"tmst"`
	Chan int       `json:"chan"`
	Rfch int       `json:"rfch"`
	Freq float64   `json:"freq"`
	Stat int       `json:"stat"`
	Modu string    `json:"modu"`
	Datr string    `json:"datr"`
	Codr string    `json:"codr"`
	Rssi int       `json:"rssi"`
	Lsnr float64   `json:"lsnr"`
	Size int       `json:"size"`
	Data string    `json:"data"`
}

type TXPX struct {
	Imme bool    `json:"imme"`
	Freq float64 `json:"freq"`
	Rfch int     `json:"rfch"`
	Powe int     `json:"powe"`
	Modu string  `json:"modu"`
	Datr int     `json:"datr"`
	Fdev int     `json:"fdev"`
	Size int     `json:"size"`
	Data string  `json:"data"`
}

func NewConn(r *net.UDPConn) *Conn {
	return &Conn{r}
}

func (c *Conn) ReadMessage() (*Message, error) {
	buf := make([]byte, 2048)
	n, addr, err := c.Raw.ReadFromUDP(buf)
	if err != nil {
		log.Print("Error: ", err)
		return nil, err
	}
	return c.parseMessage(addr, buf, n)
}

func (c *Conn) parseMessage(addr *net.UDPAddr, b []byte, n int) (*Message, error) {
	var header MessageHeader
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}
	msg := &Message{
		SourceAddr: addr,
		Conn:       c,
		Header:     &header,
	}
	if header.Identifier == PUSH_DATA {
		payload, err := parsePayload(string(b[12:n]))
		if err != nil {
			return nil, err
		}
		msg.Payload = payload
	}
	return msg, nil
}

func (m *Message) Ack() error {
	ack := &MessageHeader{
		ProtocolVersion: m.Header.ProtocolVersion,
		Token:           m.Header.Token,
		Identifier:      PUSH_ACK,
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, ack)
	if err != nil {
		return err
	}

	_, err = m.Conn.Raw.WriteToUDP(buf.Bytes(), m.SourceAddr)
	if err != nil {
		return err
	}
	return nil
}
