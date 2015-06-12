package etcd

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"github.com/gosharplite/herd/log"
)

// TODO Make it better.

var (
	ETCD_MACHINES = []string{
		"http://192.168.3.36:2379",
		"http://192.168.3.37:2379",
		"http://192.168.3.38:2379"}

	ETCD_PREFIX = "/gigacloud.com/autoscale/"

	c *etcd.Client
)

func init() {
	c = etcd.NewClient(ETCD_MACHINES)
	if c == nil {
		log.Err("etcd.NewClient()")
	}
}

func Set(key, value string) error {

	if key == "" {
		log.Err("key is empty")
		return errors.New("key is empty")
	}

	_, err := c.Set(ETCD_PREFIX+key, value, 0)
	if err != nil {
		log.Err("c.Set(): %v", err)
		return err
	}

	return nil
}

func Get(key string) (string, error) {

	if key == "" {
		log.Err("key is empty")
		return "", errors.New("key is empty")
	}

	r, err := c.Get(ETCD_PREFIX+key, false, false)
	if err != nil {
		log.Err("c.Get(): %v", err)
		return "", err
	}

	return r.Node.Value, nil
}
