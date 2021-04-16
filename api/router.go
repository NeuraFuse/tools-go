package router

import (
	"./client"
)

type F struct{}

func (f F) Infrastructure(action string) {
	client.F.Infrastructure(client.F{}, action)
}
