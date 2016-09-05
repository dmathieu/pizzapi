package pizzapi

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Constraint struct {
	Constraint string `json:"constraint"`
	Token      string `json:"token"`
}

var loadedConstraints []*Constraint

var ErrInvalidToken = errors.New("A valid token is requied in the `Authorization` header")

func buildConstraint(token, constraint string) (*Constraint, error) {
	c, err := findConstraint(token)
	if err == nil {
		return c, nil
	}

	c = &Constraint{Token: token, Constraint: constraint}
	loadedConstraints = append(loadedConstraints, c)
	return c, nil
}

func findConstraint(token string) (*Constraint, error) {
	for _, v := range loadedConstraints {
		if v.Token == token {
			return v, nil
		}
	}

	return nil, ErrInvalidToken
}

func applyConstraints(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		constraint, err := findConstraint(token)

		if r.URL.Path == "/upgrade" {
			fn(w, r)
			return
		}

		if err != nil {
			if err == ErrInvalidToken {
				w.WriteHeader(401)
			} else {
				w.WriteHeader(500)
			}
			response := &errorResponse{"error", err.Error()}

			if err := json.NewEncoder(w).Encode(response); err != nil {
				panic(err)
			}
			request_id := requestId(r)
			token := r.Header.Get("Authorization")
			log.Printf("count#http.error method=%s path=%s request_id=%s token=%s", r.Method, r.URL.Path, request_id, token)
			return
		}

		request_id := requestId(r)
		log.Printf("count#constraints method=%s path=%s constraint=%s request_id=%s token=%s", r.Method, r.URL.Path, constraint.Constraint, request_id, constraint.Token)
		switch constraint.Constraint {
		case "maintenance":
			maintenance(w, r)
			return
		case "slow":
			slow(fn, w, r)
			return
		case "erroring":
			erroring(fn, w, r)
			return
		}
		fn(w, r)
	}
}

func maintenance(w http.ResponseWriter, r *http.Request) {
	if err := serveResponse(w, 503, &errorResponse{"maintenance", "API is temporarily unavailable"}); err != nil {
		panic(err)
	}

	request_id := requestId(r)
	token := r.Header.Get("Authorization")
	log.Printf("count#http.maintenance method=%s path=%s request_id=%s token=%s", r.Method, r.URL.Path, request_id, token)
}

func slow(fn http.HandlerFunc, w http.ResponseWriter, r *http.Request) {

	rand.Seed(time.Now().Unix())
	duration := rand.Intn(30-15) + 15

	request_id := requestId(r)
	token := r.Header.Get("Authorization")
	log.Printf("count#http.slow method=%s path=%s duration=%d request_id=%s token=%s", r.Method, r.URL.Path, duration, request_id, token)
	time.Sleep(time.Duration(duration) * time.Second)
	fn(w, r)
}

func erroring(fn http.HandlerFunc, w http.ResponseWriter, r *http.Request) {

	rand.Seed(time.Now().Unix())
	randomizer := rand.Intn(10)

	if randomizer >= 7 {
		fn(w, r)
	} else {
		request_id := requestId(r)
		token := r.Header.Get("Authorization")
		log.Printf("count#http.error method=%s path=%s request_id=%s token=%s", r.Method, r.URL.Path, request_id, token)

		if err := serveResponse(w, 500, &errorResponse{"not_found", "An unknown error occured"}); err != nil {
			panic(err)
		}
	}
}
