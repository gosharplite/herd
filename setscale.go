package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosharplite/herd/log"
	"io/ioutil"
	"net/http"
)

type jScale struct {
	ClusterId       string `json:"cluster_id"`
	EnableAutoScale int64  `json:"enable_auto_scale"`
	CpuMin          int64  `json:"cpu_min"`
	CpuMax          int64  `json:"cpu_max"`
	PodMin          int64  `json:"pod_min"`
	PodMax          int64  `json:"pod_max"`
}

func setScaleHandler(w http.ResponseWriter, r *http.Request) {

	// request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, log.Err("ioutil.ReadAll: %v", err))
		return
	}

	log.Info("setscale request body:\n%v\n", string(body))

	var jss jScale
	err = json.Unmarshal(body, &jss)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("json.Unmarshal: %v", err))
		return
	}

	// write into etcd
	err = etcd.Set(jss.ClusterId, string(body))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, log.Err("etcd.Set: %v", err))
		return
	}
}
