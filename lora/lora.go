package lora

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jacobsa/crypto/cmac"
	"log"
)

type PHYPayload struct {
	MHDR       byte
	MACPayload []byte
	MIC        []byte
	FCtrl      byte
	FCnt       uint16
	FOpts      []byte
	FPort      byte
	DevAddr    uint32
}

func ParsePHYPayload(buf []byte) (*PHYPayload, error) {
	data := &PHYPayload{
		MHDR: buf[0],
	}

	//mType := (data.MHDR & 0xe0) >> 5

	majorVersion := data.MHDR & 0x3
	if majorVersion != 0 {
		return nil, errors.New(fmt.Sprintf("Major version %d not supported", majorVersion))
	}

	if len(buf) < 5 {
		return nil, errors.New("The data is should be at least 5 bytes")
	}

	data.MACPayload = buf[1 : len(buf)-4]
	data.MIC = buf[len(buf)-4 : len(buf)]

	if len(data.MACPayload) < 7 {
		return nil, errors.New("Payload should at least be 7 bytes")
	}

	binary.Read(bytes.NewReader(data.MACPayload[0:4]), binary.LittleEndian, &data.DevAddr)
	data.FCtrl = data.MACPayload[4]
	binary.Read(bytes.NewReader(data.MACPayload[5:7]), binary.LittleEndian, &data.FCnt)

	index := 7
	fOptsLen := int(data.FCtrl & 0xf)
	if fOptsLen > 0 {
		if len(data.MACPayload) < index+fOptsLen {
			return nil, errors.New("Payload does not contain indicated options length")
		}
		data.FOpts = data.MACPayload[index : index+fOptsLen]
		index += fOptsLen
	} else {
		data.FOpts = make([]byte, 0)
	}

	if len(data.MACPayload) > index+1 {
		data.FPort = data.MACPayload[index]
	}

	return data, nil
}

func (d *PHYPayload) DecryptPayload(key []byte) ([]byte, error) {
	if len(d.MACPayload) <= 8+len(d.FOpts) {
		return nil, errors.New("No data to decrypt")
	}
	data := d.MACPayload[8+len(d.FOpts):]

	// See LoRaWAN specification 1r0 4.3.3.1
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Failed to create AES cipher: %s", err.Error())
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	i := 1
	group := buf.Next(block.BlockSize())
	result := make([]byte, 0, len(data))
	for len(group) > 0 {
		a := new(bytes.Buffer)
		a.Write([]byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0})
		binary.Write(a, binary.LittleEndian, d.DevAddr)
		binary.Write(a, binary.LittleEndian, uint32(d.FCnt))
		a.WriteByte(0x0)
		a.WriteByte(byte(i))

		key := make([]byte, block.BlockSize())
		block.Encrypt(key, a.Bytes())

		for j := 0; j < len(group); j++ {
			result = append(result, group[j]^key[j])
		}

		i++
		group = buf.Next(block.BlockSize())
	}
	return result, nil
}

func (d *PHYPayload) TestIntegrity(key []byte) (bool, error) {
	// See LoRaWAN specification 1r0 4.4
	b0 := new(bytes.Buffer)
	b0.Write([]byte{0x49, 0x0, 0x0, 0x0, 0x0})
	b0.WriteByte(0x0)
	binary.Write(b0, binary.LittleEndian, d.DevAddr)
	binary.Write(b0, binary.LittleEndian, uint32(d.FCnt))
	b0.WriteByte(0x0)
	b0.WriteByte(byte(1 + len(d.MACPayload)))
	b0.WriteByte(d.MHDR)
	b0.Write(d.MACPayload)

	hash, err := cmac.New(key)
	if err != nil {
		log.Printf("Failed to initialize CMAC: %s", err.Error())
		return false, err
	}

	_, err = hash.Write(b0.Bytes())
	if err != nil {
		log.Printf("Failed to hash data: %s", err.Error())
		return false, err
	}

	calculatedMIC := hash.Sum([]byte{})[0:4]
	return bytes.Equal(calculatedMIC, d.MIC), nil
}
