package cadv

import (
	"github.com/google/cadvisor/client"
	info "github.com/google/cadvisor/info/v1"
	"github.com/gosharplite/herd/log"
	"strings"
)

var CADV_PORT = ":4194"

func GetMInfo(HostIP string) (*info.MachineInfo, error) {

	client, err := client.NewClient("http://" + HostIP + CADV_PORT)
	if err != nil {
		log.Err("client.NewClient(): %v", err)
		return nil, err
	}

	return client.MachineInfo()
}

func GetCInfo(HostIP, ContainerID string) (info.ContainerInfo, error) {
	client, err := client.NewClient("http://" + HostIP + CADV_PORT)
	if err != nil {
		log.Err("client.NewClient(): %v", err)
		return info.ContainerInfo{}, err
	}

	request := info.ContainerInfoRequest{NumStats: 10}
	sInfo, err := client.DockerContainer(strings.Replace(ContainerID, "docker://", "", -1), &request)
	if err != nil {
		log.Err("client.DockerContainer(): %v", err)
		return info.ContainerInfo{}, err
	}

	return sInfo, nil
}
