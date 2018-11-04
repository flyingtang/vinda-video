#!bin/sh

docker login --username=txg5214 registry-vpc.cn-shenzhen.aliyuncs.com

docker build -t vinda-video .

docker tag vinda-video:latest registry-vpc.cn-shenzhen.aliyuncs.com/vinda/vinda-video:1.0.0

docker push registry-vpc.cn-shenzhen.aliyuncs.com/vinda/vinda-video:1.0.0
