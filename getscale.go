package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosharplite/herd/log"
	"io/ioutil"
	"net/http"
)

func getScaleHandler(w http.ResponseWriter, r *http.Request) {

	type jReq struct {
		Clusters []string `json:"clusters"`
	}

	type jResp struct {
		Clusters []jScale `json:"clusters"`
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
	scales := make([]jScale, 0)
	for _, item := range jr.Clusters {
		v, err := etcd.Get(item)
		if err != nil {
			log.Err("etcd.Get: %v", err)
			continue
		}

		var jss jScale
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
