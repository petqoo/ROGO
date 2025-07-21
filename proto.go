package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tidwall/resp"
)

type Command interface{}

type Setcommand struct {
	Key   string
	Value string
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v,_, err := rd.ReadValue( )
		if err ==io.EOF{
			break
		}
		fmt.Println("9rina ", v, err)
		if v.Type() == resp.TypeArray {
			fmt.Printf("Received array command: %s\n", v)
	}
}