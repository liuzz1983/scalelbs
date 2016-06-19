package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/uber/ringpop-go"
	"github.com/uber/tchannel-go"

	json2 "encoding/json"
	"github.com/uber/tchannel-go/json"
	"golang.org/x/net/context"
)

type worker struct {
	ringpop *ringpop.Ringpop
	channel *tchannel.Channel
	logger  *log.Logger
}

func (w *worker) RegisterPong() error {
	hmap := map[string]interface{}{"/ping": w.PingHandler}

	return json.Register(w.channel, hmap, func(ctx context.Context, err error) {
		w.logger.Debug("error occured: %v", err)
	})
}

func (w *worker) PingHandler(ctx json.Context, ping *Ping) (*Pong, error) {
	var pong Pong
	var res []byte

	handle, err := w.ringpop.HandleOrForward(ping.Key, ping.Bytes(), &res, "ping", "/ping", tchannel.JSON, nil)
	if handle {
		identity, err := w.ringpop.WhoAmI()
		if err != nil {
			return nil, err
		}
		return &Pong{"Hello, world!", identity}, nil
	}

	if err := json2.Unmarshal(res, &pong); err != nil {
		return nil, err
	}

	// else request was forwarded
	return &pong, err
}
