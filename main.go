package main

import (
	"flag"
	"log"

	"github.com/redis-server/config"
	"github.com/redis-server/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "localhost", "host")
	flag.IntVar(&config.Port, "port", 6379, "port")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("starting the server...")
	server.RunAsyncServer()
}
