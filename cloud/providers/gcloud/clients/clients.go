package clients

import (
	"context"

	container "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/storage"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/timing"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

type F struct{}

func (f F) GetToken() *oauth2.Token {
	ctx := context.Background()
	creds, err := transport.Creds(ctx, f.GetClientOptions())
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, true, true)
	token, err := creds.TokenSource.Token()
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, true, true)
	return token
}

func (f F) GetClientOptions() option.ClientOption {
	return option.WithCredentialsFile(infraConfig.F.GetInfraGcloudAuthPath(infraConfig.F{}))
}

func (f F) GetServiceAccount() string {
	return filesystem.FileToString(infraConfig.F.GetInfraGcloudAuthPath(infraConfig.F{}))
}

func (f F) GetContainer() (context.Context, *container.ClusterManagerClient) {
	ctx, _ := context.WithTimeout(context.Background(), timing.GetTimeDuration(30, "m"))
	client, err := container.NewClusterManagerClient(ctx, f.GetClientOptions())
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create new ClusterManagerClient!", false, false, true)
	return ctx, client
}

func (f F) GetStorage() (context.Context, *storage.Client) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, f.GetClientOptions())
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, true, true)
	return ctx, client
}
