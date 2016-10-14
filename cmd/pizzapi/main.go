package main

import _ "github.com/joho/godotenv/autoload"
import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/dmathieu/pizzapi/server"
)

func main() {
	fmt.Println("Starting API...")
	flag.Parse()

	_, err := strconv.Atoi(*server.HttpPort)
	if err != nil {
		log.Printf("%s: $PORT must be an integer value. - %s\n", *server.HttpPort, err)
		os.Exit(1)
	}

	quit := server.UpdateStatus(1 * time.Minute)
	server.StartServer(*server.HttpPort, server.AwaitSignals(syscall.SIGURG))
	quit <- true
}
