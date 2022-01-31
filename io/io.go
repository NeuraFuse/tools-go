package io

import (
	"io"
	"net"
	"net/http"
	"os"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Reachable(addrs string) bool {
	var reachable bool
	_, err := net.DialTimeout("tcp", addrs+":443", timing.GetTimeDuration(1, "s"))
	if err == nil {
		reachable = true
	}
	return reachable
}

func (f F) DownloadFile(filePath string, url string) error {
	logging.Log([]string{"", vars.EmojiDir, vars.EmojiProcess}, "Downloading from url: "+url+" --> "+filePath, 0)
	if filesystem.Exists(filePath) {
		filesystem.Delete(filePath, false)
	} else {
		filesystem.CreateEmptyDir(filePath)
	}
	out, err := os.Create(filePath)
	if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create file: "+filePath, false, false, true) {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get http url: "+url, false, false, true)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("HTTP response bad status: " + resp.Status)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	logging.Log([]string{"", vars.EmojiDir, vars.EmojiSuccess}, "Download finished.", 0)
	return nil
}
