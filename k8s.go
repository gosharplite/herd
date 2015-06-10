package main

import (
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func getRedisMaster() (*api.PodList, error) {

	// k8s client package
	config := &client.Config{
		Host:    "http://192.168.4.54:8080",
		Version: "v1beta3",
	}

	client, err := client.New(config)
	if err != nil {
		return nil, err
	}

	selector := labels.Set{"name": "redis-master"}.AsSelector()
	receivedPodList, err := client.Pods("default").List(selector, fields.Everything())
	if err != nil {
		return nil, err
	}

	return receivedPodList, nil
}

func getApi() (string, error) {

	resp, err := http.Get("http://192.168.4.54:8080/api")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getVersion() (string, error) {

	resp, err := http.Get("http://192.168.4.54:8080/version")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func test_1_Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "herd test,"+r.Host+","+strconv.FormatInt(time.Now().UnixNano(), 10)+"\n")

	receivedPodList, err := getRedisMaster()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "getRedisMaster(): %v", err)
		return
	}

	b, err := json.MarshalIndent(receivedPodList, "", "    ")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "json.Marshal(receivedPodList): %v", err)
		return
	}

	fmt.Fprintf(w, "PodList:\n%v\n", string(b))
	fmt.Printf("PodList:\n%v\n", string(b))
}
