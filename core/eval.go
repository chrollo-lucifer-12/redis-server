package core

import (
	"errors"
	"net"

	"github.com/redis-server/resp"
)

func evalPING(args []string, client net.Conn) error {
	var b []byte

	if len(args) >= 2 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 0 {
		b = resp.Encode("PONG", true)
	} else {
		b = resp.Encode(args[0], false)
	}
	_, err := client.Write(b)
	return err
}

func EvalAndInput(cmd *RedisCmd, client net.Conn) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, client)
	default:
		return evalPING(cmd.Args, client)
	}
}
