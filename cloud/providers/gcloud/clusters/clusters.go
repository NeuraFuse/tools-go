package clusters

import (
	"strconv"

	"github.com/neurafuse/tools-go/cloud/providers/gcloud/clients"
	gcloudConfig "github.com/neurafuse/tools-go/cloud/providers/gcloud/config"
	"github.com/neurafuse/tools-go/config"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"k8s.io/client-go/rest"
)

type F struct{}

func (f F) Get(logResult bool) ([]*containerpb.Cluster, error) {
	var err error
	var resp *containerpb.ListClustersResponse
	var success bool
	for ok := true; ok; ok = !success {
		ctx, client := clients.F.GetContainer(clients.F{})
		projectID := config.Setting("get", "infrastructure", "Spec.Gcloud.ProjectID", "")
		zone := config.Setting("get", "infrastructure", "Spec.Gcloud.Zone", "")
		var request = &containerpb.ListClustersRequest{
			ProjectId: projectID,
			Zone:      zone,
		}
		resp, err = client.ListClusters(ctx, request)
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get clusters!", false, false, true) && resp != nil {
			success = true
		} else {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWarning}, "Trying to recover..", 0)
			logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "ProjectID: "+projectID+" / Zone: "+zone+"\n", 0)
			timing.Sleep(1, "s")
		}
	}
	logging.ProgressSpinner("stop")
	var clusters []*containerpb.Cluster
	if resp != nil {
		clusters = resp.Clusters
	}
	if logResult {
		logging.Log([]string{"", vars.EmojiKubernetes, ""}, "Clusters in your setup:\n", 0)
		if clusters != nil {
			for iC, cluster := range clusters {
				var linebreak string
				if iC == len(clusters)-1 {
					linebreak = "\n"
				}
				logging.Log([]string{"", "", ""}, "["+strconv.Itoa(iC)+"] "+cluster.Name+linebreak, 0)
			}
		} else {
			logging.Log([]string{"", vars.EmojiKubernetes, ""}, "There are no clusters deployed.\n", 0)
		}
	}
	return clusters, err
}

func (f F) Create() bool {
	var opSuccess bool
	var clusterID string = config.Setting("get", "infrastructure", "Spec.Cluster.ID", "")
	exists, _ := f.Exists()
	if !exists {
		ctx, client := clients.F.GetContainer(clients.F{})
		request := &containerpb.CreateClusterRequest{
			ProjectId: config.Setting("get", "infrastructure", "Spec.Gcloud.ProjectID", ""),
			Zone:      config.Setting("get", "infrastructure", "Spec.Gcloud.Zone", ""),
			Cluster:   gcloudConfig.F.ClusterConfig(gcloudConfig.F{}, config.Setting("get", "infrastructure", "Spec.Gcloud.MachineType", ""), ""),
		}
		var op *containerpb.Operation
		var err error
		op, err = client.CreateCluster(ctx, request)
		var opName string = op.Name
		if op.Status == 1 {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Started creation of new cluster "+clusterID+"..", 0)
		} else {
			if err == nil {
				err = errors.New("Creation operation status is not 1!")
			}
		}
		var errMsg string = "Unable to fulfill the creation operation for cluster " + clusterID + "!"
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, true, true) {
			var checkSuccess bool
			var loggedRunning bool
			for ok := true; ok; ok = !checkSuccess {
				request := &containerpb.GetOperationRequest{
					Name: opName,
				}
				op, err = client.GetOperation(ctx, request)
				if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get operation for readiness check!", false, true, true) {
					if op.Status == 2 {
						if !loggedRunning {
							logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Creation operation for new cluster "+clusterID+" is running..", 0)
							loggedRunning = true
						}
					} else if op.Status == 3 {
						opSuccess = true
						checkSuccess = true
						logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "The cluster "+clusterID+" is now fully created.", 0)
					} else if op.Status == 4 {
						errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "The cluster creation operation is aborting!", true, true, true)
					}
				}
				timing.Sleep(1, "s")
			}
		}
	} else {
		opSuccess = true
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "The cluster "+clusterID+" is available.\n", 0)
	}
	logging.ProgressSpinner("stop")
	return opSuccess
}

