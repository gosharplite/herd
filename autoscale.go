package main

import (
	"encoding/json"
	"github.com/gosharplite/herd/log"
	"time"
)

func autoLoop() {

	// TODO Consider etcd watch

	c := time.Tick(10 * time.Second)
	for _ = range c {

		go func() {
			err := fetchScales()
			if err != nil {
				log.Err("checkScale: %v", err)
			}
		}()
	}
}

func fetchScales() error {
	defer log.Un(log.Trace("fetchScales"))

	v, err := etcd.GetScales()
	if err != nil {
		return err
	}

	for _, n := range v {

		var js jScale
		err := json.Unmarshal([]byte(n.Value), &js)
		if err != nil {
			log.Err("json.Unmarshal(): %v", err)
			continue
		}

		err = checkScale(js)
		if err != nil {
			log.Err("checkScale(): %v", err)
			continue
		}
	}

	return nil
}

func checkScale(js jScale) error {

	//	log.Info("js: %v", js)

	if js.EnableAutoScale == 0 {
		// delete
		err := etcd.Delete(js.ClusterId)
		if err != nil {
			return err
		}
	} else {
		// get pods and cpu
		rc, err := getRc(js.ClusterId)
		if err != nil {
			return err
		}

		// scale
		decide := 0
		for _, c := range rc.Containers {
			if int64(c.Cpu) >= js.CpuMax {
				decide = 1
				break
			}

			if int64(c.Cpu) < js.CpuMin {
				decide = -1
				break
			}
		}

		l := int64(len(rc.Containers) + decide)
		if l <= js.PodMax && l >= js.PodMin && l != int64(len(rc.Containers)) {
			// change
			krc, err := k8s.GetRC(js.ClusterId)
			if err != nil {
				log.Err("k8s.GetRC(): %v", err)
				return err
			}

			krc.Spec.Replicas = int(l)

			krc, err = k8s.Update(krc)
			if err != nil {
				log.Err("k8s.Update(): %v", err)
				return err
			}

			log.Info("Scale %v from %v to %v", js.ClusterId, len(rc.Containers), krc.Spec.Replicas)
		}
	}

	return nil
}
