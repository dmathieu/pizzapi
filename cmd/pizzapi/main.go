package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/dmathieu/pizzapi"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	fmt.Println("Starting API...")
	flag.Parse()

	_, err := strconv.Atoi(*pizzapi.HttpPort)
	if err != nil {
		log.Printf("%s: $PORT must be an integer value. - %s\n", *pizzapi.HttpPort, err)
		os.Exit(1)
	}

	quit := pizzapi.UpdateStatus(1 * time.Minute)
	pizzapi.StartServer(*pizzapi.HttpPort, pizzapi.AwaitSignals(syscall.SIGURG))
	quit <- true
}
