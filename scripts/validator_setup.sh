# These commands can be used to setup two validators in docker containers on the same machine.

docker build -t val01 --build-arg NR=01 -f Dockerfile_val .
docker build -t val02 --build-arg NR=02 -f Dockerfile_val .

## Validator 1 ##
docker run -it --name val01 --network="pbb_pbb" --ip="192.167.10.101" -p 26656-26657:26656-26657 \
  val01 /bin/bash

## Validator 2 ##
docker run -it --name val02 --network="pbb_pbb" --ip="192.167.10.102" -p 26666-26667:26656-26657 \
  val02 /bin/bash

## Retrieve admin addresses of both validators and add to genesis block on first validator. I.e.
eacli keys show admin -a # Retrieve the addresses of the admin accounts.
# Execute on validator 1 for both the val1 and val2 admin accounts.
pbbd add-genesis-account cosmos1xhz3vwgvu37khfr85vq3ukd5mm7fq7pdq7tu0l 100000000stake,1000foo # val02

## Collect gentxs from validator 2 and copy it to validator 1.
docker cp val02:/root/.pbbd/config/gentx ~/tmp/val2
docker cp ~/tmp/val2/gentx/* val01:/root/.pbbd/config/gentx

## Collect gentxs on validator 1. Execute this on validator 1.
pbbd collect-gentxs

## Copy the generated genesis blcock from validator 1 to validator 2.
docker cp val01:/root/.pbbd/config/genesis.json ~/tmp/genesis.json
docker cp ~/tmp/genesis.json val02:/root/.pbbd/config/genesis.json

## Now on both docker containers the pbbd can be started.
pbbd start
