package config

import (
	"github.com/neurafuse/tools-go/config/cli"
	"github.com/neurafuse/tools-go/config/dev"
	"github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/config/server"
	"github.com/neurafuse/tools-go/config/user"
)

type Packages struct {
	Cli            cli.F
	Dev            dev.F
	Infrastructure infrastructure.F
	Server         server.F
	User           user.F
}
