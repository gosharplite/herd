package log

import (
	"testing"
)

func TestMarshal(t *testing.T) {

	type pod struct {
		ContainerName string `json:"container_name"`
		Cpu           uint64 `json:"cpu"`
		Mem           uint64 `json:"mem"`
	}

	if Marshal(pod{}) != `{
    "container_name": "",
    "cpu": 0,
    "mem": 0
}` {
		t.Fail()
	}
}
