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
	Key     string `json:"key"`
}

type Pos struct {
	Lat float64 `json:lat`
	Lng float64 `json:lng`
	Id  string  `json:id`
}

func (p *Pos) Bytes() []byte {
	data, _ := json2.Marshal(p)
	return data
}

type RangeQuery struct {
	Lat      float64 `json:lat`
	Lng      float64 `json:lng`
	Distance float64 `json:distance`
}

type CellQuery struct {
	CellId string `json:cellid`
}

func (p *CellQuery) Bytes() []byte {
	data, _ := json2.Marshal(p)
	return data
}

type QueryResult struct {
	Points []*Pos `json:points`
}
