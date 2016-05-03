#!/bin/sh

CSM_IMAGE_NAME="csm-mysql:latest"
MYSQL_ROOT_PASSWORD="root123"
CSM_LOG_LEVEL="debug"
MYSQL_SERVICE_PORT_MYSQL="3306"
CSM_DEV_MODE="true"


DOCKER_HOST_IP=`echo ${DOCKER_HOST} | cut -d "/" -f 3 | cut -d ":" -f 1`


docker run --name csm-mysql \
       -p 8081:8081 \
       -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
       -e CSM_LOG_LEVEL=${CSM_LOG_LEVEL} \
       -e MYSQL_SERVICE_HOST=${DOCKER_HOST_IP} \
       -e MYSQL_SERVICE_PORT_MYSQL=${MYSQL_SERVICE_PORT_MYSQL} \
       -e CSM_API_KEY=${CSM_API_KEY} \
       -e CSM_DEV_MODE=${CSM_DEV_MODE} \
       -d ${CSM_IMAGE_NAME}