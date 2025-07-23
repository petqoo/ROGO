package main

import (
	"bytes"
	"fmt"
	"io"
	"github.com/tidwall/resp"
)

type Command interface{}
const (
	Setcommand = "SET"
	GetCommand = "GET"
)
type SetCommand struct {
	Key   string
	Value string
	Peer *Peer
}
type Getcommand struct {
	Key string
	Peer *Peer
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))

	v, _, err := rd.ReadValue()
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF 
		}
		return nil, err
	}

	if v.Type() != resp.Array {
		return nil, fmt.Errorf("invalid command: expected RESP array")
	}

	arr := v.Array()
	if len(arr) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	commandName := arr[0].String()

	switch commandName {
	case Setcommand:
		if len(arr) != 3 {
			return nil, fmt.Errorf("invalid SET command: expected 3 arguments, got %d", len(arr))
		}
	
		cmd := SetCommand{
			Key:   arr[1].String(),
			Value: arr[2].String(),
		}
		return &cmd, nil
    case GetCommand:
		if len(arr) != 2 {
			return nil, fmt.Errorf("invalid GET command: expected 2 arguments, got %d", len(arr))
		}
		cmd := Getcommand{
			Key: arr[1].String(),
		}
		return &cmd, nil

	//  GET command
	default:
		return nil, fmt.Errorf("unknown command: %s", commandName)
	}
}