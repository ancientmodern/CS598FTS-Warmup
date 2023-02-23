#!/bin/bash

for i in {1..10}
do
    ./output/client --numRead=10 --numWrite=10 --clientNum="$i" &
done

wait
echo "All done"