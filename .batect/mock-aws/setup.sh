#!/bin/bash

set -euo pipefail

USERNAME=test-user

aws iam create-user \
    --user-name "${USERNAME}" \
    --endpoint http://127.0.0.1:5000

aws iam attach-user-policy \
    --endpoint http://127.0.0.1:5000 \
    --user-name "${USERNAME}" \
    --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess

CREDENTIALS=$(aws iam create-access-key \
    --endpoint http://127.0.0.1:5000 \
    --user-name "${USERNAME}")

aws configure set aws_access_key_id "$(echo "${CREDENTIALS}" | jq -r '.AccessKey.AccessKeyId')"
aws configure set aws_secret_access_key "$(echo "${CREDENTIALS}" | jq -r '.AccessKey.SecretAccessKey')"
