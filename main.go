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

type jSetScale struct {
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

	var jss jSetScale
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

func getScaleHandler(w http.ResponseWriter, r *http.Request) {

	type jReq struct {
		Clusters []string `json:"clusters"`
	}

	type jResp struct {
		Clusters []jSetScale `json:"clusters"`
	}

	// request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("ioutil.ReadAll: %v", err))
		return
	}

	var jr jReq
	err = json.Unmarshal(body, &jr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("json.Unmarshal: %v", err))
		return
	}

	// response
	scales := make([]jSetScale, 0)
	for _, item := range jr.Clusters {
		v, err := etcd.Get(item)
		if err != nil {
			log.Err("etcd.Get: %v", err)
			continue
		}

		var jss jSetScale
		err = json.Unmarshal([]byte(v), &jss)
		if err != nil {
			log.Err("json.Unmarshal: %v", err)
			continue
		}

		scales = append(scales, jss)
	}

	jresp := jResp{Clusters: scales}

	j, err := json.Marshal(jresp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, log.Err("json.Marshal: %v", err))
		return
	}

	fmt.Fprintf(w, string(j))
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	//	test_3_Handler(w, r)
}
