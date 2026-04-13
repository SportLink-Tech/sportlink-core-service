#!/bin/bash

#set -euo pipefail

print_banner() {
    local message=$1
    local line="========================================================================================================================="

    echo "$line"
    echo "$message"
    echo "$line"
}

print_banner "======================================= Seeding Database ==========================================================="

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

if [ -f "${SCRIPT_DIR}/seed-match-requests.sh" ]; then
    bash "${SCRIPT_DIR}/seed-match-requests.sh"
else
    echo "Warning: seed-match-requests.sh not found"
fi

print_banner "======================================= Seeding Ends ==========================================================="

echo "Database seeded."
