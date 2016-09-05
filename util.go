package pizzapi

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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
		request_id := requestId(r)
		token := r.Header.Get("Authorization")
		fn(w, r)
		log.Printf("count#http method=%s path=%s request_id=%s token=%s", r.Method, r.URL.Path, request_id, token)
	}
}

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid), nil
}

func requestId(r *http.Request) (id string) {
	if id = r.Header.Get("Request-Id"); id == "" {
		id = r.Header.Get("X-Request-Id")
	}

	if id == "" {
		// In the event of a rare case where uuid
		// generation fails, it's probably more
		// desirable to continue as is with an empty
		// request_id than to bubble the error up the
		// stack.
		uuid, _ := NewUUID()
		id = string(uuid)
		r.Header.Set("X-Request-Id", id)
	}

	return id
}

type errorResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func serveResponse(w http.ResponseWriter, status int, response interface{}) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}
