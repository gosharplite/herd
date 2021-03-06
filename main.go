package main

import (
	"flag"
	"fmt"
	e "github.com/gosharplite/herd/etcd"
	k "github.com/gosharplite/herd/k8s"
	"github.com/gosharplite/herd/log"
	"net/http"
)

var (
	PORT = flag.Int("port", 8090, "The server port")

	K8S_HOST      = flag.String("k8s_host", "http://192.168.4.54:8080", "k8s host")
	K8S_VERSION   = flag.String("k8s_version", "v1", "k8s version")
	K8S_NAMESPACE = flag.String("k8s_namespace", "default", "k8s namespace")
	k8s           *k.Client

	ETCD_MACHINES = flag.String(
		"etcd_machines",
		"http://192.168.3.36:2379,http://192.168.3.37:2379,http://192.168.3.38:2379",
		"etcd machines")
	ETCD_PREFIX = flag.String("etcd_prefix", "/gigacloud.com/autoscale/", "etcd prefix")
	etcd        *e.Client
)

func main() {
	flag.Parse()
	log.Info("K8S_HOST: %v", *K8S_HOST)
	log.Info("K8S_VERSION: %v", *K8S_VERSION)
	log.Info("K8S_NAMESPACE: %v", *K8S_NAMESPACE)
	log.Info("ETCD_MACHINES: %v", *ETCD_MACHINES)
	log.Info("ETCD_PREFIX: %v", *ETCD_PREFIX)

	k8s = k.NewClient(*K8S_HOST, *K8S_VERSION, *K8S_NAMESPACE)
	etcd = e.NewClient(*ETCD_MACHINES, *ETCD_PREFIX)

	http.HandleFunc("/setscale", setScaleHandler)
	http.HandleFunc("/getscale", getScaleHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/getevent", getEventHandler)

	go autoLoop()

	err := http.ListenAndServe(fmt.Sprintf(":%d", *PORT), nil)
	fmt.Printf("proxy: %v\n", err)
}
