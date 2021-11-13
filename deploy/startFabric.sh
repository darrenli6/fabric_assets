set -e

export MSYS_NO_PATHCONV=1

starttime=$(date +%s)

./start.sh

docker-compose -f ./docker-compose.yml up -d cli

printf "Success!\n"