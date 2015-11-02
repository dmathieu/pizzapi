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

	"github.com/braintree/manners"
	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

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
		response := &ErrorResponse{Id: "not_found", Message: "Order not found"}
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
		response := &ErrorResponse{Id: "invalid", Message: fmt.Sprintf("Invalid data. Expected '{\"id\":1}'. Got %s", body)}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	pizza, err := FindPizza(orderContent.Pizza)
	if err != nil {
		w.WriteHeader(422)
		response := &ErrorResponse{Id: "not_found", Message: "Pizza not found"}
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

type upgradeBody struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Key  string `json:"key"`
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	var upgradeContent upgradeBody
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
		response := &ErrorResponse{Id: "invalid", Message: "Invalid data"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	if upgradeContent.Key != os.Getenv("UPGRADE_KEY") {
		response := &ErrorResponse{Id: "unauthorized", Message: "Not authorized"}
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
	log.Printf("count#upgraded name=%s ip=%s request_id=", upgradeContent.Name, upgradeContent.Ip, request_id)
}

func app() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/pizzas", pizzasList).Methods("GET")
	r.HandleFunc("/orders/{id}", findOrder).Methods("GET")

	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/upgrade", upgrade).Methods("POST")

	return applyConstraints(logRequest(r.ServeHTTP))
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
