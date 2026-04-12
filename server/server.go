package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/redis-server/config"
)

func readCommand(client net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)

	n, err := client.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func respond(cmd string, client net.Conn) error {
	if _, err := client.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

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

		for {
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				con_clients -= 1
				log.Println("client disconnected with address: ", c.RemoteAddr(), "concurrent clients: ", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}

			log.Println("command", cmd)

			if err = respond(cmd, c); err != nil {
				log.Println("error in write: ", err)
			}
		}
	}
}
