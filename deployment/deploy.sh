#!/bin/bash

PROJECT_NAME="slowmoon/go-tools"

this_file=`pwd`"/"$0
CURRENT_DIR=`dirname $this_file`

. $CURRENT_DIR/utils.sh

echo `pwd`
echo $(new_tag)

set -ex

docker login -u slowmoon -p along665

docker build  -t ${PROJECT_NAME}:$(new_tag) .

docker push ${PROJECT_NAME}:$(new_tag)

docker rmi ${PROJECT_NAME}:$(new_tag)






