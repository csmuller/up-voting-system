#!/bin/zsh

# This script generates registration transactions but does not send them. The first two
# arguments define the
# range of ballots to produce. E.g. if you want to produce 100 transactions the range is 1 to 100.
# If after that you need to produce 100 additional ones provide the range 101 to 200. The third
# argument specifies the home directory number. This implies that you are using multiple home
# directories for the vcli pplication (e.g. ~/.vcli1 and ~/.vcli2). It is expected that the home
# directory specified with this argument already exists and contains a param.json file that contains
# the election parameters.
# The last argument denotes the next free sequence number of the voter account. If you have
# already sent 100 transaction with the account then the next sequence number is 100.
# The generated ballots are stored in the vcli home directory.

if (( $# < 4 )); then
    echo "./genvoter.sh [begin range] [end range] [home number] [sequence number]."
    exit 1
fi

let h=$3
let accNr=$3+5
let s=$4

if [ ! -e ~/.vcli$h/creds ]; then
    echo "--> Create credentials folder at ~/.vcli$h/creds."
    mkdir -p ~/.vcli$h/creds
fi

if [ ! -e ~/.vcli$h/voter-txs ]; then
    echo "--> Create transactions folder at ~/.vcli$h/voter-txs."
    mkdir -p ~/.vcli$h/voter-txs
fi

vcli query pbb params --home ~/.vcli$h

let startTime=$(date +%s)

echo "--> Creating new voter transactions starting at sequence number $s." 
for (( v = $1; v <= $2; v++ )); do

    vcli tx pbb new-voter ~/.vcli$h/creds/$v.pub ~/.vcli$h/creds/$v.priv ~/.vcli$h/params.json \
      --from $(vcli keys show voter -a --home ~/.vcli$h) --generate-only --sequence $s \
      --account-number $accNr --home ~/.vcli$h --gas 1000000000000 > ~/.vcli$h/voter-txs/$v.json

    echo "12345678" | vcli tx sign ~/.vcli$h/voter-txs/$v.json --from voter --offline \
      --account-number $accNr --sequence $s --home ~/.vcli$h > ~/.vcli$h/voter-txs/$v-signed.json

    let s++
    let currTime=$(date +%s)
    let elapsed=$((currTime - startTime))
    let nrOfTxs=$(($v - $1 + 1))
    echo "Generated $nrOfTxs voter transactions in $elapsed seconds."
done

let endtime=$(date +%s)
let elapsed=$((endtime - startTime))
let nrOfTxs=$(($2 - $1 + 1))

echo "DONE: Generated $nrOfTxs voter transactions in $elapsed seconds."
