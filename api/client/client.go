package client

import (
	"bytes"
	"io"
	"strings"

	"github.com/neurafuse/tools-go/api/client/auth"
	"github.com/neurafuse/tools-go/api/client/http"
	"github.com/neurafuse/tools-go/ci/api"
	"github.com/neurafuse/tools-go/config"
	devConfig "github.com/neurafuse/tools-go/config/dev"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	projectConfig "github.com/neurafuse/tools-go/config/project"
	userConfig "github.com/neurafuse/tools-go/config/user"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	infraID "github.com/neurafuse/tools-go/infrastructures/id"
	kubeID "github.com/neurafuse/tools-go/kubernetes/client/id"
	kubeConfig "github.com/neurafuse/tools-go/kubernetes/config"
	"github.com/neurafuse/tools-go/logging"
	projectsID "github.com/neurafuse/tools-go/projects/id"
	"github.com/neurafuse/tools-go/readers/yaml"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var APIconnected bool
var loggedConnect bool
var userTypeAdmin string = "user_admin"
var userTypeStandard string = "user_standard"
var userTypeContainer string = "container"

func (f F) Router(context, method, route, cookies, contentType, queryParams string, body io.Reader) string {
	api.F.EvalAction(api.F{}, "connect")
	if !APIconnected {
		APIconnected = true
		f.Connect(context)
		f.sync()
	}
	return http.F.RequestHandler(http.F{}, context, method, route, f.getBaseCookies()+cookies, contentType, queryParams, body, auth.F.GetJwtAuthStatus(auth.F{}))
}

func (f F) Connect(context string) {
	if !loggedConnect {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiLink}, "Creating secure connection to "+vars.NeuraKubeName+"..", 0)
		loggedConnect = true
	}
	var connectStatus string = f.ConnectStatusCache(context, "get", vars.NeuraKubeNameID, "")
	if connectStatus == "true" {
		auth.F.Check(auth.F{}, true)
	} else {
		f.checkUserAccount(context)
	}
}

func (f F) sync() {
	logging.Log([]string{"", vars.EmojiClient, vars.EmojiAPI}, "Syncing client data..", 0)
	f.syncInfra()
	f.syncProject()
}

func (f F) ConnectStatusCache(context, method, module, value string) string {
	if !f.getAPIInitStatus(context) {
		config.Setting("set", "infrastructure", "Spec."+strings.Title(module)+".Cache.ConnectStatus", "false")
		return "false"
	}
	return config.Setting(method, "infrastructure", "Spec."+strings.Title(module)+".Cache.ConnectStatus", value)
}

func (f F) getBaseCookies() string {
	var cookies string
	cookies = f.getCookieProjectID() + f.getCookieInfraID() + f.getCookieKubeID()
	return cookies
}

func (f F) getCookieInfraID() string {
	var id string = infraID.F.GetActive(infraID.F{})
	if id == "" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get infraID!", true, true, true)
	}
	return "infraID=" + id + ";"
}

func (f F) getCookieKubeID() string {
	var id string = kubeID.F.GetActive(kubeID.F{})
	if id == "" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get clusterID!", true, true, true)
	}
	return "kubeID=" + id + ";"
}

func (f F) getCookieProjectID() string {
	var id string = projectsID.GetActive()
	if id == "" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get projectID!", true, true, true)
	}
	return "projectID=" + id + ";"
}

func (f F) Inspect() {
	auth.F.Check(auth.F{}, false)
	var answer []string = strings.Split(f.Router(vars.NeuraKubeNameID, "GET", "inspect", "", "", "", nil), ",") // TODO: "," potential bug
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiInspect}, vars.NeuraKubeName+" inspection\n", 0)
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Version: "+answer[0], 0)
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Health check: "+answer[1], 0)
	var emojiInfrastructure string
	var checkAccountStatus bool
	if answer[2] == "initialized" {
		emojiInfrastructure = vars.EmojiSuccess
	} else {
		emojiInfrastructure = vars.EmojiWarning
		checkAccountStatus = true
	}
	logging.Log([]string{"", vars.EmojiAPI, emojiInfrastructure}, "Infrastructure: "+answer[2]+"\n", 0)
	if checkAccountStatus {
		f.checkAccountStatus(vars.NeuraKubeNameID)
	}
}

