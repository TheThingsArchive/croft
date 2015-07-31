package main

type Publisher interface {
	Configure() error
	Publish(message interface{}) error
}
