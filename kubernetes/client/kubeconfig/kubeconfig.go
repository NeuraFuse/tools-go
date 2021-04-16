package kubeconfig

import (
	"context"

	client ".."
	"../../../config"
	"../../../errors"
	"../../../filesystem"
	"../../../logging"
	"../../../runtime"
	"../../../vars"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type F struct{}

func (f F) Create(namespace string) { // TODO: Test compatibility for selfhosted clusters
	logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiCrypto}, "Creating kubeconfig..", 0)
	secret, err := f.getSecret()
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "There are not existing secrets for the namespace "+namespace+" to create a kubeconfig!", false, false, true) {
		clusters, clusterName := f.getClusters(secret)
		contexts := f.getContexts(clusterName, namespace)
		authinfos := f.getAuthInfos(namespace, secret)
		clientConfig := f.getClientConfig(clusters, contexts, authinfos)
		f.save(clientConfig)
	}
}

func (f F) getSecret() (apiv1.Secret, error) {
	secretList, err := client.F.GetAuth(client.F{}).CoreV1().Secrets("kube-system").List(context.TODO(), metav1.ListOptions{})
	var secret apiv1.Secret
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to fetch secrets for kubeconfig creation!", false, false, true) {
		if len(secretList.Items) >= 2 {
			secret = secretList.Items[1]
			logging.Log([]string{"", vars.EmojiCrypto, vars.EmojiProcess}, "Using secret: "+secret.ObjectMeta.Name, 0)
		} else {
			err = errors.New("")
		}
	}
	return secret, err
}

func (f F) getClusters(secret apiv1.Secret) (map[string]*clientcmdapi.Cluster, string) {
	clusters := make(map[string]*clientcmdapi.Cluster)
	clusterName := config.Setting("get", "infrastructure", "Spec.Cluster.Name", "")
	clusters[clusterName] = &clientcmdapi.Cluster{
		Server:                   "https://" + client.F.GetRestConfig(client.F{}).Host,
		CertificateAuthorityData: secret.Data["ca.crt"],
	}
	return clusters, clusterName
}

func (f F) getContexts(clusterName, namespace string) map[string]*clientcmdapi.Context {
	contexts := make(map[string]*clientcmdapi.Context)
	contexts["default-context"] = &clientcmdapi.Context{
		Cluster:   clusterName,
		Namespace: namespace,
		AuthInfo:  namespace,
	}
	return contexts
}

func (f F) getAuthInfos(namespace string, secret apiv1.Secret) map[string]*clientcmdapi.AuthInfo {
	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	if vars.InfraProviderActive == vars.InfraProviderGcloud {
		authinfos[namespace] = &clientcmdapi.AuthInfo{
			//Token: string(secret.Data["token"]),
			AuthProvider: &clientcmdapi.AuthProviderConfig{
				Name:   "gcp",
				Config: f.getAuthProviderConfig(),
			},
		}
	} else if vars.InfraProviderActive == vars.InfraProviderSelfHosted {
		authinfos[namespace] = &clientcmdapi.AuthInfo{
			Token: string(secret.Data["token"]),
		}
	}
	return authinfos
}

func (f F) getAuthProviderConfig() map[string]string {
	config := make(map[string]string)
	if vars.InfraProviderActive == vars.InfraProviderGcloud {
		config["cmd-args"] = "config config-helper --format=json"
		config["cmd-path"] = "/Users/djw/google-cloud-sdk/bin/gcloud"
		config["expiry-key"] = "{.credential.token_expiry}"
		config["token-key"] = "{.credential.access_token}"
	}
	return config
}

func (f F) getClientConfig(clusters map[string]*clientcmdapi.Cluster, contexts map[string]*clientcmdapi.Context, authinfos map[string]*clientcmdapi.AuthInfo) clientcmdapi.Config {
	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: "default-context",
		AuthInfos:      authinfos,
	}
	return clientConfig
}

func (f F) save(clientConfig clientcmdapi.Config) {
	configBasePath := f.GetConfigBasePath()
	configFilePath := configBasePath + "config"
	configFilePathBackup := configFilePath + "_backup"
	if !filesystem.Exists(configBasePath) {
		filesystem.CreateDir(configBasePath, false)
	} else if filesystem.Exists(configFilePath) {
		filesystem.Move(configFilePath, configFilePathBackup, false)
		logging.Log([]string{"", vars.EmojiCrypto, vars.EmojiSuccess}, "Backed up existing kubeconfig: "+configFilePathBackup, 0)
	}
	clientcmd.WriteToFile(clientConfig, configFilePath)
	logging.Log([]string{"", vars.EmojiCrypto, vars.EmojiSuccess}, "Created new kubeconfig: "+configFilePath, 0)
}

func (f F) GetConfigBasePath() string {
	return filesystem.GetUserHomeDir() + "/.kube/"
}