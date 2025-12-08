#!/bin/bash

# Script para insertar anuncios de partidos de ejemplo en DynamoDB
# Usa fechas relativas al día de ejecución
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

print_banner "Insertando anuncios de partidos de ejemplo en DynamoDB..."

# Obtener fecha/hora actual como base
# Intentar usar timezone de Argentina, fallback a UTC
if TZ="America/Argentina/Buenos_Aires" date +%s >/dev/null 2>&1; then
    NOW=$(TZ="America/Argentina/Buenos_Aires" date +%s)
    TODAY_DATE=$(TZ="America/Argentina/Buenos_Aires" date +%Y-%m-%d)
    TODAY_START=$(TZ="America/Argentina/Buenos_Aires" date -d "${TODAY_DATE} 00:00:00" +%s 2>/dev/null || date -d "${TODAY_DATE} 00:00:00" +%s)
else
    NOW=$(date +%s)
    TODAY_DATE=$(date +%Y-%m-%d)
    TODAY_START=$(date -d "${TODAY_DATE} 00:00:00" +%s 2>/dev/null || date -u -d "${TODAY_DATE} 00:00:00" +%s)
fi

# Función auxiliar para generar UUID (simple pero efectivo)
generate_uuid() {
    cat /proc/sys/kernel/random/uuid 2>/dev/null || \
    python3 -c "import uuid; print(uuid.uuid4())" 2>/dev/null || \
    echo "$(date +%s)-$(shuf -i 1000-9999 -n 1)-$(shuf -i 1000-9999 -n 1)-$(shuf -i 1000-9999 -n 1)-$(shuf -i 100000000000-999999999999 -n 1)"
}

# Función auxiliar para insertar un anuncio de partido
insert_match_announcement() {
    local team_name=$1
    local sport=$2
    local day_offset=$3          # Días desde hoy (0 = hoy, 1 = mañana, etc.)
    local hour=$4                # Hora de inicio (24h format)
    local duration_hours=$5      # Duración en horas
    local country=$6
    local province=$7
    local locality=$8
    local range_type=$9          # SPECIFIC, GREATER_THAN, LESS_THAN, BETWEEN
    local categories=${10}       # Array JSON de categorías o "null"
    local min_level=${11}        # Min level o "null"
    local max_level=${12}        # Max level o "null"
    local status=${13}
    
    local id=$(generate_uuid)
    local entity_id="Entity#MatchAnnouncement"
    
    # Calcular timestamps
    local day_seconds=$((TODAY_START + (day_offset * 86400)))
    local start_time=$((day_seconds + (hour * 3600)))
    local end_time=$((start_time + (duration_hours * 3600)))
    local created_at=$NOW
    local expires_at=$((created_at + (30 * 86400)))  # 30 días desde ahora
    
    # Construir el item JSON según el tipo de rango usando un archivo temporal
    local temp_file=$(mktemp)
    
    {
        echo "{"
        echo "    \"EntityId\": {\"S\": \"${entity_id}\"},"
        echo "    \"Id\": {\"S\": \"${id}\"},"
        echo "    \"TeamName\": {\"S\": \"${team_name}\"},"
        echo "    \"Sport\": {\"S\": \"${sport}\"},"
        echo "    \"Day\": {\"N\": \"${day_seconds}\"},"
        echo "    \"StartTime\": {\"N\": \"${start_time}\"},"
        echo "    \"EndTime\": {\"N\": \"${end_time}\"},"
        echo "    \"Country\": {\"S\": \"${country}\"},"
        echo "    \"Province\": {\"S\": \"${province}\"},"
        echo "    \"Locality\": {\"S\": \"${locality}\"},"
        echo "    \"RangeType\": {\"S\": \"${range_type}\"},"
        
        # Agregar Categories si existe
        if [ "$categories" != "null" ]; then
            echo -n "    \"Categories\": {\"L\": ["
            local first=true
            for cat in $categories; do
                if [ "$first" = true ]; then
                    first=false
                else
                    echo -n ", "
                fi
                echo -n "{\"N\": \"$cat\"}"
            done
            echo "]},"
        fi
        
        # Agregar MinLevel si existe
        if [ "$min_level" != "null" ]; then
            echo "    \"MinLevel\": {\"N\": \"${min_level}\"},"
        fi
        
        # Agregar MaxLevel si existe
        if [ "$max_level" != "null" ]; then
            echo "    \"MaxLevel\": {\"N\": \"${max_level}\"},"
        fi
        
        echo "    \"Status\": {\"S\": \"${status}\"},"
        echo "    \"CreatedAt\": {\"N\": \"${created_at}\"},"
        echo "    \"ExpiresAt\": {\"N\": \"${expires_at}\"}"
        echo "}"
    } > "$temp_file"
    
    awslocal dynamodb put-item \
        --table-name "$TABLE_NAME" \
        --region "$REGION" \
        --item file://"$temp_file" \
        --return-consumed-capacity TOTAL > /dev/null
    
    rm -f "$temp_file"
    
    # Formatear fecha para mostrar (usar el timezone disponible)
    if TZ="America/Argentina/Buenos_Aires" date -d "@${day_seconds}" "+%Y-%m-%d" >/dev/null 2>&1; then
        local date_str=$(TZ="America/Argentina/Buenos_Aires" date -d "@${day_seconds}" "+%Y-%m-%d")
    else
        local date_str=$(date -d "@${day_seconds}" "+%Y-%m-%d" 2>/dev/null || date -u -d "@${day_seconds}" "+%Y-%m-%d")
    fi
    echo "✓ Anuncio insertado: ${team_name} - ${sport} - ${date_str} ${hour}:00 - ${status}"
}

