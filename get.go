package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosharplite/herd/k8s"
	"github.com/gosharplite/herd/log"
	"io/ioutil"
	"net/http"
)

type req struct {
	Services []string `json:"services"`
	Clusters []string `json:"clusters"`
}

type pod struct {
	ContainerName string `json:"container_name"`
	Cpu           int64  `json:"cpu"`
	Mem           int64  `json:"mem"`
}

type rc struct {
	ClusterName string `json:"cluster_name"`
	Containers  []pod  `json:"containers"`
}

type service struct {
	ServiceName string `json:"service_name"`
	Clusters    []rc   `json:"clusters"`
}

type resp struct {
	Services []service `json:"services"`
	Clusters []rc      `json:"clusters"`
}

func getHandler(w http.ResponseWriter, r *http.Request) {

	// request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("ioutil.ReadAll: %v", err))
		return
	}

	var rq req
	err = json.Unmarshal(body, &rq)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("json.Unmarshal: %v", err))
		return
	}

	// response
	//	rsp := resp{
	//		Services: make([]service, 0),
	//		Clusters: make([]rc, 0),
	//	}

	// clusters
	for _, rc := range rq.Clusters {

		// get rc
		repcon, err := k8s.GetRC(rc)
		if err != nil {
			log.Err(" k8s.GetRC(): %v", err)
			continue
		}

		log.Info("rc, labels: %v, %v", repcon.Name, repcon.Labels)

		// get pods in rc
		pods, err := k8s.GetPods(repcon.Labels)
		if err != nil {
			log.Err(" k8s.GetPods(): %v", err)
			continue
		}

		for _, p := range pods.Items {
			log.Info("pod: %v", p.Name)
		}

		// get dockers in pod
		// get cadvisor of docker
		// cal cpu/mem of pod

	}

}
