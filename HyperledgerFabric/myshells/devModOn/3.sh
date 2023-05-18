# execute this as step 3.

# Create channel and join peer.
cd ~/HyperledgerFabric/fabric
export PATH=/usr/local/go/bin:$PATH
export PATH=$(pwd)/build/bin:$PATH
export FABRIC_CFG_PATH=$(pwd)/sampleconfig
configtxgen -channelID ch1 -outputCreateChannelTx ch1.tx -profile SampleSingleMSPChannel -configPath $FABRIC_CFG_PATH
peer channel create -o 127.0.0.1:7050 -c ch1 -f ch1.tx
peer channel join -b ch1.block

# build and start the chaincode

# cd ~/HyperledgerFabric/mycodes/atcc
cd $1
go build -o simpleChaincode .
mv simpleChaincode ~/HyperledgerFabric/fabric
cd ~/HyperledgerFabric/fabric


CORE_CHAINCODE_LOGLEVEL=debug CORE_PEER_TLS_ENABLED=false CORE_CHAINCODE_ID_NAME=mycc:1.0 ./simpleChaincode -peer.address 127.0.0.1:7052
