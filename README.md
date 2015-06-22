# herd
Animals playground.

CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .

sudo docker build -t gosharplite/herd:v5 .

sudo docker run --publish 8090:8090 gosharplite/herd:v5 -etcd_machines="http://192.168.6.2:2379,http://192.168.6.3:2379,http://192.168.6.4:2379" -k8s_host=192.168.6.13:8080