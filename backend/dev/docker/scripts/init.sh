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
  --region "${AWS_DEFAULT_REGION}" \
  --output table


print_banner "========================================= Creating SQS Queue ============================================================"

awslocal cloudformation create-stack \
  --stack-name "sqs-queue-stack" \
  --template-body file://sqs-queue.yml \
  --region "${AWS_DEFAULT_REGION}" \
  --output table

print_banner "Waiting for DyanamoDb CloudFormation stack to be created..."

awslocal cloudformation wait stack-create-complete --stack-name core-dynamo-table-stack

print_banner "Waiting for SQS CloudFormation stack to be created..."

awslocal cloudformation wait stack-create-complete --stack-name sqs-queue-stack

print_banner "======================================= Seeding Database ==========================================================="

# Ejecutar scripts de seeding
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -f "${SCRIPT_DIR}/seed-teams.sh" ]; then
    bash "${SCRIPT_DIR}/seed-teams.sh"
else
    echo "Warning: seed-teams.sh not found"
fi

if [ -f "${SCRIPT_DIR}/seed-match-announcements.sh" ]; then
    bash "${SCRIPT_DIR}/seed-match-announcements.sh"
else
    echo "Warning: seed-match-announcements.sh not found"
fi

print_banner "======================================= Localstack Setup Ends ==========================================================="

echo "All services created and database seeded."