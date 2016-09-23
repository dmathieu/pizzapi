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
	r.PathPrefix("/doc").Handler(http.StripPrefix("/doc", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/upgrade", upgrade).Methods("POST")

	r.HandleFunc("/pizzas", applyConstraints(pizzasList)).Methods("GET")
	r.HandleFunc("/orders", applyConstraints(findOrders)).Methods("GET")
	r.HandleFunc("/orders/{id}", applyConstraints(findOrder)).Methods("GET")
	r.HandleFunc("/orders", applyConstraints(createOrder)).Methods("POST")

	return logRequest(r.ServeHTTP)
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
