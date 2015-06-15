package etcd

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"github.com/gosharplite/herd/log"
	"strings"
)

type Client struct {
	etcd_machines []string
	etcd_prefix   string
	clt           *etcd.Client
}

func NewClient(etcd_machines, etcd_prefix string) *Client {

	c := Client{
		etcd_machines: strings.Split(etcd_machines, ","),
		etcd_prefix:   etcd_prefix,
	}

	c.clt = etcd.NewClient(c.etcd_machines)
	if c.clt == nil {
		log.Err("etcd.NewClient()")
		return nil
	}

	return &c
}

func (c *Client) Set(key, value string) error {

	if key == "" {
		log.Err("key is empty")
		return errors.New("key is empty")
	}

	_, err := c.clt.Set(c.etcd_prefix+key, value, 0)
	if err != nil {
		log.Err("c.Set(): %v", err)
		return err
	}

	return nil
}

func (c *Client) Get(key string) (string, error) {

	r, err := c.clt.Get(c.etcd_prefix+key, false, false)
	if err != nil {
		log.Err("c.Get(): %v", err)
		return "", err
	}

	return r.Node.Value, nil
}

func (c *Client) GetScales() (etcd.Nodes, error) {

	r, err := c.clt.Get(c.etcd_prefix, false, false)
	if err != nil {
		log.Err("c.Get(): %v", err)
		return nil, err
	}

	return r.Node.Nodes, nil
}

func (c *Client) Delete(key string) error {

	_, err := c.clt.Delete(c.etcd_prefix+key, false)
	if err != nil {
		log.Err("(): %v", err)
		return err
	}

	return nil
}
