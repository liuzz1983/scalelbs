package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}

type HttpServer struct {
	host   string
	worker *Worker
}

func NewHttpServer(host string, worker *Worker) *HttpServer {
	return &HttpServer{
		host:   host,
		worker: worker,
	}
}

func (s *HttpServer) Serv() {

	http.HandleFunc("/hello", s.Ping)
	http.HandleFunc("/add", s.Add)
	http.HandleFunc("/search", s.Search)

	fmt.Println("begin to listen on " + s.host)
	err := http.ListenAndServe(s.host, nil)
	if err != nil {
		fmt.Errorf("channel did not create successfully: %v", err)
	}
}

func (s *HttpServer) Add(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		OutputJson(res, 0, "参数错误", nil)
		return
	}

	lat, err := strconv.ParseFloat(req.FormValue("lat"), 64)
	lng, err := strconv.ParseFloat(req.FormValue("lng"), 64)
	id := req.FormValue("id")

	pos := &Pos{
		Lat: lat,
		Lng: lng,
		Id:  id,
	}
	err = s.worker.Add(pos)
	if err != nil {
		OutputJson(res, 404, err.Error(), nil)
	} else {
		OutputJson(res, 0, "", "success")
	}
}

func (s *HttpServer) Search(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		OutputJson(res, 0, "参数错误", nil)
		return
	}

	lat, err := strconv.ParseFloat(req.FormValue("lat"), 64)
	lng, err := strconv.ParseFloat(req.FormValue("lng"), 64)

	result, err := s.worker.Search(lat, lng)
	if err != nil {
		OutputJson(res, 404, err.Error(), nil)
	}
	OutputJson(res, 0, "", result)
}

func (s *HttpServer) Ping(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		OutputJson(w, 0, "参数错误", nil)
		return
	}
	key := req.FormValue("key")
	ping := &Ping{
		Key: key,
	}

	pong, err := s.worker.Ping(ping)
	if err != nil {
		OutputJson(w, 404, err.Error(), nil)
		return
	}
	OutputJson(w, 0, "", pong)

}
