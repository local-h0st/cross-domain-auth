# execute this as step 1

# clean
sudo rm -rf /var/hyperledger/
rm ~/HyperledgerFabric/fabric/simpleChaincode

# generate genesis block
export PATH=/usr/local/go/bin/:$PATH
cd ~/HyperledgerFabric/fabric
make orderer peer configtxgen
export PATH=$(pwd)/build/bin:$PATH
export FABRIC_CFG_PATH=$(pwd)/sampleconfig
sudo mkdir /var/hyperledger
sudo chown ubuntu /var/hyperledger
configtxgen -profile SampleDevModeSolo -channelID syschannel -outputBlock genesisblock -configPath $FABRIC_CFG_PATH -outputBlock "$(pwd)/sampleconfig/genesisblock"

# start the orderer
ORDERER_GENERAL_GENESISPROFILE=SampleDevModeSolo orderer

# now switch to another new terminal window -->
