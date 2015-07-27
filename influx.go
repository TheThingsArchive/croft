package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thethingsnetwork/croft/lora"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var influxUrl = fmt.Sprintf(
	"http://%s:%s/db/%s/series?u=%s&p=%s",
	os.Getenv("INFLUXDB_HOST"),
	os.Getenv("INFLUXDB_PORT"),
	os.Getenv("INFLUXDB_DB"),
	os.Getenv("INFLUXDB_USER"),
	os.Getenv("INFLUXDB_PWD"),
)

func WriteData(msg lora.Message) error {
	log.Print("ATTEMPTING TO WRITE DATA")
	if msg.Payload == nil {
		return errors.New("No payload provided")
	}
	log.Printf("ORIGINAL JSON WAS: %s", string(msg.Payload))

	resultBytes, err := FormatJsonBytesForInflux(msg.Payload)
	if err != nil {
		return err
	}
	log.Printf("POSTING FORMATTED DATA: %s", string(resultBytes))
	req, err := http.NewRequest("POST", influxUrl, bytes.NewBuffer(resultBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
}

func FormatJsonBytesForInflux(original []byte) ([]byte, error) {
	payload := make(map[string]map[string]interface{})
	err := json.Unmarshal(original, &payload)
	if err != nil {
		return nil, err
	}
	_, sok := payload["stat"]
	_, rok := payload["rxpk"]
	if sok {
		return formatStat(payload)
	} else if rok {
		log.Print("TODO: formatRXPK")
		return nil, errors.New("formatRXPK not handled")
	} else {
		return nil, errors.New("Unknown key")
	}
}

func formatStat(payload map[string]map[string]interface{}) ([]byte, error) {
	vals := make([][]interface{}, 1)

	keys := []string{
		"timex",
		"lati",
		"long",
		"alti",
	}

	s := payload["stat"]

	vals[0] = make([]interface{}, 4)
	vals[0][0] = s["time"]
	vals[0][1] = s["lati"]
	vals[0][2] = s["long"]
	vals[0][3] = s["alti"]

	byts, err := json.Marshal([]interface{}{
		map[string]interface{}{
			"name":    "stat",
			"columns": keys,
			"points":  vals,
		},
	})

	if err != nil {
		return nil, err
	}

	return byts, nil
}

func formatRXPK(payload map[string]map[string]interface{}) []interface{} {
	return nil
}
