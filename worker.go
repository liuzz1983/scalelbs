package main

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/uber/ringpop-go"
	_ "github.com/uber/ringpop-go/replica"
	"github.com/uber/tchannel-go"
)

type WorkerOptions struct {
	CellMeters  float64
	ExpiredTime time.Duration
}

type Worker struct {
	ringpop  *ringpop.Ringpop
	channel  *tchannel.Channel
	logger   *log.Logger
	httpHost string

	geoIndexer *GeoIndexer
}

var DefaultOptions *WorkerOptions = &WorkerOptions{
	CellMeters:  2.0,
	ExpiredTime: time.Minute * 1,
}

func MergeOptions(options *WorkerOptions) *WorkerOptions {
	result := *DefaultOptions
	if options == nil {
		return &result
	}

	result.CellMeters = options.CellMeters
	result.ExpiredTime = options.ExpiredTime
	return &result
}

func NewWorker(ringpop *ringpop.Ringpop, channel *tchannel.Channel, logger *log.Logger, httpHost string, options *WorkerOptions) *Worker {

	values := MergeOptions(options)

	return &Worker{
		ringpop:    ringpop,
		channel:    channel,
		logger:     logger,
		httpHost:   httpHost,
		geoIndexer: NewGeoIndexer(Km(values.CellMeters),values.ExpiredTime),
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

type Reponses struct {
	successes []*QueryResult
	errors    []error
	sync.Mutex
}

func (res *Reponses) add(r *QueryResult, err error) {
	res.Lock()
	go res.Unlock()

	if r != nil {
		res.successes = append(res.successes, r)
	}

	if err != nil {
		res.errors = append(res.errors, err)
	}
}

func (w *Worker) Search(lat float64, lng float64) (*QueryResult, error) {

	cells := w.geoIndexer.Cells(lat, lng)
	r := &QueryResult{
		Points: make([]*Pos, 0),
	}

	//replicator := replica.NewReplicator(w.ringpop, w.tchannel.GetSubChannel(ServiceName), nil, nil)

	var wg sync.WaitGroup

	responses := &Reponses{}

	for _, cell := range cells {

		wg.Add(1)
		go func(c Cell) {

			defer wg.Done()

			dest, err := w.ringpop.Lookup(c.Id())
			if err != nil {
				w.logger.Errorf("error in query id:%v with error", c.Id(), err)
				responses.add(nil, err)
				return
			}

			query := CellQuery{
				CellId: c.Id(),
			}

			res, err := w.ringpop.Forward(dest, []string{c.Id()}, query.Bytes(), ServiceName, QueryPath, tchannel.JSON, nil)
			if err != nil {
				w.logger.Errorf("error in forward dest")
				responses.add(nil, err)
				return

			}

			result := &QueryResult{}
			if err := json.Unmarshal(res, result); err != nil {
				w.logger.Errorf("error in unmarshal result %v", err)
				responses.add(nil, err)
				return
			}

			responses.add(result, nil)

		}(cell)
	}

	wg.Wait()

	//TODO need to remove duplication
	for _, res := range responses.successes {
		if res.Points != nil {
			r.Points = append(r.Points, res.Points...)
		}
	}

	// TODO need return it to the
	for _, err := range responses.errors {
		if err != nil {
			w.logger.Errorf("error in process request for error %v", err)
		}
	}

	return r, nil
}
