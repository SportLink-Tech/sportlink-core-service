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

<<<<<<< HEAD
print_banner "Waiting for the CloudFormation stack to be created..."

awslocal cloudformation wait stack-create-complete --stack-name core-dynamo-table-stack

print_banner "======================================= Localstack Setup Ends ==========================================================="
=======
print_banner "======================================= Localstack Setup Ends ==========================================================="

>>>>>>> ede636af5f05fcc09b639d934c1122b83ee8747b
