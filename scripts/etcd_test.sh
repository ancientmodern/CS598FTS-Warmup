#!/bin/bash

# COPY AND PASTE ONLY. DO NOT DIRECTLY RUN THIS

HOST_1=128.110.217.160
HOST_2=128.110.217.137
HOST_3=128.110.217.131
HOST_4=128.110.217.155
HOST_5=128.110.217.120
ENDPOINTS=$HOST_1:2379,$HOST_2:2379,$HOST_3:2379,$HOST_4:2379,$HOST_5:2379

# detect leader
etcdctl --endpoints=${ENDPOINTS} -w table endpoint status

# change leader
etcdctl --endpoints=${ENDPOINTS} move-leader 155bef1f5d10a316

# write
benchmark --endpoints=${ENDPOINTS} --clients=32 put --key-size=24 --total=50000 --val-size=10

# read
benchmark --endpoints=${ENDPOINTS} --clients=32 range YOUR_KEY --consistency=l --total=50000

# mixed (final choice)
benchmark --endpoints=${ENDPOINTS} --clients=100 txn-mixed --key-size=24 --total=100000 --val-size=10 --rw-ratio=1 --consistency=l

# limit cpu
cpulimit -l 100 -p 11399 -b