#!/bin/zsh
# This script sends ballot transactions created with 'generate_ballot_tx.sh'. The arguments for
# this command have to be consistent with the arguments used when producing ballots with
# 'generate_ballot_tx.sh'.

if (( $# < 4 )); then
    echo "./sendballot.sh [begin range] [end range] [home number] [sequence number]."
    exit 1
fi

let h=$3
let accNr=$3+5
let s=$4

if [ ! -e ~/.vcli$h/ballot-txs ]; then
  echo "Directory with ballot transactions doesn't exist."
  exit 0
fi

vcli query pbb poly --home ~/.vcli$h

let startTime=$(date +%s)

echo "--> Sending ballot transactions starting at sequence number $s." 
for (( v = $1; v <= $2; v++ )); do

    vcli tx broadcast ~/.vcli$h/ballot-txs/$v-signed.json --account-number $accNr --sequence $s \
      --home ~/.vcli$h

    let s++
    let currTime=$(date +%s)
    let elapsed=$((currTime - startTime))
    let nrOfTxs=$(($v - $1 + 1))
    echo "Broadcasted $nrOfTxs ballot transactions in $elapsed seconds."
    sleep .5s
done

let endtime=$(date +%s)
let elapsed=$((endtime - startTime))
let nrOfTxs=$(($2 - $1 + 1))

echo "Broadcasted $nrOfTxs ballot transactions in $elapsed seconds."
