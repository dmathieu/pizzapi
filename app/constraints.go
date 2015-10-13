package app

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Constraint struct {
	Constraint string `json:"constraint"`
	Ip         string `json:"ip"`
}

var loadedConstraints []*Constraint

func findConstraint(ip string) (*Constraint, error) {
	for _, v := range loadedConstraints {
		if v.Ip == ip {
			return v, nil
		}
	}

	constraint := &Constraint{Ip: ip, Constraint: "none"}
	loadedConstraints = append(loadedConstraints, constraint)
	return constraint, nil
}

func applyConstraints(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		constraint, _ := findConstraint(r.RemoteAddr)

		log.Printf("count#constraints method=%s path=%s constraint=%s ip=%s", r.Method, r.URL.Path, constraint.Constraint, constraint.Ip)
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
	w.WriteHeader(503)
	response := &ErrorResponse{Id: "maintenance", Message: "API is temporarily unavailable."}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
	log.Printf("count#http.maintenance method=%s path=%s", r.Method, r.URL.Path)
}

func slow(fn http.HandlerFunc, w http.ResponseWriter, r *http.Request) {

	rand.Seed(time.Now().Unix())
	duration := rand.Intn(60-30) + 30

	log.Printf("count#http.slow method=%s path=%s duration=%d", r.Method, r.URL.Path, duration)
	time.Sleep(time.Duration(duration) * time.Second)
	fn(w, r)
}

func erroring(fn http.HandlerFunc, w http.ResponseWriter, r *http.Request) {

	rand.Seed(time.Now().Unix())
	randomizer := rand.Intn(10)

	if randomizer >= 7 {
		fn(w, r)
	} else {
		log.Printf("count#http.error method=%s path=%s", r.Method, r.URL.Path)

		w.WriteHeader(500)
		response := &ErrorResponse{Id: "error", Message: "An unknown error occured."}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}
}
