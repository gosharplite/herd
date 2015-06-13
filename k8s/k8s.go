package k8s

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/gosharplite/herd/log"
)

// TODO Make it better.

var (
	K8S_HOST      = "http://192.168.4.54:8080"
	K8S_VERSION   = "v1beta3"
	K8S_NAMESPACE = "default"
	c             *client.Client
)

func init() {

	config := &client.Config{
		Host:    K8S_HOST,
		Version: K8S_VERSION,
	}

	var err error
	c, err = client.New(config)
	if err != nil {
		log.Err("client.New()")
	}
}

func GetService(name string) (*api.Service, error) {
	return c.Services(K8S_NAMESPACE).Get(name)
}

func GetRCList(lbs map[string]string) (*api.ReplicationControllerList, error) {
	selector := labels.Set(lbs).AsSelector()
	log.Info("selector: %v", selector)

	return c.ReplicationControllers(K8S_NAMESPACE).List(selector)
}

func GetRC(name string) (*api.ReplicationController, error) {
	return c.ReplicationControllers(K8S_NAMESPACE).Get(name)
}

func GetPods(lbs map[string]string) (*api.PodList, error) {
	selector := labels.Set(lbs).AsSelector()
	return c.Pods(K8S_NAMESPACE).List(selector, fields.Everything())
}
