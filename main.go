package main

import (
	"flag"
	"log"

	"github.com/redis-server/config"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host")
	flag.IntVar(&config.Port, "port", 6379, "port")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("starting the server...")
}