# Anuncios para hoy (día 0)
insert_match_announcement "Los Leones FC" "Football" 0 18 2 "Argentina" "Buenos Aires" "Capital Federal" "BETWEEN" "null" "4" "6" "PENDING"
insert_match_announcement "Rocket Pádel" "Paddle" 0 20 2 "Argentina" "Buenos Aires" "Palermo" "GREATER_THAN" "null" "5" "null" "PENDING"
insert_match_announcement "Tennis Club Elite" "Tennis" 0 16 2 "Argentina" "Buenos Aires" "Recoleta" "SPECIFIC" "5 6" "null" "null" "PENDING"

# Anuncios para mañana (día 1)
insert_match_announcement "Real Madrid Local" "Football" 1 19 2 "Argentina" "Buenos Aires" "Belgrano" "GREATER_THAN" "null" "5" "null" "PENDING"
insert_match_announcement "Smash Team" "Paddle" 1 21 2 "Argentina" "Buenos Aires" "Palermo" "BETWEEN" "null" "4" "7" "PENDING"
insert_match_announcement "Ace Masters" "Tennis" 1 17 2 "Argentina" "Buenos Aires" "Núñez" "LESS_THAN" "null" "null" "6" "PENDING"
insert_match_announcement "FC Barcelona Fans" "Football" 1 18 2 "Argentina" "Buenos Aires" "Caballito" "SPECIFIC" "4 5" "null" "null" "PENDING"

# Anuncios para pasado mañana (día 2)
insert_match_announcement "Atlético de Madrid" "Football" 2 20 2 "Argentina" "Buenos Aires" "San Telmo" "GREATER_THAN" "null" "6" "null" "PENDING"
insert_match_announcement "Paddle Masters" "Paddle" 2 19 2 "Argentina" "Buenos Aires" "Villa Crespo" "SPECIFIC" "6 7" "null" "null" "PENDING"
insert_match_announcement "Pro Tennis Team" "Tennis" 2 18 2 "Argentina" "Buenos Aires" "Barracas" "GREATER_THAN" "null" "6" "null" "PENDING"

# Anuncios para 3 días
insert_match_announcement "River Plate B" "Football" 3 17 2 "Argentina" "Buenos Aires" "Villa Lugano" "BETWEEN" "null" "3" "5" "PENDING"
insert_match_announcement "Club Pádel Buenos Aires" "Paddle" 3 20 2 "Argentina" "Buenos Aires" "Palermo" "LESS_THAN" "null" "null" "5" "PENDING"

# Anuncios para 5 días
insert_match_announcement "Boca Juniors Local" "Football" 5 19 2 "Argentina" "Buenos Aires" "La Boca" "SPECIFIC" "4 5 6" "null" "null" "PENDING"
insert_match_announcement "Pádel Pro" "Paddle" 5 21 2 "Argentina" "Buenos Aires" "Palermo" "GREATER_THAN" "null" "6" "null" "PENDING"
insert_match_announcement "Tennis Buenos Aires" "Tennis" 5 16 2 "Argentina" "Buenos Aires" "Recoleta" "BETWEEN" "null" "4" "6" "PENDING"

# Un anuncio confirmado
insert_match_announcement "Los Leones FC" "Football" 7 18 2 "Argentina" "Buenos Aires" "Capital Federal" "GREATER_THAN" "null" "5" "null" "CONFIRMED"

print_banner "Anuncios de partidos insertados correctamente"

