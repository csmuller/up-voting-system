rm -rf ~/.vcli
rm -rf ~/.acli
rm -rf ~/.pbbd

pbbd init val --chain-id pbb

sed -i -e 's/^max_body_bytes.*/max_body_bytes = 10000000/' ~/.pbbd/config/config.toml
sed -i -e 's/^max_packet_msg_payload_size.*/max_packet_msg_payload_size = 10000000/' \
  ~/.pbbd/config/config.toml 
sed -i -e 's/^max_tx_bytes.*/max_tx_bytes = 10000000/' ~/.pbbd/config/config.toml

echo "12345678" | acli keys add admin
echo "12345678" | vcli keys add voter

acli config chain-id pbb
acli config output json
acli config indent true
acli config trust-node true

vcli config chain-id pbb
vcli config output json
vcli config indent true
vcli config trust-node true

pbbd add-genesis-account $(acli keys show admin -a) 100000000stake,1000foo
# pbbd add-genesis-account $(vcli keys show voter -a) 1foo

echo "12345678" | pbbd gentx --details val --name admin
pbbd collect-gentxs

# pbbd start

# vcli keys show voter -a

# echo "12345678" | acli tx bank send $(acli keys show admin -a) \ 
# $(vcli keys show voter -a) 1foo -y