func (f F) checkUserAccount(context string) {
	createUser, userType := f.checkAccountStatus(vars.NeuraKubeNameID)
	if createUser {
		f.createAccount(userType)
		f.ConnectStatusCache(context, "set", vars.NeuraKubeNameID, "true")
	}
}

func (f F) getAPIInitStatus(context string) bool {
	var status bool
	var resp string = http.F.RequestHandler(http.F{}, context, "GET", "inspect/api/init", "", "", "", nil, false)
	if resp == "initialized" {
		status = true
	}
	return status
}

func (f F) checkAccountStatus(context string) (bool, string) {
	var createUser bool
	auth.F.Check(auth.F{}, false)
	var userType string
	if f.getAPIInitStatus(context) {
		if strings.Contains(usersID.F.GetActive(usersID.F{}), "container") {
			userType = userTypeContainer
		} else {
			userType = userTypeStandard
		}
		auth.F.Check(auth.F{}, true)
		var userList string = f.Router(context, "GET", "users", "", "", "", nil)
		if !strings.Contains(userList, usersID.F.GetActive(usersID.F{})) {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiUser}, "User "+usersID.F.GetActive(usersID.F{})+" is not registered yet.", 0)
			createUser = true
		}
	} else {
		userType = userTypeAdmin
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWarning}, vars.NeuraKubeName+" is not initialized yet.\n", 0)
		logging.Log([]string{"", vars.EmojiClient, vars.EmojiAPI}, "Starting initialization..", 0)
		createUser = true
	}
	return createUser, userType
}

func (f F) createAccount(userType string) {
	f.createUser(userType)
	auth.F.Check(auth.F{}, true)
	if strings.Contains(userType, "user") {
		f.createDevConfig()
	}
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Account "+usersID.F.GetActive(usersID.F{})+" is now synced.\n", 0)
	if userType == "admin" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, vars.NeuraKubeName+" infrastructure is now initialized.\n", 0)
	}
}

func (f F) createUser(userType string) {
	logging.Log([]string{"", vars.EmojiUser, vars.EmojiProcess}, "Creating new "+userType+" user: "+usersID.F.GetActive(usersID.F{}), 0)
	var userConfig *userConfig.Default = userConfig.F.GetConfig(userConfig.F{})
	var config []byte = yaml.ToJSON(userConfig)
	var body *bytes.Reader = bytes.NewReader(config)
	var status string = f.Router(vars.NeuraKubeNameID, "POST", "user/create", "", "application/yaml", "", body)
	if status == "success" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "User "+usersID.F.GetActive(usersID.F{})+" synced.", 0)
	} else {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiError}, status, 0)
		terminal.Exit(1, "Failed to create user!")
	}
}

func (f F) Infrastructure(action string) {
	if action == "delete" || action == "del" {
		var status string = f.Router(vars.NeuraKubeNameID, "GET", "user/infrastructure/setup", "", "", "", nil)
		if status == "success" {
			f.ResetCaches(vars.NeuraKubeNameID)
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Cluster deleted.", 0)
		} else {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete the infrastructure: "+status, true, true, true)
		}
	}
}

func (f F) ResetCaches(context string) {
	http.F.EndpointCache(http.F{}, "reset", vars.NeuraKubeNameID, "")
	f.ConnectStatusCache(context, "reset", vars.NeuraKubeNameID, "")
	f.DeleteKubeAuth()
}

func (f F) DeleteKubeAuth() {
	filesystem.Delete(infraConfig.F.GetInfraKubeAuthPath(infraConfig.F{}, true), false)
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Deleted cluster authentication.", 0)
}

