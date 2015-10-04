package main

import _ "github.com/joho/godotenv/autoload"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/dmathieu/pizzapi/app"
)

func main() {
	fmt.Println("Starting API...")
	flag.Parse()

	_, err := strconv.Atoi(*app.HttpPort)
	if err != nil {
		log.Printf("%s: $PORT must be an integer value. - %s\n", app.HttpPort, err)
		os.Exit(1)
	}

	app.StartServer(*app.HttpPort, app.AwaitSignals(syscall.SIGURG))
}
