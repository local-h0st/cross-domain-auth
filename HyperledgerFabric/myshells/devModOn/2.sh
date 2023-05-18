# execute this as step 2.

# start the peer in dev Mod.
cd ~/HyperledgerFabric/fabric
export CORE_OPERATIONS_LISTENADDRESS=127.0.0.1:9444
export PATH=$(pwd)/build/bin:$PATH
export FABRIC_CFG_PATH=$(pwd)/sampleconfig
FABRIC_LOGGING_SPEC=chaincode=debug CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052 peer node start --peer-chaincodedev=true

# now switch to another terminal window -->