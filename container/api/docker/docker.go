package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	gcloudClient "github.com/neurafuse/tools-go/cloud/providers/gcloud/clients"
	"github.com/neurafuse/tools-go/config"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
)

var contextPack string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

var registryUsername string
var registryPassword string

func BuildImage(filePath string, imgTags []string) {
	logging.Log([]string{"", vars.EmojiContainer, vars.EmojiProcess}, "Building container "+contextPack+" image: "+imgTags[0]+"\n", 0)
	ctx, client := createClient()
	buildResponse, err := client.ImageBuild(
		ctx,
		//dockerFileTarReader(filePath, fileName),
		getContext(filePath),
		types.ImageBuildOptions{
			//Context:    dockerFileTarReader(filePath, fileName),
			Context:     getContext(filePath),
			Dockerfile:  "Dockerfile",
			Tags:        imgTags,
			AuthConfigs: GetAuthConfigMap(),
		})
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to build container "+contextPack+" image: "+imgTags[0], false, false, true)
	responseHandler(buildResponse.Body)
	defer buildResponse.Body.Close()
}

func PushImage(imgAddrs string) {
	logging.Log([]string{"", vars.EmojiContainer, vars.EmojiProcess}, "Pushing container "+contextPack+" image: "+imgAddrs+"\n", 0)
	ctx, client := createClient()
	pushResponse, err := client.ImagePush(ctx,
		imgAddrs,
		types.ImagePushOptions{
			RegistryAuth: GetAuth(),
		})
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to push container "+contextPack+" image: "+imgAddrs, false, false, true)
	responseHandler(pushResponse)
}

func responseHandler(reader io.ReadCloser) error {
	defer reader.Close()
	rd := bufio.NewReader(reader)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		logging.Log([]string{"", vars.EmojiContainer, vars.EmojiInfo}, string(n), 0)
	}
	fmt.Println()
	return nil
}

func RenameImage(imgAddrsSource, imgAddrsTarget string) {
	logging.Log([]string{"", vars.EmojiContainer, vars.EmojiProcess}, "Renaming "+contextPack+" image: "+imgAddrsSource+" --> "+imgAddrsTarget, 0)
	ctx, client := createClient()
	err := client.ImageTag(ctx, imgAddrsSource, imgAddrsTarget)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
}

func Initialize(release bool) {
	daemon("start")
	var username string
	var password string
	var configID string
	var configKeyPrefix string = "Spec.Containers.Registry.Auth."
	if release {
		configID = "dev"
	} else {
		configID = "project"
		if infraConfig.F.ProviderIDIsActive(infraConfig.F{}, "gcloud") {
			if !config.ValidSettings(configID, "containers/registry", false) {
				config.Setting("set", configID, configKeyPrefix+"Username", "_json_key")
				config.Setting("set", configID, configKeyPrefix+"Password", gcloudClient.F.GetServiceAccount(gcloudClient.F{}))
			}
		}
	}
	username = config.Setting("get", configID, configKeyPrefix+"Username", "")
	password = config.Setting("get", configID, configKeyPrefix+"Password", "")
	SetRegistryCred(username, password)
}

func SetRegistryCred(username, password string) {
	registryUsername = username
	registryPassword = password
}

func GetAuth() string {
	encodedJSON, err := json.Marshal(GetAuthConfig())
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get "+contextPack+" registry auth.!", false, false, true)
	return base64.URLEncoding.EncodeToString(encodedJSON)
}

func GetAuthConfig() types.AuthConfig {
	return types.AuthConfig{Username: registryUsername,
		Password: registryPassword}
}

func GetAuthConfigMap() map[string]types.AuthConfig {
	m := make(map[string]types.AuthConfig)
	m["Username"] = types.AuthConfig{Username: registryUsername}
	m["Password"] = types.AuthConfig{Password: registryPassword}
	return m
}

func CreateContainer(imgAddrs string) {
	ctx, client := createClient()
	resp, err := client.ContainerCreate(ctx, &container.Config{
		Image: imgAddrs,
	}, nil, nil, nil, "")
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create container!", false, true, true) {
		StartContainer(resp.ID)
	}
}

func StartContainer(id string) {
	ctx, client := createClient()
	err := client.ContainerStart(ctx, id, types.ContainerStartOptions{})
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to start container!", false, true, true) {
		logging.Log([]string{"", vars.EmojiContainer, vars.EmojiInfo}, "Created container: "+id, 0)
	}
}

func PullImage(imgAddrs string) {
	ctx, client := createClient()
	out, err := client.ImagePull(ctx, imgAddrs, types.ImagePullOptions{})
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to pull image: "+imgAddrs+" !", false, false, true)
	io.Copy(os.Stdout, out)
}

func ListImages() {
	ctx, client := createClient()
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to list images!", false, false, true)
	if len(images) != 0 {
		logging.Log([]string{"", vars.EmojiContainer, vars.EmojiInfo}, "List of local images:", 0)
		for _, image := range images {
			logging.Log([]string{"", vars.EmojiContainer, vars.EmojiInfo}, strings.Join(image.RepoTags, ","), 0)
		}
	} else {
		logging.Log([]string{"", vars.EmojiContainer, vars.EmojiInfo}, "There are no "+contextPack+" images to list.", 0)
	}
}

func createClient() (context.Context, *client.Client) {
	ctx := context.Background()
	// const defaultDockerAPIVersion = "v1.37"
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()) // client.WithVersion(defaultDockerAPIVersion)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create a new client with opts!", false, false, true)
	return ctx, client
}

func daemon(action string) {
	if action == "start" {
		logging.Log([]string{"", vars.EmojiContainer, vars.EmojiWaiting}, "Starting "+contextPack+" daemon..", 0)
		ctx, client := createClient()
		var success bool
		var triggeredStart bool
		for ok := true; ok; ok = !success {
			ping, err := client.Ping(ctx)
			if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, false) {
				if !triggeredStart {
					execProgram, execProgramArgs := runtime.F.GetExecParams(runtime.F{}, contextPack, "start")
					exec.WithLiveLogs(execProgram, execProgramArgs, true)
					triggeredStart = true
				}
				logging.Log([]string{"", vars.EmojiContainer, vars.EmojiWaiting}, "Waitig for "+contextPack+" daemon to be ready..", 0)
			} else {
				success = true
				logging.Log([]string{"", vars.EmojiContainer, vars.EmojiSuccess}, "Docker daemon running (v"+ping.APIVersion+").\n", 0)
			}
			timing.Sleep(200, "ms")
		}
	}
}

func getContext(filePath string) io.Reader {
	// Use homedir.Expand to resolve paths like '~/repos/myrepo'
	//filePath, _ = homedir.Expand(filePath)
	ctx, _ := archive.TarWithOptions(filePath, &archive.TarOptions{})
	return ctx
}

func dockerFileTarReader(dockerfilePath string, dockerfileName string) *bytes.Reader {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	dockerFileReader, err := os.Open(dockerfilePath)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to open dockerfilePath: "+dockerfilePath+" !", false, false, true)
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to read dockerFileReader!", false, false, true)
	tarHeader := &tar.Header{
		Name: dockerfileName,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to write tar header!", false, false, true)
	_, err = tw.Write(readDockerFile)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to write dockerfile!", false, false, true)
	return bytes.NewReader(buf.Bytes())
}
