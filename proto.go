package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

type Command interface{}
const (
	SetCommand = "SET"
	GetCommand = "GET"
)
type SetCommand struct {
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
		if err != nil {
			log.Fatal("Error reading value:", err)
		}
		fmt.Println("9rina ", v, err)
		if v.Type() == resp.Array {
			for _ ,value := v.Array() {
				switch value.String() {
					case SetCommand:
						fmt.Println("Set command detected",len(value.Array()))
						if len(v.Array()) != 3 {
							return nil, fmt.Errorf("invalid SET command format")
						}
						cmd:=SetCommand{
							Key:   v.Array()[1].String(),
							Value: v.Array()[2].String(),
						}
						return cmd, nil
	}
}

}
	}
	return  "foo bar", nil

}