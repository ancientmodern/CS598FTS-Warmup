#!/bin/bash

make

ssh node-1 "pkill -9 ryu-manager; pkill -9 replica; pkill -9 proxy"
ssh node-2 "pkill -9 ryu-manager; pkill -9 replica; pkill -9 proxy"
ssh node-3 "pkill -9 ryu-manager; pkill -9 replica; pkill -9 proxy"

sleep 2

scp -r output/ node-1:~ &
scp -r output/ node-2:~ &
scp -r output/ node-3:~ &

wait
echo "finish"
