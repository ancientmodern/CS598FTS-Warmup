#!/bin/bash

scp -r output/ node-1:~ &
scp -r output/ node-2:~ &
scp -r output/ node-3:~ &

wait
echo "finish"
