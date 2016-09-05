package main

import _ "github.com/joho/godotenv/autoload"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/dmathieu/pizzapi"
)

func main() {
	fmt.Println("Starting API...")
	flag.Parse()

	_, err := strconv.Atoi(*pizzapi.HttpPort)
	if err != nil {
		log.Printf("%s: $PORT must be an integer value. - %s\n", *pizzapi.HttpPort, err)
		os.Exit(1)
	}

	pizzapi.StartServer(*pizzapi.HttpPort, pizzapi.AwaitSignals(syscall.SIGURG))
}
