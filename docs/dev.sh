#!/bin/sh

myFile="./app"
 if [ ! -f "$myFile" ]; then
	rm $myFile
 fi

 (go build -o app; docker build -t cmdb-collector:1.0.0 .)
 (docker tag cmdb-collector:1.0.0 registry-jinan-lab.inspurcloud.cn/service/lma/cmdb-collector:1.0.0)
 (docker push registry-jinan-lab.inspurcloud.cn/service/lma/cmdb-collector:1.0.0)

