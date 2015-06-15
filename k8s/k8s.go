package k8s

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	clt "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/gosharplite/herd/log"
)

type Client struct {
	host      string
	version   string
	namespace string
	clt       *clt.Client
}

func NewClient(host, version, namespace string) *Client {

	c := Client{
		host:      host,
		version:   version,
		namespace: namespace,
	}

	config := &clt.Config{
		Host:    host,
		Version: version,
	}

	var err error
	c.clt, err = clt.New(config)
	if err != nil {
		log.Err("client.New()")
		return nil
	}

	return &c
}

func (c *Client) GetService(name string) (*api.Service, error) {
	return c.clt.Services(c.namespace).Get(name)
}

func (c *Client) GetRCList(lbs map[string]string) (*api.ReplicationControllerList, error) {
	selector := labels.Set(lbs).AsSelector()
	log.Info("selector: %v", selector)

	return c.clt.ReplicationControllers(c.namespace).List(selector)
}

func (c *Client) GetRC(name string) (*api.ReplicationController, error) {
	return c.clt.ReplicationControllers(c.namespace).Get(name)
}

func (c *Client) GetPods(lbs map[string]string) (*api.PodList, error) {
	selector := labels.Set(lbs).AsSelector()
	return c.clt.Pods(c.namespace).List(selector, fields.Everything())
}

func (c *Client) Update(ctrl *api.ReplicationController) (*api.ReplicationController, error) {
	return c.clt.ReplicationControllers(c.namespace).Update(ctrl)
}
