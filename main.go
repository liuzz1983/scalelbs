// Copyright (c) 2015 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"flag"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/uber-common/bark"
	"github.com/uber/ringpop-go"
	"github.com/uber/ringpop-go/discovery/jsonfile"
	"github.com/uber/ringpop-go/logging"
	"github.com/uber/ringpop-go/swim"
	"github.com/uber/tchannel-go"
)

var (
	hostport = flag.String("listen", "127.0.0.1:3000", "hostport to start ringpop on")
	httpport = flag.String("http", "127.0.0.1:8000", "hostport to start ringpop on")
	hostfile = flag.String("hosts", "./hosts.json", "path to hosts file")
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	flag.Parse()

	logger := log.StandardLogger()

	l := bark.NewLoggerFromLogrus(logger)

	ch, err := tchannel.NewChannel("geo", &tchannel.ChannelOptions{
		Logger: ProxyLogger{l},
	})
	if err != nil {
		log.Fatalf("channel did not create successfully: %v", err)
	}

	logOpts := ringpop.LogLevels(map[string]logging.Level{
		"damping":       logging.Debug,
		"dissemination": logging.Debug,
		"gossip":        logging.Debug,
		"join":          logging.Debug,
		"membership":    logging.Debug,
		"ring":          logging.Debug,
		"suspicion":     logging.Debug,
	})

	rp, err := ringpop.New("geo-app",
		ringpop.Channel(ch),
		ringpop.Identity(*hostport),
		ringpop.Logger(l),
		logOpts,
	)
	if err != nil {
		log.Fatalf("Unable to create Ringpop: %v", err)
	}

	options := &WorkerOptions{
		CellMeters:  2.0,
		ExpiredTime: time.Minute * 1,
	}

	worker := NewWorker(rp, ch, logger, *httpport,options)

	if err := worker.Register(); err != nil {
		log.Fatalf("could not register pong handler: %v", err)
	}

	if err := worker.channel.ListenAndServe(*hostport); err != nil {
		log.Fatalf("could not listen on given hostport: %v", err)
	}

	opts := new(swim.BootstrapOptions)
	opts.DiscoverProvider = jsonfile.New(*hostfile)

	if _, err := worker.ringpop.Bootstrap(opts); err != nil {
		log.Fatalf("ringpop bootstrap failed: %v", err)
	}

	httpServer := NewHttpServer(*httpport, worker)
	httpServer.Serv()

	select {}
}
