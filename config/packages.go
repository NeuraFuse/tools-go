package config

import (
	"./cli"
	"./dev"
	"./server"
	"./infrastructure"
	"./user"
)

type Packages struct {
	Cli            cli.F
	Dev            dev.F
	Infrastructure infrastructure.F
	Server         server.F
	User           user.F
}
