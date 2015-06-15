package main

import (
	e "github.com/gosharplite/herd/etcd"
	"testing"
)

func TestAutoScale(t *testing.T) {
	t.Skip()
	etcd = e.NewClient("http://123.51.216.20:2379", "/gigacloud.com/autoscale/")

	err := fetchScales()
	if err != nil {
		t.Fail()
	}
}
