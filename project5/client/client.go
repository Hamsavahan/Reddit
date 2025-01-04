package client

import (
	"fmt"
	"protoactor-simulation/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type Client struct {
	username string
	engine   *actor.PID
}

func NewClient(username string, engine *actor.PID) *Client {
	return &Client{
		username: username,
		engine:   engine,
	}
}

func (c *Client) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Printf("Client %s registered, Welcome to reddit %s\n ", c.username, c.username)
		context.Send(c.engine, &messages.Register{Username: c.username})
	case *messages.Response:
		fmt.Printf("Client %s received response: %s\n", c.username, msg.Message)

	}
}
