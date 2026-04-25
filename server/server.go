package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/redis-server/config"
	"github.com/redis-server/core"
	"github.com/redis-server/resp"
)

func readCommand(client net.Conn) (*core.RedisCmd, error) {
	var buf []byte = make([]byte, 512)

	n, err := client.Read(buf)
	if err != nil {
		return nil, err
	}

	tokens, err := resp.DecodeArrayString(buf[:n])

	if err != nil {
		return nil, err
	}

	fmt.Println(tokens)

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondWithError(client net.Conn, err error) {
	client.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCmd, client net.Conn) {
	err := core.EvalAndInput(cmd, client)
	if err != nil {
		respondWithError(client, err)
	}
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

			respond(cmd, c)
		}
	}
}
