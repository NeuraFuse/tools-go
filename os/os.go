package os

import (
	"os/user"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/runtime"
)

func GetHostUID(username string) string {
	user, err := user.Lookup(username)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get UID for username: "+username+" !", false, false, true)
	return user.Uid
}

func GetHostGID(username string) string {
	user, err := user.Lookup(username)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get GID for username: "+username+" !", false, false, true)
	return user.Gid
}
