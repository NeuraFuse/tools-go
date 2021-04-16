package http

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"

	"../../../../neurakube/infrastructure/ci"
	"../../../../neurakube/infrastructure/ci/api"
	"../../../../neurakube/infrastructure/ci/base"
	"../../../config"
	"../../../crypto/jwt"
	"../../../env"
	"../../../errors"
	"../../../kubernetes/pods"
	"../../../kubernetes/services"
	"../../../logging"
	"../../../objects/strings"
	"../../../runtime"
	"../../../timing"
	"../../../vars"
	"../../../users"
)

type F struct{}

var extAPI bool
var loggedAPIEndpoint bool

func (f F) getAPIEndpoint(context string) string {
	var apiEndpoint string
	switch context {
	case vars.NeuraKubeNameRepo:
		apiEndpoint = f.endpointneurakube()
	case "inference/gpt":
		apiEndpoint = f.endpointGPT()
	}
	if !loggedAPIEndpoint {
		logging.Log([]string{"", vars.EmojiLink, vars.EmojiInfo}, "Endpoint: "+apiEndpoint, 0)
		loggedAPIEndpoint = true
	}
	return apiEndpoint
}

func (f F) endpointGPT() string {
	var apiIP string
	apiIP = services.F.GetClusterIP(services.F{}, base.F.GetNamespace(base.F{}), ci.F.GetContextID(ci.F{}, env.F.GetContext(env.F{}, "inference", false)))
	return f.generateAPIEndpoint(apiIP, base.F.GetContainerPorts(base.F{}, "inference"))
}

func (f F) endpointneurakube() string {
	var apiIP string
	if f.getAPIExternal() {
		endpointCache := f.EndpointCache("get", vars.NeuraKubeNameRepo, "")
		if endpointCache == "" {
			api.F.EvalAction(api.F{}, "connect")
			apiIP = services.F.GetLoadBalancerIP(services.F{}, api.F.GetNamespace(api.F{}), ci.F.GetContextID(ci.F{}, vars.NeuraKubeNameRepo))
			f.EndpointCache("set", vars.NeuraKubeNameRepo, apiIP)
		} else {
			apiIP = endpointCache
		}
	}
	if apiIP == "" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get "+vars.NeuraKubeName+" endpoint (address is empty)!", true, true, true)
	}
	return f.generateAPIEndpoint(apiIP, api.F.GetContainerPorts(api.F{}))
}

func (f F) getAPIExternal() bool {
	var extAPI bool
	if config.Setting("get", "dev", "Spec.Status", "") == "active" {
		apiIP := config.Setting("get", "dev", "Spec.API.Address", "")
		if apiIP == "cluster" {
			extAPI = true
		}
	} else {
		extAPI = true
	}
	return extAPI
}

func (f F) EndpointCache(method, module, value string) string {
	return config.Setting(method, "infrastructure", "Spec."+strings.Title(module)+".Cache.Endpoint", value)
}

func (f F) generateAPIEndpoint(apiIP string, containerPorts [][]string) string {
	return "https://" + apiIP + ":" + containerPorts[0][0]
}

var APIStatusPhaseRunning bool = false

func (f F) RequestHandler(context, method, route, cookie, contentType, queryParams string, body io.Reader, auth bool) string {
	var response string
	var success bool = false
	var waitingCounter int = 0
	var loggedWaiting bool = false
	for ok := true; ok; ok = !success {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWaiting}, "Sending "+vars.NeuraKubeName+" request to route "+route+"..", 2)
		logging.ProgressSpinner("start")
		responseReq, err := f.sendRequest(context, method, route, cookie, contentType, queryParams, body, auth)
		logging.ProgressSpinner("stop")
		if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, false) {
			if waitingCounter == 5 {
				logging.Log([]string{"\n", vars.EmojiAPI, vars.EmojiWarning}, "The response from "+vars.NeuraKubeName+" takes longer than usual.", 0)
				if config.Setting("get", "dev", "Spec.API.Address", "") != "localhost" {
					pods.F.Logs(pods.F{}, api.F.GetNamespace(api.F{}), api.F.GetContext(api.F{}), "", true, 2)
				} else {
					logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWarning}, "Your local "+vars.NeuraKubeName+" instance seems not to be running.", 0)
					logging.Log([]string{"", vars.EmojiAPI, vars.EmojiInfo}, "Start it with: "+vars.NeuraKubeNameRepo+" server", 0)
				}
			} else {
				if !APIStatusPhaseRunning && extAPI {
					pods.F.WaitForPhase(pods.F{}, api.F.GetNamespace(api.F{}), api.F.GetContext(api.F{}), "running", 2)
					APIStatusPhaseRunning = true
				}
				if !loggedWaiting {
					logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWaiting}, "Waiting for response from "+vars.NeuraKubeName+"..", 0)
					loggedWaiting = true
				}
				timing.TimeOut(1, "s")
			}
			waitingCounter++
		} else {
			success = true
			pods.InterruptPodLogsLive = true
			logging.ProgressSpinner("stop")
			response = f.readResponse(route, responseReq)
		}
	}
	return response
}

func (f F) sendRequest(context, method, route, cookie, contentType, queryParams string, body io.Reader, auth bool) (*http.Response, error) {
	client := f.getClient()
	if queryParams != "" {
		queryParams = "?" + queryParams
	}
	requestRoute := vars.RESTRoutePreamble + "/" + route + queryParams
	req, err := http.NewRequest(method, f.getAPIEndpoint(context)+requestRoute, body)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create new http request!", false, false, true)
	req.Header.Set("From", users.GetIDActive())
	validToken, err := jwt.GenerateToken(jwt.SigningKeyActive)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create JWT token!", false, true, true)
	req.Header.Set("Token", validToken)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return client.Do(req)
}

var loggedResponseStatus bool = false

func (f F) readResponse(route string, res *http.Response) string {
	response, err := ioutil.ReadAll(res.Body)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to read repsonse!", false, false, true)
	respString := string(response)
	if respString == "signature is invalid" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "User not authorized (JWT signing key invalid)!", true, true, true)
	} else if respString == "404 page not found" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "The requested route is unknown: "+route, true, true, true)
	} else {
		if !loggedResponseStatus {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Connected to "+vars.NeuraKubeName+".\n", 0)
			loggedResponseStatus = true
		}
	}
	res.Body.Close()
	return respString
}

func (f F) getClient() *http.Client {
	tr := &http.Transport{
		IdleConnTimeout: timing.GetTimeDuration(10, "m"),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}
