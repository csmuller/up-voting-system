#!/bin/zsh
# This script sends registration transactions created with 'generate_voter_tx.sh'. The arguments
# for this command have to be consistent with the arguments used when producing registration
# transactions with 'generate_voter_tx.sh'.

if (( $# < 4 )); then
    echo "./sendvoter.sh [begin range] [end range] [home number] [sequence number]."
    exit 1
fi

let h=$3
let accNr=$3+5
let s=$4

if [ ! -e ~/.vcli$h/voter-txs ]; then
  echo "Directory with voter transactions doesn't exist."
  exit 0
fi

vcli query pbb params --home ~/.vcli$h

let startTime=$(date +%s)

echo "--> Sending voter transactions starting at sequence number $s." 
for (( v = $1; v <= $2; v++ )); do

    vcli tx broadcast ~/.vcli$h/voter-txs/$v-signed.json --account-number $accNr --sequence $s \
      --home ~/.vcli$h --gas 100000000000

    let s++
    let currTime=$(date +%s)
    let elapsed=$((currTime - startTime))
    let nrOfTxs=$(($v - $1 + 1))
    echo "Broadcasted $nrOfTxs voter transactions in $elapsed seconds."
done

let endtime=$(date +%s)
let elapsed=$((endtime - startTime))
let nrOfTxs=$(($2 - $1 + 1))

echo "DONE: Broadcasted $nrOfTxs voter transactions in $elapsed seconds."