func (f F) Delete() {
	var clusterID string = f.getClusterID()
	exists, _ := f.Exists()
	if exists {
		ctx, client := clients.F.GetContainer(clients.F{})
		request := &containerpb.DeleteClusterRequest{
			ProjectId: config.Setting("get", "infrastructure", "Spec.Gcloud.ProjectID", ""),
			Zone:      config.Setting("get", "infrastructure", "Spec.Gcloud.Zone", ""),
			ClusterId: config.Setting("get", "infrastructure", "Spec.Cluster.ID", ""),
			Name:      config.Setting("get", "infrastructure", "Spec.Cluster.ID", ""),
		}
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiWarning}, "Deleting cluster "+clusterID+"!", 0)
		_, err := client.DeleteCluster(ctx, request)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete cluster "+clusterID+"!", false, false, true)
		logging.ProgressSpinner("start")
		var success bool
		for ok := true; ok; ok = !success {
			if f.getStatus(clusterID, true) == 4 {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWaiting}, "Waiting for the cluster "+clusterID+" to be deleted..", 0)
				timing.Sleep(1, "s")
			} else {
				logging.ProgressSpinner("stop")
				logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Cluster deleted.\n", 0)
				success = true
			}
		}
	} else {
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiSuccess}, "The cluster "+clusterID+" is already deleted.", 0)
	}
}

func (f F) getClusterID() string {
	var clusterID string
	clusterID = config.Setting("get", "infrastructure", "Spec.Cluster.ID", "")
	return clusterID
}

func (f F) Exists() (bool, error) {
	var exists bool
	var clusters []*containerpb.Cluster
	var err error
	clusters, err = f.Get(false)
	if len(clusters) > 0 {
		for _, cluster := range clusters {
			if cluster.Name == f.getClusterID() {
				exists = true
				break
			}
		}
	}
	return exists, err
}

func (f F) getStatus(id string, ignoreDoesNotExist bool) containerpb.Cluster_Status {
	cluster := f.getCluster(id, ignoreDoesNotExist)
	var status containerpb.Cluster_Status = 0
	if cluster != nil {
		status = cluster.Status
	}
	return status
}

func (f F) getCluster(id string, ignoreDoesNotExist bool) *containerpb.Cluster {
	var clusterSelection *containerpb.Cluster
	clusters, err := f.Get(false)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true) {
		for _, cluster := range clusters {
			if cluster.Name == id {
				clusterSelection = cluster
				break
			}
		}
	} else if !ignoreDoesNotExist {
		f.ResourceMissing()
	}
	return clusterSelection
}

func (f F) ResourceMissing() bool {
	var clusterID string = config.Setting("get", "infrastructure", "Spec.Cluster.ID", "")
	var devConfigSkip bool
	if !infraConfig.F.GetClusterRecentlyDeleted(infraConfig.F{}) {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWarning}, "The pre-configured cluster "+clusterID+" does not exist yet.", 0)
		var askCreation bool
		if !config.DevConfigActive() {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWarning}, "The process can't be continued without an existing cluster.", 0)
			askCreation = true
		} else {
			var action string = terminal.GetUserSelection("Do you want to continue by aborting the overlay process?", []string{}, false, true)
			if action == "Yes" {
				devConfigSkip = true
			} else {
				askCreation = true
			}
		}
		if askCreation {
			var action string = terminal.GetUserSelection("Do you want to create the cluster?", []string{}, false, true)
			if action == "Yes" {
				f.Create()
			} else {
				terminal.Exit(0, "")
			}
		}
	} else {
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiSuccess}, "The specified cluster "+clusterID+" was successfully deleted.", 0)
		terminal.Exit(0, "")
	}
	return devConfigSkip
}

var loggedConnected bool

func (f F) GetAuthConfig() *rest.Config {
	var success bool
	var cluster *containerpb.Cluster
	logging.ProgressSpinner("start")
	for ok := true; ok; ok = !success {
		cluster = f.getCluster(config.Setting("get", "infrastructure", "Spec.Cluster.ID", ""), false)
		if f.getStatus(config.Setting("get", "infrastructure", "Spec.Cluster.ID", ""), false) == 0 {
			f.ResourceMissing()
		} else if f.getStatus(config.Setting("get", "infrastructure", "Spec.Cluster.ID", ""), false) != 2 {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWaiting}, "Waiting for cluster "+config.Setting("get", "infrastructure", "Spec.Cluster.ID", "")+" to be ready..", 0)
			timing.Sleep(1, "s")
		} else {
			success = true
			logging.ProgressSpinner("stop")
			if !loggedConnected {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Connected to cluster "+config.Setting("get", "infrastructure", "Spec.Cluster.ID", "")+".", 0)
				loggedConnected = true
			}
		}
	}
	restconfig := &rest.Config{}
	restconfig.Host = cluster.Endpoint
	restconfig.ServerName = config.Setting("get", "infrastructure", "Spec.Cluster.ID", "")
	// restconfig.TLSClientConfig.CAData = []byte(masterAuth.ClusterCaCertificate)
	// restconfig.TLSClientConfig.CertData = []byte(masterAuth.ClientCertificate)
	// restconfig.KeyData = []byte(masterAuth.ClientKey)
	// errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	return restconfig
}
