package main

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"net/http"
	"strconv"
	"time"
)

func test_3_Handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "herd test,"+r.Host+","+strconv.FormatInt(time.Now().UnixNano(), 10)+"\n")

	machines := []string{"http://192.168.3.36:2379", "http://192.168.3.37:2379", "http://192.168.3.38:2379"}
	client := etcd.NewClient(machines)

	result, err := client.Get("/", true, true)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "client.Get: %v", err)
		return
	}

	b, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "json.Marshal(receivedPodList): %v", err)
		return
	}

	fmt.Fprintf(w, "etcd get:\n%v\n", string(b))
	fmt.Printf("etcd get:\n%v\n", string(b))
}
