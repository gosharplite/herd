package main

import (
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/gosharplite/herd/cadv"
	"github.com/gosharplite/herd/log"
	"io/ioutil"
	"net/http"
)

type req struct {
	Services   []string `json:"services"`
	Clusters   []string `json:"clusters"`
	Containers []string `json:"containers"`
}

type pod struct {
	ContainerName string `json:"container_name"`
	Cpu           uint64 `json:"cpu"`
	Mem           uint64 `json:"mem"`
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
	Services   []service `json:"services"`
	Clusters   []rc      `json:"clusters"`
	Containers []pod     `json:"containers"`
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	defer log.Un(log.Trace("getHandler"))

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

	//	response
	rsp := resp{
		Services:   make([]service, 0),
		Clusters:   make([]rc, 0),
		Containers: make([]pod, 0),
	}

	// services
	for _, se := range rq.Services {
		ses, err := getService(se)
		if err != nil {
			log.Err("getService(): %v", err)
			continue
		}

		rsp.Services = append(rsp.Services, ses)
	}

	// clusters
	for _, rc := range rq.Clusters {
		rcs, err := getRc(rc)
		if err != nil {
			log.Err("getRc(): %v", err)
			continue
		}

		rsp.Clusters = append(rsp.Clusters, rcs)
	}

	// containers
	for _, co := range rq.Containers {
		pd, err := k8s.GetPod(co)
		if err != nil {
			log.Err("k8s.GetPod(): %v", err)
			continue
		}

		cos, err := getPod(*pd)
		if err != nil {
			log.Err("getPod(): %v", err)
			continue
		}

		rsp.Containers = append(rsp.Containers, cos)
	}

	// return response
	j, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("json.Marshal: %v", err))
		return
	}

	log.Marshal(rsp)

	fmt.Fprintf(w, string(j))
}

func getService(name string) (service, error) {

	se, err := k8s.GetService(name)
	if err != nil {
		log.Err("k8s.GetService(): %v", err)
		return service{}, err
	}

	//	log.Info("service: %v", se.Name)

	rcList, err := k8s.GetRCList(se.Spec.Selector)
	if err != nil {
		log.Err("k8s.GetRCList(): %v", err)
		return service{}, err
	}

	rspClusters := make([]rc, 0)

	// clusters
	for _, rcl := range rcList.Items {
		rc := rcl.Name
		rcs, err := getRc(rc)
		if err != nil {
			log.Err("getRc(): %v", err)
			continue
		}

		rspClusters = append(rspClusters, rcs)
	}

	return service{
		ServiceName: name,
		Clusters:    rspClusters,
	}, nil
}

func getRc(name string) (rc, error) {

	repcon, err := k8s.GetRC(name)
	if err != nil {
		log.Err("k8s.GetRC(): %v", err)
		return rc{}, err
	}

	//	log.Info("rc: %v", repcon.Name)

	// get pods in rc
	pods, err := k8s.GetPods(repcon.Spec.Selector)
	if err != nil {
		log.Err("k8s.GetPods(): %v", err)
		return rc{}, err
	}

	rcs := rc{
		ClusterName: name,
		Containers:  make([]pod, 0),
	}

	for _, p := range pods.Items {

		pd, err := getPod(p)
		if err != nil {
			log.Err("getPod(): %v", err)
			continue
		}

		rcs.Containers = append(rcs.Containers, pd)
	}

	return rcs, nil
}

func getPod(p api.Pod) (pod, error) {

	//	log.Info("pod: %v", p.Name)

	// get machine info for total mem
	mInfo, err := cadv.GetMInfo(p.Status.HostIP)
	if err != nil {
		log.Err("cadv.GetMInfo(): %v", err)
		return pod{}, err
	}

	var podCpu, podMem uint64

	// get dockers in pod
	for _, d := range p.Status.ContainerStatuses {

		// get cadvisor of docker
		info, err := cadv.GetCInfo(p.Status.HostIP, d.ContainerID)
		if err != nil {
			log.Err("k8s.GetInfo(): %v", err)
			continue
		}

		// cpu, mem
		var timeStart, timeStop int64
		var cpuStart, cpuStop, mem_high uint64
		for i, s := range info.Stats {
			if i == 0 {
				timeStart = s.Timestamp.UnixNano()
				cpuStart = s.Cpu.Usage.Total
			} else {
				timeStop = s.Timestamp.UnixNano()
				cpuStop = s.Cpu.Usage.Total
			}

			if s.Memory.Usage > mem_high {
				mem_high = s.Memory.Usage
			}
		}

		var percent uint64
		if timeStop > timeStart {
			percent = (cpuStop - cpuStart) * 100 / (uint64)(timeStop-timeStart)
		}

		var memPercent uint64
		if mInfo.MemoryCapacity > 0 {
			memPercent = mem_high * 100 / (uint64)(mInfo.MemoryCapacity)
		}

		// cal cpu/mem of pod
		podCpu += percent
		podMem += memPercent
	}

	return pod{
		ContainerName: p.Name,
		Cpu:           podCpu,
		Mem:           podMem,
	}, nil
}
