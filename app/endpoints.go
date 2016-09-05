package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func pizzasList(w http.ResponseWriter, r *http.Request) {
	list, _ := PizzaList()

	if err := json.NewEncoder(w).Encode(list); err != nil {
		panic(err)
	}
}

func findOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	order, err := FindOrder(id)
	if err != nil {
		w.WriteHeader(404)
		response := &errorResponse{Id: "not_found", Message: "Order not found"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(order); err != nil {
		panic(err)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var orderContent struct {
		Pizza int `json:"id"`
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &orderContent); err != nil {
		w.WriteHeader(422)
		response := &errorResponse{Id: "invalid", Message: fmt.Sprintf("Invalid data. Expected '{\"id\":1}'. Got %s", body)}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	pizza, err := FindPizza(orderContent.Pizza)
	if err != nil {
		w.WriteHeader(422)
		response := &errorResponse{Id: "not_found", Message: "Pizza not found"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	order := CreateOrder(pizza)
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		panic(err)
	}
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	var upgradeContent struct {
		Name string `json:"name"`
		Ip   string `json:"ip"`
		Key  string `json:"key"`
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &upgradeContent); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		response := &errorResponse{Id: "invalid", Message: "Invalid data"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	if upgradeContent.Key != os.Getenv("UPGRADE_KEY") {
		response := &errorResponse{Id: "unauthorized", Message: "Not authorized"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	var constraint *Constraint
	if upgradeContent.Ip == "all" {
		constraint = globalConstraint
	} else {
		constraint, _ = findConstraint(upgradeContent.Ip, false)
	}
	constraint.Constraint = upgradeContent.Name

	request_id := requestId(r)
	log.Printf("count#upgraded name=%s ip=%s request_id=%s", upgradeContent.Name, upgradeContent.Ip, request_id)
}
