package main

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/uber/ringpop-go"
	"github.com/uber/tchannel-go"

	"encoding/json"
)

type Worker struct {
	ringpop  *ringpop.Ringpop
	channel  *tchannel.Channel
	logger   *log.Logger
	httpHost string

	geoIndexer *GeoIndexer
}

func NewWorker(ringpop *ringpop.Ringpop, channel *tchannel.Channel, logger *log.Logger, httpHost string) *Worker {
	return &Worker{
		ringpop:    ringpop,
		channel:    channel,
		logger:     logger,
		httpHost:   httpHost,
		geoIndexer: NewGeoIndexer(Km(2.0)),
	}
}

func (w *Worker) Ping(ping *Ping) (*Pong, error) {

	if !w.ringpop.Ready() {
		w.logger.Errorf("rong is not ready")
		return nil, errors.New("worker not ready")
	}

	dest, err := w.ringpop.Lookup(ping.Key)
	if err != nil {
		w.logger.Errorf("can not find dest for %v with error %v", ping.Key, err)
		return nil, err
	}

	res, err := w.ringpop.Forward(dest, []string{ping.Key}, ping.Bytes(), ServiceName, PingPath, tchannel.JSON, nil)
	if err != nil {
		w.logger.Errorf("can not forward to dest for error %v", err)
		return nil, err
	}

	var pong Pong
	if err := json.Unmarshal(res, &pong); err != nil {
		w.logger.Errorf("error in unmarshal result %v", err)
		return nil, err
	}
	// else request was forwarded
	return &pong, err

}
func (w *Worker) Add(pos *Pos) error {

	cell := w.geoIndexer.Cell(pos.Lat, pos.Lng)
	key := cell.Id()

	dest, err := w.ringpop.Lookup(key)
	if err != nil {
		w.logger.Errorf("can not find dest %v for key %v", dest, key)
		return errors.New("cant not find dest")
	}
	_, err = w.ringpop.Forward(dest, []string{key}, pos.Bytes(), ServiceName, AddPath, tchannel.JSON, nil)
	return err
}

func (w *Worker) Search(lat float64, lng float64) (*QueryResult, error) {

	cells := w.geoIndexer.Cells(lat, lng)
	r := &QueryResult{
		Points: make([]*Pos, 0),
	}

	for _, cell := range cells {

		dest, err := w.ringpop.Lookup(cell.Id())
		if err != nil {
			w.logger.Errorf("error in query id:%v with error", cell.Id(), err)
			continue
		}
		query := CellQuery{
			CellId: cell.Id(),
		}

		res, err := w.ringpop.Forward(dest, []string{cell.Id()}, query.Bytes(), "geo", "/geo_query", tchannel.JSON, nil)
		if err != nil {
			continue
		}
		result := QueryResult{}
		if err := json.Unmarshal(res, &result); err != nil {
			w.logger.Errorf("error in unmarshal result %v", err)
			continue
		}
		if result.Points != nil {
			r.Points = append(r.Points, result.Points...)
		}
	}

	return r, nil
}
