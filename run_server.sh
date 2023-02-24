#!/bin/bash

port_num=50055=1
size=10

./output/replica --port=$port_num  --size=$size &

sleep 100
while [[ $(ps aux | grep ./output/client | grep -v grep) ]]
do
    sleep 1
    echo "running"
done
pkill -f "./output/replica"