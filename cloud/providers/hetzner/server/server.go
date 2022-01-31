package server

import (
	"context"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"../client"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/runtime"
)

type F struct{}

func (f F) Create(id string) {
	client := client.F.Get(client.F{})
	hcloud.ServerClient.Create(context.Background(), f.getCreateOpts(id))
}

func (f F) Delete() {

}

func (f F) GetByID(id int) *hcloud.Server {
	client := client.F.Get(client.F{})
	server, _, err := client.Server.GetByID(context.Background(), id)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to retrieve server (ID: "+id+")!", false, false, true)
	return server
}

func (f F) getCreateOpts(id string) *hcloud.ServerCreateOpts {
	createOpts := &hcloud.ServerCreateOpts{}
	createOpts.Name = id
	createOpts.ServerType = f.getType()
	createOpts.Image = f.getImage()
	createOpts.SSHKeys = f.getSSHKeys()
	createOpts.Datacenter = f.getDatacenter()
	createOpts.StartAfterCreate = &true
	/*type ServerCreateOpts struct {
		Name             id
		ServerType       getServerType()
		Image            *Image
		SSHKeys          []*SSHKey
		Location         *Location
		Datacenter       *Datacenter
		UserData         string
		StartAfterCreate *bool
		Labels           map[string]string
		Automount        *bool
		Volumes          []*Volume
		Networks         []*Network
	}*/
	return createOpts
}

func (f F) getType() *hcloud.ServerType {
	typ := &hcloud.ServerType{}
	typ.Name = "cx11"
	/*serverType := {
		ID          int
		Name        string
		Description string
		Cores       int
		Memory      float32
		Disk        int
		StorageType StorageType
		CPUType     CPUType
		Pricings    []ServerTypeLocationPricing
	}*/
	return typ
}

func (f F) getImage() *hcloud.Image {
	img := &hcloud.Datacenter{}
	img.Name = "ubuntu-20.04"
	return img
}

func (f F) getSSHKeys() []*hcloud.SSHKey {
	key := []*&hcloud.SSHKey{}
	key[0].Name = "ssh-1"
	publicKeyFilePath, privateKeyFilePath := rsa.GenerateKeys(vars.NeuraKubeName, "ssh/keys/", false)
	key[0].PublicKey = filesystem.FileToString(publicKeyFilePath)
	return key
}

func (f F) getDatacenter() *hcloud.Datacenter {
	dc := &hcloud.Datacenter{}
	dc.ID = "nbg1-dc3"
	return dc
}