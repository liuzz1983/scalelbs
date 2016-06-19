package main

import (
	json2 "encoding/json"
)

// Ping is a request
type Ping struct {
	Key string `json:"key"`
}

// Bytes returns the byets for a ping
func (p Ping) Bytes() []byte {
	data, _ := json2.Marshal(p)
	return data
}

// Pong is a ping response
type Pong struct {
	Message string `json:"message"`
	From    string `json:"from"`
}
