package podman

import (
	"context"
	"fmt"
	"os"

	"../../errors"
	"../../logging"
	"../../objects/strings"
	"../../runtime"
	"../../vars"
	"github.com/containers/libpod/v2/pkg/bindings"
	"github.com/containers/libpod/v2/pkg/bindings/images"
	"github.com/containers/libpod/v2/pkg/domain/entities"
)

type F struct{}

func (f F) getContext() string {
	// Get Podman socket location
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	socket := "unix:" + sock_dir + "/podman/podman.sock"
	// Connect to Podman socket
	connText, err := bindings.NewConnection(context.Background(), socket)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get new connection!", false, true, true)
	return connText
}

func (f F) imagesPull(addrs string) {
	_, err := images.Pull(f.getContext(), addrs, entities.ImagePullOptions{})
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to pull image: "+addrs, false, true, true)
}

func (f F) imagesList() {
	imageSummary, err := images.List(f.getContext(), nil, nil)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to list images!", false, true, true)
	var names []string
	for _, i := range imageSummary {
		names = append(names, i.RepoTags...)
	}
	fmt.Println()
	fmt.Println(names)
	logging.Log([]string{"", vars.EmojiContainer, vars.EmojiInfo}, "Listing images:", 0)
	logging.Log([]string{"", vars.EmojiInfo, vars.EmojiInspect}, strings.Join(names, "\n"), 0)
}