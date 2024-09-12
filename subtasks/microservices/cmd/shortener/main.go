package main

import (
	"log"
	k "maleykovich/internal/kafka_client"
	lc "maleykovich/internal/link_converter"

	"net/http"
)

func main() {

	kf, err := k.New(
		[]string{"kafka:29092"},
		"shortener-topic",
		"ms-group",
	)

	if err != nil {
		log.Fatal(err) //TODO: change to simple logging mechanism, o wait kafka
	}
	s := lc.New(kf)

	http.HandleFunc("/add", s.AddLink)
	http.HandleFunc("/get", s.GetOriginalLink)

	http.ListenAndServe(":8080", nil)
}
