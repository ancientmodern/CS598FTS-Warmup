#!/bin/bash

read=10
write=10
initial=10

while [[ $# -gt 0 ]]; do
  case $1 in
    -r|--read)
      read="$2"
      shift # past argument
      shift # past value
      ;;
    -w|--write)
      write="$2"
      shift # past argument
      shift # past value
      ;;
    -i|--initial)
      initial="$2"
      shift # past argument
      shift # past value
      ;;
  esac
done

for i in {1..10}
do
    ./output/client --numRead=$read --numWrite=$write --numInitial=$initial --cid=$i &
done

wait
echo "All done"
pkill -f "./output/client"
