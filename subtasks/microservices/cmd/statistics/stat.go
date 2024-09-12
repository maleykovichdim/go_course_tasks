package main

import (
	"context"
	"encoding/json"
	"log"
	k "maleykovich/internal/kafka_client"
	"maleykovich/package/api"
	"time"

	// lc "maleykovich/internal/link_converter"

	"net/http"
)

type Service struct {
	NumReference uint64  `json:"id"`
	AverLength   float64 `json:"average_length"`
	totalLength  uint64
	dp           api.Dispatcher
}

var s Service

func (s *Service) run() {
	for {
		msg, err := s.dp.ReadConfirmed(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		s.NumReference++
		s.totalLength += uint64(len(msg))
		s.AverLength = float64(s.totalLength) / float64(s.NumReference)
	}
}

func main() {

	for {
		kf, err := k.New(
			[]string{"kafka:29092"},
			"shortener-topic",
			"ms-group",
		)
		if err != nil {
			log.Println(err) //TODO: change to simple logging mechanism
			time.Sleep(1 * time.Second)
			continue
		}
		s.dp = kf
		break
	}

	go s.run()

	http.HandleFunc("/", GetStatictics)

	http.ListenAndServe(":8081", nil)
}

func GetStatictics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}
