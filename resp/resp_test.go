package resp

import (
	"reflect"
	"testing"
)

func TestReadSimpleString(t *testing.T) {
	input := []byte("+OK\r\n")

	val, n, err := readSimpleString(input)
	if err != nil {
		t.Fatal(err)
	}

	if val != "OK" {
		t.Errorf("expected OK, got %s", val)
	}
	if n != len(input) {
		t.Errorf("expected %d, got %d", len(input), n)
	}
}

func TestReadError(t *testing.T) {
	input := []byte("-ERR something went wrong\r\n")

	val, _, _ := readError(input)
	if val != "ERR something went wrong" {
		t.Errorf("unexpected error string: %s", val)
	}
}

func TestReadInt(t *testing.T) {
	input := []byte(":12345\r\n")

	val, _, _ := readInt(input)
	if val != 12345 {
		t.Errorf("expected 12345, got %d", val)
	}
}

func TestReadBulkString(t *testing.T) {
	input := []byte("$5\r\nhello\r\n")

	val, _, _ := readBulkString(input)
	if val != "hello" {
		t.Errorf("expected hello, got %s", val)
	}
}

func TestReadArray(t *testing.T) {
	input := []byte("*2\r\n+OK\r\n:42\r\n")

	val, _, err := readArray(input)
	if err != nil {
		t.Fatal(err)
	}

	expected := []interface{}{"OK", int64(42)}

	if !reflect.DeepEqual(val, expected) {
		t.Errorf("expected %v, got %v", expected, val)
	}
}

func TestDecodeOne(t *testing.T) {
	tests := []struct {
		input    []byte
		expected interface{}
	}{
		{[]byte("+OK\r\n"), "OK"},
		{[]byte(":99\r\n"), int64(99)},
		{[]byte("$3\r\nhey\r\n"), "hey"},
		{[]byte("*2\r\n+OK\r\n:1\r\n"), []interface{}{"OK", int64(1)}},
	}

	for _, tt := range tests {
		val, _, err := DecodeOne(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(val, tt.expected) {
			t.Errorf("input %q: expected %v, got %v", tt.input, tt.expected, val)
		}
	}
}

func TestDecode(t *testing.T) {
	input := []byte("+PONG\r\n")

	val, err := Decode(input)
	if err != nil {
		t.Fatal(err)
	}

	if val != "PONG" {
		t.Errorf("expected PONG, got %v", val)
	}
}

func TestDecodeEmpty(t *testing.T) {
	_, err := Decode([]byte{})
	if err == nil {
		t.Error("expected error for empty input")
	}
}
