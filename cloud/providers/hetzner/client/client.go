package client

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type F struct{}

func (f F) Get() *hcloud.Client {
	return hcloud.NewClient(hcloud.WithToken("token"))
}
