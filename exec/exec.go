package exec

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"../errors"
	"../filesystem"
	"../objects/strings"
	"../runtime"
)

func Run(program string, args []string) {
	argsStr := strings.Join(args, " ")
	sourceFileExists(program)
	c := exec.Command(program, args...)
	err := c.Run()
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to execute command: "+program+" "+argsStr, false, false, true)
}

func WithLiveLogs(program string, args []string, printLogs bool) error {
	argsStr := strings.Join(args, " ")
	sourceFileExists(program)
	cmd := exec.Command(program, args...)
	var err error
	var errMsg string
	var out io.Reader
    {
		stdout, err := cmd.StdoutPipe()
		errMsg = "Unable to get StdoutPipe pipe: "+program+" "+argsStr
		if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true) {
			err = errors.New(errMsg)
		}
		stderr, err := cmd.StderrPipe()
		errMsg = "Unable to get StderrPipe pipe: "+program+" "+argsStr
        if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true) {
			err = errors.New(errMsg)
		}
        out = io.MultiReader(stdout, stderr)
    }
	err = cmd.Start()
	errMsg = "Unable to execute command: "+program+" "+argsStr
	if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true) {
		err = errors.New(errMsg)
	}
	scannerMulti := bufio.NewScanner(out)
	if printLogs {
		for scannerMulti.Scan() {
			fmt.Println(scannerMulti.Text())
		}
	}
	err = cmd.Wait()
	errMsg = "The command has quit with errors: "+program+" "+argsStr
	if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true) {
		err = errors.New(errMsg)
	}
	return err
}

func sourceFileExists(program string) {
	if strings.Contains(program, ".") && !filesystem.Exists(program) {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to execute command that is based on a program that doesn't exist on this host: "+program, true, true, true)
	}
}