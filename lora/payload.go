package lora

import (
	"errors"
	"strings"
)

type PayloadObject struct {
	Key   string
	Value string
}
type Payload []*PayloadObject

//This is  a quick hack and does not support nested json in *other* keys
//as defined by the spec, it is simply to get round duplicate keys in the json
//this also only really takes into account edge cases
func parsePayload(input string) (*Payload, error) {
	p := new(Payload)
	p, err := getNextPayloadObject(p, input)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getNextPayloadObject(p *Payload, input string) (*Payload, error) {
	startKey := strings.Index(input, `"`)
	if startKey == -1 {
		return nil, errors.New("Invalid input")
	}
	input = input[startKey+1 : len(input)]

	closeKey := strings.Index(input, `"`)
	key := input[0:closeKey]

	startJSON := strings.Index(input, `{`)
	stopJSON := strings.Index(input, `}`)
	startARRAY := strings.Index(input, `[`)
	stopARRAY := strings.Index(input, `]`)
	var startVAL int
	var stopVAL int
	if startJSON > 0 {
		if startARRAY > 0 {
			if startJSON > startARRAY {
				startVAL = startARRAY
				stopVAL = stopARRAY
			} else {
				startVAL = startJSON
				stopVAL = stopJSON
			}
		} else {
			startVAL = startARRAY
			stopVAL = stopARRAY
		}
	}

	value := input[startVAL:(stopVAL + 1)]
	po := &PayloadObject{
		Key:   key,
		Value: value,
	}
	*p = append(*p, po)
	remaining := input[stopVAL+1 : len(input)]
	if remaining[0] != ',' {
		return p, nil
	}
	return getNextPayloadObject(p, remaining)
}
