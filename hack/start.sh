#!/bin/bash

server="./build/velero-volume-controller/velero-volume-controller"
let item=0
item=`ps -ef | grep $server | grep -v grep | wc -l`

if [ $item -eq 1 ]; then
	echo "The velero-volume-controller is running, shut it down..."
	pid=`ps -ef | grep $server | grep -v grep | awk '{print $2}'`
	kill -9 $pid
fi

echo "Start velero-volume-controller now ..."
make src.build
./build/velero-volume-controller/velero-volume-controller -c ./examples/config.yaml -logtostderr=true -v=5 >> ./velero-volume-controller.log 2>&1 &
