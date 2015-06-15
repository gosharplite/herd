package main

import (
	"fmt"
	"github.com/gosharplite/herd/log"
	"net/http"
)

func getEventHandler(w http.ResponseWriter, r *http.Request) {

	events, err := k8s.ListEvents()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, log.Err("k8s.ListEvents(): %v", err))
		return
	}

	fmt.Fprintf(w, log.Marshal(events.Items))
}
