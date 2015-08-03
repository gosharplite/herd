package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosharplite/herd/log"
	"net/http"
)

func getEventHandler(w http.ResponseWriter, r *http.Request) {
	defer log.Un(log.Trace("getEventHandler"))

	events, err := k8s.ListEvents()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, log.Err("k8s.ListEvents(): %v", err))
		return
	}

	j, err := json.MarshalIndent(events.Items, "", "    ")
	if err != nil {
		log.Err("json.MarshalIndent: %v", err)
		return
	}

	fmt.Fprintf(w, string(j))
}
