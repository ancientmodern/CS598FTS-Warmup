#!/bin/bash

# COPY AND PASTE ONLY. DO NOT DIRECTLY RUN THIS

wget https://github.com/etcd-io/etcd/releases/download/v3.5.7/etcd-v3.5.7-linux-amd64.tar.gz
tar xzvf etcd-v3.5.7-linux-amd64.tar.gz
sudo mv etcd-v3.5.7-linux-amd64/etcd* /bin/
rm etcd-v3.5.7-linux-amd64* -rf

TOKEN=token-01
CLUSTER_STATE=new
NAME_1=machine-1
NAME_2=machine-2
NAME_3=machine-3
NAME_4=machine-4
NAME_5=machine-5
HOST_1=128.110.217.160
HOST_2=128.110.217.137
HOST_3=128.110.217.131
HOST_4=128.110.217.155
HOST_5=128.110.217.120
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2380,${NAME_3}=http://${HOST_3}:2380,${NAME_4}=http://${HOST_4}:2380,${NAME_5}=http://${HOST_5}:2380
THIS_NAME=${NAME_1}
THIS_IP=${HOST_1}
nohup etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN} &