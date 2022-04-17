#!/bin/bash

set -euo pipefail

YEAR=$(date --date='-1 month' +%Y)
MONTH=$(date --date='-1 month' +%m)

put_csv() {
  local report=$1
  local cluster=$2
  local file_name="${report}-${cluster}.csv"

  echo $'REPORT,CLUSTER\n'${report},${cluster} > /tmp/${file_name}

  local key="${YEAR}/${MONTH}/${report}/${cluster}.csv"
  aws s3api put-object \
    --endpoint ${AWS_ENDPOINT} \
    --bucket ${BUCKET} \
    --key "${key}" \
    --body /tmp/${file_name} > /dev/null

  echo "Put test csv to ${BUCKET}/${key}"
  aws s3 cp \
    --endpoint ${AWS_ENDPOINT} \
    s3://${BUCKET}/${key} \
    -
}

aws s3api create-bucket \
    --endpoint ${AWS_ENDPOINT} \
    --bucket ${BUCKET}

put_csv single cluster-1
put_csv single cluster-2
put_csv cumulative cluster-1
put_csv cumulative cluster-2
