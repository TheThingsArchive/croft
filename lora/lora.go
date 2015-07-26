package lora

import (
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

var buf = make([]byte, 2048)

type Message struct {
	ProtocolVersion int
	Token           []byte
	Identifier      int
	SourceAddr      *net.UDPAddr
	Conn            *Conn
	Payload         []byte
	GatewayEUI      net.HardwareAddr
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
	msg := &Message{
		SourceAddr:      addr,
		Conn:            c,
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
	p := make([]byte, 4)
	p[0] = byte(m.ProtocolVersion)
	p[1] = m.Token[0]
	p[2] = m.Token[1]
	p[3] = PUSH_ACK
	_, err := m.Conn.Raw.WriteToUDP(p, m.SourceAddr)
	if err != nil {
		return err
	}
	return nil
}
