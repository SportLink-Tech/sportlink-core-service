#!/bin/bash

#set -euo pipefail

print_banner() {
    local message=$1
    local line="========================================================================================================================="

    echo "$line"
    echo "$message"
    echo "$line"
}

print_banner "======================================= Creating Dynamo Table ==========================================================="

awslocal cloudformation create-stack \
  --stack-name "core-dynamo-table-stack" \
  --template-body file://core-dynamo-table.yml \
  --region "${AWS_REGION}" \
  --output table

print_banner "======================================= Localstack Setup Ends ==========================================================="

