package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gosharplite/herd/etcd"
	"github.com/gosharplite/herd/log"
	"io/ioutil"
	"net/http"
)

var (
	port = flag.Int("port", 8090, "The server port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/setscale", setScaleHandler)
	http.HandleFunc("/getscale", getScaleHandler)
	http.HandleFunc("/test", testHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	fmt.Printf("proxy: %v\n", err)
}

func setScaleHandler(w http.ResponseWriter, r *http.Request) {

	type jSetScale struct {
		ClusterId       string `json:"cluster_id"`
		EnableAutoScale int64  `json:"enable_auto_scale"`
		CpuMin          int64  `json:"cpu_min"`
		CpuMax          int64  `json:"cpu_max"`
		PodMin          int64  `json:"pod_min"`
		PodMax          int64  `json:"pod_max"`
	}

	// request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, log.Err("ioutil.ReadAll: %v", err))
		return
	}

	log.Info("setscale request body:\n%v\n", string(body))

	var jss jSetScale
	err = json.Unmarshal(body, &jss)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("json.Unmarshal: %v", err))
		return
	}

	// echo request
	j, err := json.Marshal(jss)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("json.Marshal: %v", err))
		return
	}

	// write into etcd
	err = etcd.Set(jss.ClusterId, string(j))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, log.Err("etcd.Set: %v", err))
		return
	}
}

func getScaleHandler(w http.ResponseWriter, r *http.Request) {

	// body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("ioutil.ReadAll: %v", err))
		return
	}

	fmt.Fprintf(w, log.Info("getscale body:\n%v\n", string(body)))
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	//	test_3_Handler(w, r)
}
