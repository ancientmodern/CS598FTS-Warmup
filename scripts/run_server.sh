#!/bin/bash

port_num=50051
size=10

while [[ $# -gt 0 ]]; do
  case $1 in
    -s|--size)
      size="$2"
      shift # past argument
      shift # past value
      ;;
    -p|--port)
      port_num="$2"
      shift # past argument
      shift # past value
      ;;
    -i|--ip)
      ip="$2"
      shift # past argument
      shift # past value
      ;;
  esac
done

./output/replica --port=$port_num --size=$size --ip=$ip &

# sleep 100
# while [[ $(ps aux | grep ./output/client | grep -v grep) ]]
# do
#     sleep 1
#     echo "running"
# done
# pkill -f "./output/replica"
