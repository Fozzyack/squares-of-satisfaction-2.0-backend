package main

import "flag"



func main() {
	var port int
	flag.IntVar(&port, "port", 8800, "Starting port of the app")
	flag.Parse()
}
