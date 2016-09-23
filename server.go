package pizzapi

import (
	"log"
	"net/http"

	"github.com/braintree/manners"
	"github.com/gorilla/mux"
)

func app() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/pizzas", pizzasList).Methods("GET")
	r.HandleFunc("/orders", findOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", findOrder).Methods("GET")

	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/upgrade", upgrade).Methods("POST")

	return applyConstraints(logRequest(r.ServeHTTP))
}

func StartServer(port string, shutdown <-chan struct{}) {
	log.Printf("http.start.port=%s\n", port)
	go listenForShutdown(shutdown)

	if err := manners.ListenAndServe(":"+port, app()); err != nil {
		log.Fatalf("server.server error=%v", err)
	}
}

func listenForShutdown(shutdown <-chan struct{}) {
	log.Println("http.graceful.await")
	<-shutdown
	log.Println("http.graceful.shutdown")
	manners.Close()
}
