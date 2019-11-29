#!/bin/bash

PROJECT_NAME="go-tools"

this_file=`pwd`"/"$0
CURRENT_DIR=`dirname $this_file`

. $CURRENT_DIR/utils.sh

echo `pwd`
echo $(new_tag)

set -ex
set -eo pipefail

docker login -u slowmoon -p along665

#docker build --squash -t ${PROJECT_NAME}: -f







