#!/bin/bash

# Script para insertar equipos de ejemplo en DynamoDB
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

print_banner "Insertando equipos de ejemplo en DynamoDB..."

# Función auxiliar para insertar un equipo
insert_team() {
    local sport=$1
    local name=$2
    local category=$3
    local entity_id="Entity#Team"
    local id="SPORT#${sport}#NAME#${name}"
    
    awslocal dynamodb put-item \
        --table-name "$TABLE_NAME" \
        --region "$REGION" \
        --item "{
            \"EntityId\": {\"S\": \"${entity_id}\"},
            \"Id\": {\"S\": \"${id}\"},
            \"Category\": {\"N\": \"${category}\"},
            \"Sport\": {\"S\": \"${sport}\"}
        }" \
        --return-consumed-capacity TOTAL > /dev/null
    
    echo "✓ Equipo insertado: ${name} (${sport}, Categoría ${category})"
}

# Equipos de Fútbol
insert_team "Football" "Los Leones FC" 5
insert_team "Football" "Real Madrid Local" 6
insert_team "Football" "FC Barcelona Fans" 4
insert_team "Football" "Atlético de Madrid" 7
insert_team "Football" "River Plate B" 3
insert_team "Football" "Boca Juniors Local" 4

# Equipos de Pádel
insert_team "Paddle" "Rocket Pádel" 5
insert_team "Paddle" "Smash Team" 6
insert_team "Paddle" "Paddle Masters" 7
insert_team "Paddle" "Club Pádel Buenos Aires" 4
insert_team "Paddle" "Pádel Pro" 6

# Equipos de Tenis
insert_team "Tennis" "Tennis Club Elite" 5
insert_team "Tennis" "Ace Masters" 6
insert_team "Tennis" "Tennis Buenos Aires" 4
insert_team "Tennis" "Pro Tennis Team" 7

print_banner "Equipos insertados correctamente"

