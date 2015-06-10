package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/cadvisor/client"
	info "github.com/google/cadvisor/info/v1"
	"net/http"
	"strconv"
	"time"
)

func test_2_Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "herd test,"+r.Host+","+strconv.FormatInt(time.Now().UnixNano(), 10)+"\n")

	client, err := client.NewClient("http://192.168.4.54:4194/")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "client.NewClient: %v", err)
		return
	}

	// mInfo
	mInfo, err := client.MachineInfo()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "client.MachineInfo: %v", err)
		return
	}

	b, err := json.MarshalIndent(mInfo, "", "    ")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "json.Marshal(receivedPodList): %v", err)
		return
	}

	//	fmt.Fprintf(w, "MachineInfo:\n%v\n", string(b))
	//	fmt.Printf("MachineInfo:\n%v\n", string(b))

	// DockerContainer
	request := info.ContainerInfoRequest{NumStats: 1}
	sInfo, err := client.DockerContainer("cee96c31361dfe7083eae2930afc9c5cb796fb21d6eeca853c52be2bfe5ce563", &request)

	// cpu
	b, err = json.MarshalIndent(sInfo.Stats[0].Cpu, "", "    ")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "json.Marshal(receivedPodList): %v", err)
		return
	}

	fmt.Fprintf(w, "Cpu:\n%v\n", string(b))
	fmt.Printf("Cpu:\n%v\n", string(b))

	// memory
	b, err = json.MarshalIndent(sInfo.Stats[0].Memory, "", "    ")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "json.Marshal(receivedPodList): %v", err)
		return
	}

	fmt.Fprintf(w, "Memory:\n%v\n", string(b))
	fmt.Printf("Memory:\n%v\n", string(b))

	// Network
	b, err = json.MarshalIndent(sInfo.Stats[0].Network, "", "    ")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "json.Marshal(receivedPodList): %v", err)
		return
	}

	fmt.Fprintf(w, "Network:\n%v\n", string(b))
	fmt.Printf("Network:\n%v\n", string(b))
}
