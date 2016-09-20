package bot

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

const (
	listenAddr      = "127.0.0.1:5000"
	defaultHostname = "localhost"
)

type Info struct {
	Hostname  string    `json:"hostname"`
	Pid       int       `json:"pid"`
	Version   string    `json:"version"`
	StartTime time.Time `json:"start_time"`
	Uptime    string    `json:"uptime"`
}

func (w *WeburgBot) CheckHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	hostname, err := os.Hostname()
	if err != nil {
		hostname = defaultHostname
	}

	var info = Info{
		Hostname:  hostname,
		Pid:       os.Getpid(),
		Version:   Version(),
		StartTime: w.StartTime,
		Uptime:    time.Now().Sub(w.StartTime).String(),
	}

	encoder := json.NewEncoder(rw)
	encoder.Encode(info)
}

func (w *WeburgBot) Listen() {
	r := mux.NewRouter()

	r.HandleFunc("/", w.CheckHandler).Methods("GET")
	http.Handle("/", r)

	logrus.Info(listenAddr)
	logrus.Fatal(http.ListenAndServe(listenAddr, nil))
}
