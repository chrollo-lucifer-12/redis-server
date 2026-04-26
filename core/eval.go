package core

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/redis-server/resp"
)

func evalPING(args []string, client io.ReadWriter) error {
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

func evalSET(args []string, client io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(error) ERR wrong number of arguments for 'set' command")
	}

	var key, value string
	var exDurationsMs int64 = -1

	key, value = args[0], args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			i++
			if i == len(args) {
				return errors.New("(error) ERR syntax error")
			}

			exDurationSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			exDurationsMs = exDurationSec * 1000

		default:
			return errors.New("(error) ERR syntax error")
		}
	}
	Put(key, NewObj(value, exDurationsMs))
	client.Write([]byte("+OK\r\n"))
	return nil
}

func evalGET(args []string, client io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'get' command")
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		client.Write([]byte("$-1\r\n"))
		return nil
	}

	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		client.Write([]byte("$-1\r\n"))
		return nil
	}

	client.Write(resp.Encode(obj.Value, false))
	return nil
}

func evalTTL(args []string, client io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'get' command")
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		client.Write([]byte(":-2\r\n"))
		return nil
	}

	if obj.ExpiresAt == -1 {
		client.Write([]byte(":-1\r\n"))
		return nil
	}

	durationMs := obj.ExpiresAt - time.Now().UnixMilli()

	if durationMs < 0 {
		client.Write([]byte(":-2\r\n"))
		return nil
	}

	client.Write(resp.Encode(durationMs/1000, false))
	return nil
}

func EvalAndInput(cmd *RedisCmd, client io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, client)
	case "SET":
		return evalSET(cmd.Args, client)
	case "GET":
		return evalGET(cmd.Args, client)
	case "TTL":
		return evalTTL(cmd.Args, client)
	default:
		return evalPING(cmd.Args, client)
	}

}
