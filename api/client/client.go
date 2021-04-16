package client

import (
	"bytes"
	"io"
	"strings"

	"../../../neurakube/infrastructure/ci/api"
	"../../config"
	devConfig "../../config/dev"
	infraConfig "../../config/infrastructure"
	projectConfig "../../config/project"
	userConfig "../../config/user"
	"../../errors"
	"../../filesystem"
	"../../logging"
	"../../projects"
	"../../readers/yaml"
	"../../runtime"
	"../../terminal"
	"../../users"
	"../../vars"
	"./auth"
	"./http"
)

type F struct{}

var APIconnected bool = false
var loggedConnect bool = false
var userTypeAdmin string = "user_admin"
var userTypeStandard string = "user_standard"
var userTypeContainer string = "container"

func (f F) Router(context, method, route, cookies, contentType, queryParams string, body io.Reader) string {
	api.F.EvalAction(api.F{}, "connect")
	if !APIconnected {
		APIconnected = true
		f.Connect()
	}
	return http.F.RequestHandler(http.F{}, context, method, route, f.getCookieProjectID()+cookies, contentType, queryParams, body, auth.F.GetJwtAuthStatus(auth.F{}))
}

func (f F) Connect() {
	if !loggedConnect {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiLink}, "Creating secure connection to "+vars.NeuraKubeName+"..", 0)
		loggedConnect = true
	}
	connectStatus := f.ConnectStatusCache("get", vars.NeuraKubeNameRepo, "")
	if connectStatus == "true" {
		auth.F.Check(auth.F{}, true)
	} else {
		f.checkUserAccount()
	}
}

func (f F) ConnectStatusCache(method, module, value string) string {
	return config.Setting(method, "infrastructure", "Spec."+strings.Title(module)+".Cache.ConnectStatus", value)
}

func (f F) getCookieProjectID() string {
	return "projectID=" + projects.IDActive + ";"
}

func (f F) Inspect() {
	auth.F.Check(auth.F{}, false)
	answer := strings.Split(f.Router(vars.NeuraKubeNameRepo, "GET", "inspect", "", "", "", nil), ",")
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiInspect}, vars.NeuraKubeName+" inspection\n", 0)
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiInfo}, "Version: "+answer[0], 0)
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiInfo}, "Health check: "+answer[1], 0)
	emojiInfrastructure := ""
	if answer[2] == "initialized" {
		emojiInfrastructure = vars.EmojiSuccess
	} else {
		emojiInfrastructure = vars.EmojiWarning
	}
	logging.Log([]string{"", vars.EmojiAPI, emojiInfrastructure}, "Infrastructure: "+answer[2]+"\n", 0)
}

func (f F) checkUserAccount() {
	createUser, userType := f.checkAccountStatus(vars.NeuraKubeNameRepo)
	if createUser {
		f.createAccount(userType)
		f.ConnectStatusCache("set", vars.NeuraKubeNameRepo, "true")
	}
}

func (f F) checkAccountStatus(context string) (bool, string) {
	createUser := false
	auth.F.Check(auth.F{}, false)
	init := http.F.RequestHandler(http.F{}, context, "GET", "inspect/infrastructure/init", "", "", "", nil, false)
	var userType string
	if init == "initialized" {
		if strings.Contains(users.GetIDActive(), "container") {
			userType = userTypeContainer
		} else {
			userType = userTypeStandard
		}
		auth.F.Check(auth.F{}, true)
		userList := f.Router(context, "GET", "users", "", "", "", nil)
		if !strings.Contains(userList, users.GetIDActive()) {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiUser}, "User "+users.GetIDActive()+" is not registered yet.", 0)
			createUser = true
		}
	} else {
		userType = userTypeAdmin
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWarning}, vars.NeuraKubeName+" is not initialized yet.\n", 0)
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiProcess}, "Starting initialization..", 0)
		createUser = true
	}
	return createUser, userType
}

