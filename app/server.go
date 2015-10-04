package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/braintree/manners"
	"github.com/gorilla/mux"
)

type orderBody struct {
	Pizza int `json:"id"`
}

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
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
    fmt.Fprintf(w, "Order not found")
    return
	}

	if err := json.NewEncoder(w).Encode(order); err != nil {
		panic(err)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var orderContent orderBody
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &orderContent); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
    fmt.Fprintf(w, "Invalid data")
		return
	}

	pizza, err := FindPizza(orderContent.Pizza)
	if err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
    fmt.Fprintf(w, "Pizza id not found")
		return
	}

	order := CreateOrder(pizza)
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		panic(err)
	}
}

func app() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/pizzas", pizzasList).Methods("GET")
	r.HandleFunc("/orders/{id}", findOrder).Methods("GET")

	r.HandleFunc("/orders", createOrder).Methods("POST")

	return logRequest(r.ServeHTTP)
}

func StartServer(port string, shutdown <-chan struct{}) {
	log.Printf("http.start.port=%s\n", port)
	http.Handle("/", app())
	go listenForShutdown(shutdown)

	if err := manners.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("server.server error=%v", err)
	}
}

func listenForShutdown(shutdown <-chan struct{}) {
	log.Println("http.graceful.await")
	<-shutdown
	log.Println("http.graceful.shutdown")
	manners.Close()
}
