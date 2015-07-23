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

	resultBytes, err := formatJsonForInflux(msg.Payload)
	if err != nil {
		return err
	}
	log.Printf("POSTING FORMATTED DATA: %s", string(resultBytes))
	req, err := http.NewRequest("POST", influxUrl, bytes.NewBuffer(resultBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
}

func formatJsonForInflux(original []byte) ([]byte, error) {
	var tableName string
	payload := make(map[string]map[string]interface{})
	err := json.Unmarshal(original, &payload)

	keys := make([]string, 0)
	vals := make([]interface{}, 0)

	sobj, sok := payload["stat"]
	robj, rok := payload["rxpk"]
	if sok {
		tableName = "stat"
		for key, val := range sobj {
			keys = append(keys, key)
			vals = append(vals, val)
		}
	} else if rok {
		tableName = "rxpk"
		for key, val := range robj {
			keys = append(keys, key)
			vals = append(vals, val)
		}
	} else {
		return nil, errors.New("Unknown key")
	}

	result := []interface{}{
		map[string]interface{}{
			"name":    tableName,
			"columns": keys,
			"points":  [][]interface{}{vals},
		},
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return resultBytes, nil
}
