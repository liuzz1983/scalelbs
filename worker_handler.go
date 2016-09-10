package main

import (
	"github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/json"
	"golang.org/x/net/context"
)

func (w *Worker) Register() error {
	hmap := map[string]interface{}{
		PingPath:  w.PingHandler,
		AddPath:   w.AddHandler,
		QueryPath: w.QueryHandler,
	}

	return json.Register(w.channel, hmap, func(ctx context.Context, err error) {
		w.logger.Errorf("error occured: %v", err)
	})
}

func (w *Worker) AddHandler(ctx json.Context, pos *Pos) (map[string]interface{}, error) {

	cell := w.geoIndexer.PosCell(pos)
	key := cell.Id()
	var res []byte

	handle, err := w.ringpop.HandleOrForward(key, pos.Bytes(), &res, ServiceName, AddPath, tchannel.JSON, nil)
	if handle {
		w.logger.WithField("pos", pos).Error(" begin to add point")
		if err != nil {
			return nil, err
		}
		w.geoIndexer.AddPos(pos)
		return nil, nil
	}

	return nil, nil
}

func (w *Worker) QueryHandler(ctx json.Context, query *CellQuery) (*QueryResult, error) {

	var res []byte
	handle, err := w.ringpop.HandleOrForward(query.CellId, query.Bytes(), &res, ServiceName, QueryPath, tchannel.JSON, nil)
	if handle {
		w.logger.WithField("query", query).Debug(" begin to query point")
		result := &QueryResult{
			Points: make([]*Pos, 0),
		}

		values := w.geoIndexer.Get(query.CellId)

		for _, v := range values {
			result.Points = append(result.Points, v)
		}

		return result, nil
	}

	return nil, err
}

func (w *Worker) PingHandler(ctx json.Context, ping *Ping) (*Pong, error) {
	var res []byte

	handle, err := w.ringpop.HandleOrForward(ping.Key, ping.Bytes(), &res, ServiceName, PingPath, tchannel.JSON, nil)
	if handle {
		identity, err := w.ringpop.WhoAmI()
		if err != nil {
			return nil, err
		}
		return &Pong{"Hello, world!", identity, ping.Key}, nil
	}
	return nil, err
}
