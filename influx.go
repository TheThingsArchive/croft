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

func WriteData(msg lora.Message) error {
	log.Print("ATTEMPTING TO WRITE DATA")
	log.Printf("ORIGINAL JSON WAS: %s", string(msg.Payload))

	if msg.Payload == nil {
		return errors.New("No payload provided")
	}
	payload := make(map[string]interface{})
	err := json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		return err
	}

	keys := make([]string, 0)
	vals := make([]interface{}, 0)

	for key, val := range payload {
		keys = append(keys, key)
		vals = append(vals, val)
	}

	result := []interface{}{
		map[string]interface{}{
			"name":    "push",
			"columns": keys,
			"points":  [][]interface{}{vals},
		},
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	log.Printf("POSTING FORMATTED DATA: %s", string(resultBytes))

	url := fmt.Sprintf(
		"http://%s:%s/db/%s/series?u=%s&p=%s",
		os.Getenv("INFLUXDB_HOST"),
		os.Getenv("INFLUXDB_PORT"),
		os.Getenv("INFLUXDB_DB"),
		os.Getenv("INFLUXDB_USER"),
		os.Getenv("INFLUXDB_PWD"),
	)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(resultBytes))
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
