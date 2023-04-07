#!/bin/bash
# Path: scripts/benchmark.sh
# This script is used to benchmark Mainflux migration from version 0.13.0 to 0.14.0.
# It provisions users and things on Mainflux and then runs the migration command.
# The results are saved to a JSON file.
# The script can be run with different parameters to benchmark different scenarios.
# The script requires hyperfine to be installed.
# The script requires docker and docker-compose to be installed.
# The script requires mainflux to be cloned in the same parent directory as migrations.
# The script requires mainflux-migrate to be built.

# Wait time for docker compose to start
WAIT_TIME=5

# Path to Docker Compose file
DOCKER_COMPOSE_FILE="../mainflux/docker/docker-compose.yml"

# Prefix for user accounts
USER_PREFIX="testa"

# Prefix for things and channels
TC_PREFIX="seed"

# Migration command
MIGRATION_COMMAND="$1"

# Operation. It can be import or export
OPERATION=""

# Hyperfine prepare command on import
HYPERFINE_PREPARE_IMPORT="echo 'delete from clients; delete from groups; delete from policies;' | docker exec -i mainflux-things-db psql -U mainflux -d things"

# Stop and remove existing Docker Compose project and volumes
function stop_and_remove_docker_compose() {
    docker-compose -f "$DOCKER_COMPOSE_FILE" down --remove-orphans > /dev/null 2>&1
    docker volume rm $(docker volume ls -qf "name=docker_mainflux-*") > /dev/null 2>&1
}

# Start docker compose and waits for 2 seconds
function start_docker_compose() {
    docker-compose -f "$DOCKER_COMPOSE_FILE" up -d > /dev/null 2>&1
    sleep "$WAIT_TIME"
}

# Provision users and things on mainflux
function provision() {
    local maxusers=$1
    local maxthings=$2
    for i in $(seq 1 $maxusers); do
        ./../mainflux/tools/provision/provision -u "$USER_PREFIX"$i@example.com -p 12345678 --num $maxthings --prefix "$TC_PREFIX"  > /dev/null 2>&1
    done
}

# Run hyperfine and save results to file
function benchmark_migrate() {
    local output_file=$1
    hyperfine --runs 10 --export-json "$output_file" "$MIGRATION_COMMAND"> /dev/null 2>&1
}

function benchmark_migrate() {
    local output_file=$1
    local filename=$(basename $output_file)
    local prefix_filename="${output_file%/*}/$OPERATION$filename"

    if [[ "$OPERATION" == "export" ]]; then
        hyperfine --runs 10 --export-json "$prefix_filename" "$MIGRATION_COMMAND" > /dev/null 2>&1  
    elif [[ "$OPERATION" == "import" ]]; then
        hyperfine --runs 10 --export-json "$prefix_filename" --prepare "$HYPERFINE_PREPARE_IMPORT" "$MIGRATION_COMMAND" > /dev/null 2>&1
    fi
}

# Run the script
function run_script() {
    local -r max_users="$1"
    local -r max_things="$2"
    local -r output_file="$3"
    if [[ "$OPERATION" == "export" ]]; then
        printf "Starting Docker Compose...\n"
        start_docker_compose
        printf "Provisioning %d users, %d things and %d channels on Mainflux...\n" "$1" "$2" "$2"
        provision "$max_users" "$max_things"
        benchmark_migrate "$output_file"
        printf "Stopping Docker Compose...\n"
        stop_and_remove_docker_compose
    elif [[ "$OPERATION" == "import" ]]; then
        printf "Starting Docker Compose...\n"
        export MF_RELEASE_TAG=0.13.0
        start_docker_compose
        provision "$max_users" "$max_things"
        ./build/mainflux-migrate -f 0.13.0 -o export > /dev/null 2>&1
        stop_and_remove_docker_compose
        export MF_RELEASE_TAG=latest
        start_docker_compose
        printf "Importing %d users, %d things and %d channels on Mainflux...\n" "$1" "$2" "$2"
        benchmark_migrate "$output_file"
        printf "Stopping Docker Compose...\n"
        stop_and_remove_docker_compose
    fi
}

# Check if the migration command is provided
if [[ -z "$MIGRATION_COMMAND" ]]; then
    echo "Please provide the migration command as an input parameter."
    echo "Example: ./benchmark.sh \"./build/mainflux-migrate -f 0.13.0 -o export\""
    exit 1
elif echo "$MIGRATION_COMMAND" | grep -q "export"; then
    OPERATION="export"
elif echo "$MIGRATION_COMMAND" | grep -q "import"; then
    OPERATION="import"
else
  echo "The migration command does not contain the word 'export' or 'import'."  
fi

# Make sure docker compose is stopped and volumes are removed
stop_and_remove_docker_compose

# Run the script with different parameters
run_script 1 10 "scripts/1user_10things_each.json"
run_script 1 100 "scripts/1user_100things_each.json"
run_script 1 1000 "scripts/1user_1000things_each.json"
run_script 10 10 "scripts/10users_10things_each.json"
run_script 10 100 "scripts/10users_100things_each.json"
run_script 10 1000 "scripts/10users_1000things_each.json"
run_script 100 10 "scripts/100users_10things_each.json"
run_script 100 100 "scripts/100users_100things_each.json"
