#!/bin/bash

read=10
write=10
initial=10
num=1
correct=false

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
    -n|--num)
      num="$2"
      shift # past argument
      shift # past value
      ;;
    -c|--correct)
      correct="$2"
      shift # past argument
      shift # past value
      ;;
  esac
done

for (( i = 1; i <= $num; i++ ))
do
    ./output/client --numRead=$read --numWrite=$write --numInitial=$initial --cid=$i --isForCorrect=$correct &
done

wait
echo "All done"
pkill -f "./output/client"
