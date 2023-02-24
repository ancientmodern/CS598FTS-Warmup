#!/bin/bash

port_num=50051
size=10

./output/replica --port=$port_num  --size=$size &

sleep 5
while  ps aux | grep ./output/client | grep -v grep;
do
    sleep 1
    echo "running"
done
pkill -f "./output/replica"