package server

import (
	"log"
	"net"
	"strconv"

	"github.com/redis-server/config"
)

func RunServer() {
	log.Println("starting the server on ", config.Host, config.Port)

	var con_clients int = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))

	if err != nil {
		panic(err)
	}

	for {
		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}
		con_clients += 1
		log.Println("client connected with address: ", c.RemoteAddr(), "concurrent clients: ", con_clients)
	}
}
