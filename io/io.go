package io

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"../errors"
	"../filesystem"
	"../logging"
	"../runtime"
	"../timing"
	"../vars"
)

type F struct{}

func (f F) Reachable(addrs string) bool {
	var reachable bool = false
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
	}
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	logging.Log([]string{"", vars.EmojiDir, vars.EmojiSuccess}, "Download finished.", 0)
	return nil
}
