package id

import (
	"github.com/neurafuse/tools-go/filesystem"
)

var active string = filesystem.GetWorkingDir(true)

func GetActive() string {
	return active
}

func SetActive(id string) {
	active = id
}
