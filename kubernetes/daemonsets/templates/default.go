package templates

import (
	"../../../errors"
	"../../../filesystem"
	"../../../io"
	"../../../runtime"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func GetConfig(platform, provider, id, nodeImage string) *appsv1.DaemonSet {
	basePath := "../tools-go/kubernetes/daemonsets/templates/"
	templatePath := basePath + platform + "/" + provider + "/" + id + "/" + nodeImage
	templateFile := "/template.yaml"
	templateFilePath := templatePath + templateFile
	if !filesystem.Exists(templatePath) {
		filesystem.CreateDir(templatePath, false)
		templateURL := "https://raw.githubusercontent.com/GoogleCloudPlatform/container-engine-accelerators/master/nvidia-driver-installer/cos/daemonset-preloaded.yaml"
		io.F.DownloadFile(io.F{}, templateFilePath, templateURL)
		//errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "There are no daemonset templates that match your provided arguments!", true, true, true)
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(filesystem.FileToBytes(templateFilePath), nil, nil)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, true, true)
	ds := obj.(*appsv1.DaemonSet)
	return ds
}
