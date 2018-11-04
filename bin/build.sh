#!bin/sh

docker login --username=txg5214 registry.cn-shenzhen.aliyuncs.com

docker build -t vinda-video .

docker tag vinda-video:latest registry.cn-shenzhen.aliyuncs.com/vinda/vinda-video:1.0.0

docker push registry.cn-shenzhen.aliyuncs.com/vinda/vinda-video:1.0.0
