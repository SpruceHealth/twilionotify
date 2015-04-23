package main

import (
	"encoding/xml"
	"flag"
	"log"
	"net/http"
	"strings"
)

var (
	flagNumbers = flag.String("n", "", "Phone numbers to forward SMS to")
	flagListen  = flag.String("l", "127.0.0.1:8086", "Address:port on which to listen")
)

type Response struct {
	XMLName  xml.Name   `xml:"Response"`
	Messages []*Message `xml:"Message"`
}

type Message struct {
	To   string `xml:"to,attr,omitempty"`
	From string `xml:"from,attr,omitempty"`
	Body string
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	var numbers []string
	for _, s := range strings.Split(*flagNumbers, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			numbers = append(numbers, s)
		}
	}
	if len(numbers) == 0 {
		log.Fatal("No numbers set")
	}

	http.HandleFunc("/twilio/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		body := r.FormValue("Body")
		if body == "" {
			log.Printf("Missing body")
			http.NotFound(w, r)
			return
		}
		res := &Response{}
		for _, n := range numbers {
			res.Messages = append(res.Messages, &Message{
				To:   n,
				Body: body,
			})
		}
		w.Header().Set("Content-Type", "application/xml")
		if err := xml.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Failed to marshal XML: %s", err)
		}
	})

	log.Fatal(http.ListenAndServe(*flagListen, nil))
}
