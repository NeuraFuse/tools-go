package runtime

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/neurafuse/tools-go/errors"
)

type F struct{}

func (f F) GetCallerInfo(packageOnly bool) string {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if !ok || details == nil {
		errors.Check(nil, "", "", true, false, true)
	}
	info := strings.SplitAfter(details.Name(), "/")
	infoCaller := info[len(info)-1] + "()"
	if packageOnly {
		infoCaller = strings.Split(infoCaller, ".")[0]
	}
	return infoCaller
}

func (f F) GetExecParams(serviceID, action string) (string, []string) {
	var execProgram string
	var execProgramArgs []string
	if f.OSActive("linux") {
		execProgram = "service"
		execProgramArgs = []string{action, serviceID}
	} else if f.OSActive("darwin") {
		var actionT string
		switch action {
		case "start":
			actionT = "open"
		case "stop":
			actionT = "close"
		}
		execProgram = actionT
		execProgramArgs = []string{"-a", strings.Title(serviceID)}
	}
	return execProgram, execProgramArgs
}

func (f F) GetOS() string {
	return runtime.GOOS
}

func (f F) OSActive(os string) bool {
	var active bool
	if f.GetOS() == os {
		active = true
	}
	return active
}

func (f F) GetOSArchitecture() string {
	return runtime.GOARCH
}

func (f F) GetOSArchitecturePairs(os string) [][]string {
	var archiPairs [][]string
	switch os {
	case "linux":
		archiPairs = [][]string{{"linux", "386"}, {"linux", "amd64"}, {"linux", "arm"}, {"linux", "arm64"},
			{"linux", "mips"}, {"linux", "mips64"}, {"linux", "mips64le"}, {"linux", "mipsle"}, {"linux", "ppc64"}, {"linux", "ppc64le"},
			{"linux", "riscv64"}, {"linux", "s390x"}}
	case "macos":
		archiPairs = [][]string{{"darwin", "amd64"}}
	}
	return archiPairs
}

func (f F) GetOSUsername() string {
	user, err := user.Current()
	errors.Check(err, f.GetCallerInfo(false), "Unable to get current user!", false, false, true)
	return user.Username
}

func (f F) GetSudoUser() string {
	var username string = os.Getenv("SUDO_USER")
	return username
}

func (f F) GetOSUserGid() string {
	var user *user.User
	return user.Gid
}

func (f F) GetOSInstallDir() string {
	return "/usr/local/"
}

func (f F) GetRunningExecutable() string {
	ex, _ := os.Executable()
	return filepath.Base(ex)
}
