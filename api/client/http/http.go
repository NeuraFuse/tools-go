package http

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/neurafuse/tools-go/ci"
	"github.com/neurafuse/tools-go/ci/api"
	"github.com/neurafuse/tools-go/ci/base"
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/crypto/jwt"
	"github.com/neurafuse/tools-go/errors"
	toolsIO "github.com/neurafuse/tools-go/io"
	"github.com/neurafuse/tools-go/kubernetes/pods"
	"github.com/neurafuse/tools-go/kubernetes/services"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/timing"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var extAPI bool
var loggedAPIEndpoint bool

func (f F) getAPIEndpoint(context, requestRoute string) string {
	var apiEndpoint string
	switch context {
	case vars.NeuraKubeNameID:
		apiEndpoint = f.endpointNeuraKube()
	case "inference/gpt":
		apiEndpoint = f.endpointGPT()
	}
	if !loggedAPIEndpoint {
		logging.Log([]string{"", vars.EmojiLink, vars.EmojiInfo}, "Endpoint: "+apiEndpoint, 0)
		loggedAPIEndpoint = true
	}
	return apiEndpoint + requestRoute
}

func (f F) endpointGPT() string {
	var apiIP string
	apiIP = services.F.GetClusterIP(services.F{}, base.F.GetNamespace(base.F{}), ci.F.GetContextID(ci.F{}))
	return f.generateAPIEndpoint(apiIP, base.F.GetContainerPorts(base.F{}, "inference"))
}

func (f F) endpointNeuraKube() string {
	var apiIP string
	if config.APILocationCluster() {
		var checkSetup bool
		var endpointCacheIP string = f.EndpointCache("get", vars.NeuraKubeNameID, "")
		if endpointCacheIP == "" {
			checkSetup = true
		} else {
			apiIP = f.generateAPIEndpoint(endpointCacheIP, api.F.GetContainerPorts(api.F{}))
			if toolsIO.F.Reachable(toolsIO.F{}, apiIP) {
				apiIP = endpointCacheIP
			} else {
				checkSetup = true
			}
		}
		if checkSetup {
			api.F.EvalAction(api.F{}, "connect")
			apiIP = services.F.GetLoadBalancerIP(services.F{}, api.F.GetNamespace(api.F{}), ci.F.GetContextID(ci.F{}))
			f.EndpointCache("set", vars.NeuraKubeNameID, apiIP)
		}
	} else {
		apiIP = "localhost"
	}
	return f.generateAPIEndpoint(apiIP, api.F.GetContainerPorts(api.F{}))
}

func (f F) EndpointCache(method, module, value string) string {
	return config.Setting(method, "infrastructure", "Spec."+strings.Title(module)+".Cache.Endpoint", value)
}

func (f F) generateAPIEndpoint(apiIP string, containerPorts [][]string) string {
	return "https://" + apiIP + ":" + containerPorts[0][0]
}

var APIStatusPhaseRunning bool

func (f F) RequestHandler(context, method, route, cookie, contentType, queryParams string, body io.Reader, auth bool) string {
	var response string
	var success bool
	var waitingCounter int
	var loggedWaiting bool
	for ok := true; ok; ok = !success {
		logging.ProgressSpinner("start")
		responseReq, err := f.sendRequest(context, method, route, cookie, contentType, queryParams, body, auth)
		logging.ProgressSpinner("stop")
		if errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, false) {
			if waitingCounter >= 2 {
				if waitingCounter == 2 {
					logging.Log([]string{"\n", vars.EmojiAPI, vars.EmojiWarning}, "The response from "+vars.NeuraKubeName+" takes longer than usual.", 0)
				}
				if config.APILocationCluster() {
					f.EndpointCache("reset", vars.NeuraKubeNameID, "")
					if waitingCounter == 3 {
						pods.F.Logs(pods.F{}, api.F.GetNamespace(api.F{}), api.F.GetContext(api.F{}), "", true, 2)
					}
				} else {
					logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "Your local "+vars.NeuraKubeName+" instance seems not to be running.", 0)
					logging.Log([]string{"", vars.EmojiAPI, vars.EmojiRocket}, "Start it in another terminal with: "+vars.NeuraKubeNameID+" server\n", 0)
				}
				config.Setting("set", "infrastructure", "Spec.Neurakube.Cache.ConnectStatus", "false")
			} else {
				if !APIStatusPhaseRunning && extAPI {
					pods.F.WaitForPhase(pods.F{}, api.F.GetNamespace(api.F{}), api.F.GetContext(api.F{}), "running", 2)
					APIStatusPhaseRunning = true
				}
				if !loggedWaiting {
					logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWaiting}, "Waiting for response from "+vars.NeuraKubeName+"..", 0)
					loggedWaiting = true
				}
				timing.Sleep(500, "ms")
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
	var client *http.Client = f.getClient()
	if queryParams != "" {
		queryParams = "?" + queryParams
	}
	var requestRoute string = vars.RESTRoutePreamble + "/" + route + queryParams
	var apiEndpoint string = f.getAPIEndpoint(context, requestRoute)
	logging.Log([]string{"", vars.EmojiClient, vars.EmojiAPI}, "apiEndpoint: "+apiEndpoint, 2)
	req, err := http.NewRequest(method, apiEndpoint, body)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create new http request!", false, false, true)
	req.Header.Set("From", usersID.F.GetActive(usersID.F{}))
	validToken, err := jwt.GenerateToken(jwt.SigningKeyActive)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create JWT token!", false, true, true)
	req.Header.Set("Token", validToken)
	if cookie != "" {
		logging.Log([]string{"", vars.EmojiClient, vars.EmojiAPI}, "Cookie: "+cookie, 2)
		req.Header.Set("Cookie", cookie)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return client.Do(req)
}

var loggedResponseStatus bool

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