func (f F) syncInfra() {
	logging.Log([]string{"", vars.EmojiClient, vars.EmojiKubernetes}, "Syncing infrastructure..", 0)
	var infraConfig []byte = yaml.ToJSON(infraConfig.F.GetConfig(infraConfig.F{}))
	var body *bytes.Reader = bytes.NewReader(infraConfig)
	var status string = f.Router(vars.NeuraKubeNameID, "POST", "user/infrastructure", "", "application/yaml", "", body)
	if status == "success" {
		f.syncGcloudAuth()
		f.syncKubeAuth()
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Infrastructure configuration synced.\n", 0)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to synchronize the infrastructure config: "+status, true, true, true)
	}
}

func (f F) syncProject() {
	var projectID string = projectsID.GetActive()
	var projectConfigFilePath string = projectConfig.F.GetFilePath(projectConfig.F{})
	var projectConfigB []byte = filesystem.FileToBytes(projectConfigFilePath)
	var body *bytes.Reader = bytes.NewReader(projectConfigB)
	logging.Log([]string{"", vars.EmojiClient, vars.EmojiProject}, "Syncing project..", 0)
	var status string = f.Router(vars.NeuraKubeNameID, "POST", "user/project", "", "application/yaml", "", body)
	if status == "success" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Project "+projectID+" synced.", 0)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Failed to sync project "+projectID+": "+status, true, true, true)
	}
}

func (f F) syncKubeAuth() {
	logging.Log([]string{"", vars.EmojiClient, vars.EmojiKubernetes}, "Preparing to synchronize cluster authentication..", 0)
	var status string
	if kubeConfig.F.CheckKubeAuth(kubeConfig.F{}) {
		var kubeConfigFilePath string = infraConfig.F.GetInfraKubeAuthPath(infraConfig.F{}, true)
		var kubeConfig []byte = filesystem.FileToBytes(kubeConfigFilePath)
		var body *bytes.Reader = bytes.NewReader(kubeConfig)
		logging.Log([]string{"", vars.EmojiClient, vars.EmojiKubernetes}, "Syncing cluster authentication..", 0)
		status = f.Router(vars.NeuraKubeNameID, "POST", "user/infrastructure/auth/kubeconfig/create", "", "application/json", "", body)
		if status == "success" {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Synced kubernetes authentication.", 0)
		} else {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to sync kubernetes authentication: "+status, true, true, true)
		}
	} else {
		logging.Log([]string{"", vars.EmojiClient, vars.EmojiKubernetes}, "Skipped synchronization of kubernetes authentication.", 0)
	}
}

func (f F) syncGcloudAuth() {
	logging.Log([]string{"", vars.EmojiClient, vars.EmojiInfra}, "Syncing "+vars.InfraProviderGcloud+" authentication..", 0)
	if config.ValidSettings("infrastructure", vars.InfraProviderGcloud, false) {
		var gCloudSaFilePath string = infraConfig.F.GetInfraGcloudAuthPath(infraConfig.F{})
		var gCloudSa []byte = filesystem.FileToBytes(gCloudSaFilePath)
		var body *bytes.Reader = bytes.NewReader(gCloudSa)
		var status string
		status = f.Router(vars.NeuraKubeNameID, "POST", "user/infrastructure/auth/gcloud/create", "", "application/json", "", body)
		if status == "success" {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Synced "+vars.InfraProviderGcloud+" authentication.\n", 0)
		} else {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to sync "+vars.InfraProviderGcloud+" authentication: "+status, true, true, true)
		}
	}
}

func (f F) createDevConfig() {
	logging.Log([]string{"", vars.EmojiClient, vars.EmojiDev}, "Syncing devConfig..", 0)
	if devConfig.F.Exists(devConfig.F{}) {
		var status string = f.Router(vars.NeuraKubeNameID, "POST", "user/devconfig/create", "", "application/yaml", "", bytes.NewReader(filesystem.FileToBytes(devConfig.F.GetFilePath(devConfig.F{}))))
		if status == "success" {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "devConfig synced.", 0)
		} else {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWarning}, status, 0)
			terminal.Exit(1, "")
		}
	}
}
