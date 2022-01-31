package app

import (
	acc "github.com/neurafuse/tools-go/ci/accelerator"
	"github.com/neurafuse/tools-go/ci/base"
	"github.com/neurafuse/tools-go/container"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/runtime"
)

type F struct{}

func (f F) Prepare() string {
	return acc.F.Prepare(acc.F{}, f.GetContext(), base.F.GetResType(base.F{}, f.GetContext()))
}

func (f F) Create() string {
	return acc.F.Create(acc.F{}, f.GetContext(), base.F.GetNamespace(base.F{}), base.F.GetResType(base.F{}, f.GetContext()), container.F.GetImgAddrs(container.F{}, f.GetContext(), false, false), base.F.GetResources(base.F{}, f.GetContext()), base.F.GetVolumes(base.F{}, f.GetContext()))
}

func (f F) update() {
	acc.F.Update(acc.F{}, f.GetContext(), base.F.GetNamespace(base.F{}), base.F.GetResType(base.F{}, f.GetContext()), container.F.GetImgAddrs(container.F{}, f.GetContext(), false, false), base.F.GetResources(base.F{}, f.GetContext()), base.F.GetVolumes(base.F{}, f.GetContext()))
}

func (f F) Delete() string {
	volumes := [][]string{} // Don't delete volumes (recycle for inference)
	return acc.F.Delete(acc.F{}, f.GetContext(), base.F.GetNamespace(base.F{}), base.F.GetResType(base.F{}, f.GetContext()), volumes)
}

func (f F) GetContext() string {
	return env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)
}
