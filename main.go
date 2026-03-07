package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not parse .env file %v", err)
		panic(err)
	}

	var port int
	flag.IntVar(&port, "port", 8800, "Starting port of the app")
	flag.Parse()
}
