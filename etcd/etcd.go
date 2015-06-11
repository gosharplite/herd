package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/gosharplite/herd/log"
	"net/http"
	"strconv"
	"time"
)

var (
	ETCD_PREFIX = "/gigacloud.com/autoscale/"
	c           *etcd.Client
)

func init() {
	machines := []string{"http://192.168.3.36:2379", "http://192.168.3.37:2379", "http://192.168.3.38:2379"}
	c = etcd.NewClient(machines)
	if c == nil {
		log.Err("etcd.NewClient(machines)")
	}
}

func Set(key, value string) error {

	if key == "" {
		log.Err("key is empty")
		return errors.New("key is empty")
	}

	_, err := c.Set(ETCD_PREFIX+key, value, 0)
	if err != nil {
		log.Err("c.Set(): %v", err)
		return err
	}

	return nil
}

func Get(key string) (string, error) {

	if key == "" {
		log.Err("key is empty")
		return "", errors.New("key is empty")
	}

	r, err := c.Get(ETCD_PREFIX+key, false, false)
	if err != nil {
		log.Err("c.Get(): %v", err)
		return "", err
	}

	return r.Node.Value, nil
}

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
