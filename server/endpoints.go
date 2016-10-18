package server

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
	http.Redirect(w, r, "/doc", http.StatusFound)
}

func pizzasList(w http.ResponseWriter, r *http.Request) {
	list, _ := PizzaList()

	if err := serveResponse(w, 200, list); err != nil {
		panic(err)
	}
}

func findOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	order, err := FindOrder(id)
	if err != nil {
		if err := serveResponse(w, 404, &errorResponse{"not_found", err.Error()}); err != nil {
			panic(err)
		}
		return
	}

	if err := serveResponse(w, 200, order); err != nil {
		panic(err)
	}
}

func findOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := AllOrders()
	if err != nil {
		if err := serveResponse(w, 404, &errorResponse{"error", err.Error()}); err != nil {
			panic(err)
		}
		return
	}

	if err := serveResponse(w, 200, orders); err != nil {
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
		if err := serveResponse(w, 422, &errorResponse{"invalid", fmt.Sprintf("Invalid data. Expected '{\"id\":1}'. Got %s", body)}); err != nil {
			panic(err)
		}

		return
	}

	pizza, err := FindPizza(orderContent.Pizza)
	if err != nil {
		if err := serveResponse(w, 404, &errorResponse{"not_found", err.Error()}); err != nil {
			panic(err)
		}
		return
	}

	order := CreateOrder(pizza)
	if err := serveResponse(w, http.StatusCreated, order); err != nil {
		panic(err)
	}
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	var upgradeContent struct {
		Name  string `json:"name"`
		Token string `json:"token"`
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &upgradeContent); err != nil {
		if err := serveResponse(w, 422, err); err != nil {
			panic(err)
		}
		return
	}

	token := r.Header.Get("Authorization")
	if token != os.Getenv("UPGRADE_KEY") {
		if err := serveResponse(w, 401, &errorResponse{"unauthorized", "Not authorized"}); err != nil {
			panic(err)
		}
		return
	}

	buildConstraint(upgradeContent.Token, upgradeContent.Name)

	request_id := requestId(r)
	log.Printf("count#upgraded name=%s token=%s request_id=%s", upgradeContent.Name, upgradeContent.Token, request_id)
	if err := serveResponse(w, 200, upgradeContent); err != nil {
		panic(err)
	}
}
