package main

import (
	"time"
)

type Publisher interface {
	Configure() error
	Publish(bindingKey string, json string, timestamp time.Time) error
}
