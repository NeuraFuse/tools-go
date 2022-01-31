package id

import (
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/runtime"
)

type F struct{}

var idActive string

func (f F) SetActive(id string) {
	idActive = id
}

func (f F) GetActive() string {
	if idActive == "" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "There is no active kubeID set!", true, true, true)
	}
	return idActive
}
