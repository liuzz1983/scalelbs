package main

import (
	json "encoding/json"
)

type ToBytes struct{}

func (p *ToBytes) Bytes() []byte {
	data, _ := json.Marshal(p)
	return data
}

// Ping is a request
type Ping struct {
	Key string `json:"key"`
}

func (p *Ping) Bytes() []byte {
	data, _ := json.Marshal(p)
	return data
}

// Pong is a ping response
type Pong struct {
	Message string `json:"message"`
	From    string `json:"from"`
	Key     string `json:"key"`
}

func (p *Pong) Bytes() []byte {
	data, _ := json.Marshal(p)
	return data
}

type Pos struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Id  string  `json:"id"`
}

func (p *Pos) Bytes() []byte {
	data, _ := json.Marshal(p)
	return data
}

type RangeQuery struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Distance float64 `json:"distance"`
}

func (p *RangeQuery) Bytes() []byte {
	data, _ := json.Marshal(p)
	return data
}

type CellQuery struct {
	CellId string `json:"cellid"`
}

func (p *CellQuery) Bytes() []byte {
	data, _ := json.Marshal(p)
	return data
}

type QueryResult struct {
	Points []*Pos `json:"points"`
}

func (p *QueryResult) Bytes() []byte {
	data, _ := json.Marshal(p)
	return data
}
