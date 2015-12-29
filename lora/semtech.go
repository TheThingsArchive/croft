package lora

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	PUSH_DATA = iota
	PUSH_ACK  = iota
	PULL_DATA = iota
	PULL_RESP = iota
	PULL_ACK  = iota
)

type Conn struct {
	Raw *net.UDPConn
}

type Message struct {
	SourceAddr *net.UDPAddr
	Conn       *Conn
	Header     *MessageHeader
	GatewayEui []byte
	Payload    interface{}
}

type MessageHeader struct {
	ProtocolVersion byte
	Token           uint16
	Identifier      byte
}

type PushMessagePayload struct {
	RXPK []*RXPK `json:"rxpk,omitempty"`
	Stat *Stat   `json:"stat,omitempty"`
}

type Stat struct {
	Time string  `json:"time"`
	Lati float64 `json:"lati"`
	Long float64 `json:"long"`
	Alti float64 `json:"alti"`
	Rxnb uint    `json:"rxnb"`
	Rxok uint    `json:"rxok"`
	Rxfw uint    `json:"rxfw"`
	Ackr float64 `json:"ackr"`
	Dwnb uint    `json:"dwnb"`
	Txnb uint    `json:"txnb"`
}

type RXPK struct {
	Time time.Time `json:"time"`
	Tmst uint      `json:"tmst"`
	Chan uint      `json:"chan"`
	Rfch uint      `json:"rfch"`
	Freq float64   `json:"freq"`
	Stat int       `json:"stat"`
	Modu string    `json:"modu"`
	Datr string    `json:"datr"`
	Codr string    `json:"codr"`
	Rssi int       `json:"rssi"`
	Lsnr float64   `json:"lsnr"`
	Size uint      `json:"size"`
	Data string    `json:"data"`
}

type TXPX struct {
	Imme bool    `json:"imme"`
	Freq float64 `json:"freq"`
	Rfch uint    `json:"rfch"`
	Powe uint    `json:"powe"`
	Modu string  `json:"modu"`
	Datr uint    `json:"datr"`
	Fdev uint    `json:"fdev"`
	Size uint    `json:"size"`
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
		if n < 12 {
			return nil, errors.New("Parse message failed, invalid size")
		}
		msg.GatewayEui = b[4:12]
		var payload PushMessagePayload
		err := json.Unmarshal(b[12:n], &payload)
		if err != nil {
			log.Printf("Parse message failed: %s\nMessage: %s", err.Error(), string(b[12:n]))
			return nil, err
		}
		msg.Payload = payload
	}
	return msg, nil
}

func (m *Message) Ack() error {
	var id byte
	switch m.Header.Identifier {
	case PUSH_DATA:
		id = PUSH_ACK
	case PULL_DATA:
		id = PULL_ACK
	default:
		return fmt.Errorf("Unknown message identifier %d to acknowledge", m.Header.Identifier)
	}

	ack := &MessageHeader{
		ProtocolVersion: m.Header.ProtocolVersion,
		Token:           m.Header.Token,
		Identifier:      id,
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

func (rxpk *RXPK) ParseData() (*PHYPayload, error) {
	buf, err := base64.StdEncoding.DecodeString(rxpk.Data)
	if err != nil {
		log.Printf("Failed to decode base64 data: %s", err.Error())
		return nil, err
	}

	return ParsePHYPayload(buf)
}
