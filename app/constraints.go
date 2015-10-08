package app

import (
	"encoding/json"
	"log"
	"net/http"
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
