#!/bin/bash

# Script para insertar solicitudes de partido de ejemplo en DynamoDB
# Usa awslocal para interactuar con LocalStack

set -euo pipefail

TABLE_NAME="${TABLE_NAME:-SportLinkCore}"
REGION="${AWS_DEFAULT_REGION:-us-east-1}"

print_banner() {
    local message=$1
    local line="========================================================================================================================="
    echo "$line"
    echo "$message"
    echo "$line"
}

print_banner "Insertando solicitudes de partido de ejemplo en DynamoDB..."

NOW=$(date +%s)

# Función para insertar una solicitud de partido
insert_match_request() {
    local id=$1
    local match_offer_id=$2
    local owner_account_id=$3       # Dueño de la oferta (recibe la solicitud)
    local requester_account_id=$4   # Quien envía la solicitud
    local status=$5

    awslocal dynamodb put-item \
        --table-name "$TABLE_NAME" \
        --region "$REGION" \
        --item "{
            \"EntityId\": {\"S\": \"Entity#MatchRequest\"},
            \"Id\": {\"S\": \"${id}\"},
            \"MatchOfferId\": {\"S\": \"${match_offer_id}\"},
            \"OwnerAccountId\": {\"S\": \"${owner_account_id}\"},
            \"RequesterAccountId\": {\"S\": \"${requester_account_id}\"},
            \"Status\": {\"S\": \"${status}\"},
            \"CreatedAt\": {\"N\": \"${NOW}\"}
        }" \
        --return-consumed-capacity TOTAL > /dev/null

    echo "✓ Solicitud insertada: ${id} [offer: ${match_offer_id}] de ${requester_account_id} → ${owner_account_id} (${status})"
}

# Solicitud de rivalfc para la oferta fija de cabrerajjorge
insert_match_request "req-rival-001" "offer-cabrerajjorge-001" "cabrerajjorge" "rivalfc" "PENDING"

print_banner "Solicitudes de partido insertadas correctamente"
