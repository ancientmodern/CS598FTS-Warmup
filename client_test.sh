#!/bin/bash

for i in {1..10}
do
    ./output/client --numRead=10 --numWrite=10 --numInitial=10 --cid=$i &
done

wait
echo "All done"
pkill -f "./output/client"