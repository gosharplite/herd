# herd
Animals playground.

CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .

sudo docker build -t gosharplite/herd:v8 .

sudo docker push gosharplite/herd:v8

sudo docker run --publish 8090:8090 gosharplite/herd:v8 -etcd_machines="http://192.168.6.2:2379,http://192.168.6.3:2379,http://192.168.6.4:2379" -k8s_host=192.168.6.13:8080

-------------
[Unit]
Description=K8s auto scale daemon

[Service]
ExecStartPre=-/usr/bin/docker kill herd1
ExecStartPre=-/usr/bin/docker rm herd1
ExecStart=/usr/bin/docker run --publish 8090:8090 --name herd1 gosharplite/herd:v8 -etcd_machines="http://192.168.6.2:2379,http://192.168.6.3:2379,http://192.168.6.4:2379" -k8s_host=192.168.6.13:8080
ExecStop=/usr/bin/docker stop herd1
Restart=always
RestartSec=10s

[X-Fleet]
MachineMetadata=role=service
