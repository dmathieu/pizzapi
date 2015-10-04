package app

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	HttpPort         = flag.String("httpPort", os.Getenv("PORT"), "HTTP port for the server.")
	HttpReadTimeout  = flag.Duration("httpReadTimeout", time.Hour, "Timeout for HTTP request reading")
	HttpWriteTimeout = flag.Duration("httpWriteTimeout", time.Hour, "Timeout for HTTP request writing")
)

func AwaitSignals(signals ...os.Signal) <-chan struct{} {
	s := make(chan os.Signal, 1)
	signal.Notify(s, signals...)
	log.Printf("signals.await signals=%v\n", signals)

	received := make(chan struct{})
	go func() {
		log.Printf("signals.received signal=%v\n", <-s)
		close(received)
	}()

	return received
}

func logRequest(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
		log.Printf("count#http method=%s path=%s", r.Method, r.URL.Path)

	}
}