func (f F) createAccount(userType string) {
	f.createUser(userType)
	auth.F.Check(auth.F{}, true)
	if strings.Contains(userType, "user") {
		f.createdevConfig()
	}
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Account "+users.GetIDActive()+" is now synced.", 0)
	if userType == "admin" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, vars.NeuraKubeName+" infrastructure is now initialized.\n", 0)
	}
}

func (f F) Infrastructure(action string) {
	if action == "delete" || action == "del" {
		f.Router(vars.NeuraKubeNameRepo, "GET", "infrastructure", "", "", "", nil)
		f.ResetCaches()
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Cluster deleted.", 0)
	}
}

func (f F) ResetCaches() {
	http.F.EndpointCache(http.F{}, "reset", vars.NeuraKubeNameRepo, "")
	f.ConnectStatusCache("reset", vars.NeuraKubeNameRepo, "")
}

func (f F) createUser(userType string) {
	logging.Log([]string{"", vars.EmojiUser, vars.EmojiProcess}, "Registering new "+userType+" user: "+users.GetIDActive(), 0)
	status := f.Router(vars.NeuraKubeNameRepo, "POST", "user/create", "", "application/yaml", "", bytes.NewReader(yaml.ToJSON(userConfig.F.GetConfig(userConfig.F{}))))
	if status == "success" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "User "+users.GetIDActive()+" synced.", 0)
	} else {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWarning}, status, 0)
		terminal.Exit(1, "Failed to create user!")
	}
}

func (f F) Sync() {
	f.syncInfra()
	f.syncProject()
}

func (f F) syncInfra() {
	f.projectAuth()
	status := f.Router(vars.NeuraKubeNameRepo, "POST", "user/infrastructure", "", "application/yaml", "", bytes.NewReader(yaml.ToJSON(infraConfig.F.GetConfig(infraConfig.F{}))))
	if status == "success" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Infrastructure synced.", 0)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Failed to sync infrastructure config: "+status, true, true, true)
	}
}

func (f F) syncProject() {
	var projectID string = projects.IDActive
	f.projectAuth()
	status := f.Router(vars.NeuraKubeNameRepo, "POST", "user/project", "", "application/yaml", "", bytes.NewReader(yaml.ToJSON(projectConfig.F.GetConfig(projectConfig.F{}))))
	if status == "success" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Project "+projectID+" synced.", 0)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Failed to sync project "+projectID+": "+status, true, true, true)
	}
}

func (f F) projectAuth() {
	var status string
	if config.ValidSettings("infrastructure", "kube", false) {
		status = f.Router(vars.NeuraKubeNameRepo, "POST", "user/infrastructure/auth/kubeconfig/create", "", "application/json", "", bytes.NewReader(filesystem.FileToBytes(config.Setting("get", "infrastructure", "Spec.Cluster.Auth.KubeConfigPath", ""))))
	} else if config.ValidSettings("infrastructure", vars.InfraProviderGcloud, false) {
		status = f.Router(vars.NeuraKubeNameRepo, "POST", "user/infrastructure/auth/gcloud/create", "", "application/json", "", bytes.NewReader(filesystem.FileToBytes(config.Setting("get", "infrastructure", "Spec.Gcloud.Auth.ServiceAccountJSONPath", ""))))
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Invalid Project auth!", true, true, true)
	}
	if status == "success" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Infrastructure auth synced.", 0)
	} else {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWarning}, status, 0)
		terminal.Exit(1, "")
	}
}

func (f F) createdevConfig() {
	if devConfig.F.Exists(devConfig.F{}) {
		status := f.Router(vars.NeuraKubeNameRepo, "POST", "user/devConfig/create", "", "application/yaml", "", bytes.NewReader(filesystem.FileToBytes(devConfig.F.GetFilePath(devConfig.F{}))))
		if status == "success" {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "devConfig synced.", 0)
		} else {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWarning}, status, 0)
			terminal.Exit(1, "")
		}
	}
}
