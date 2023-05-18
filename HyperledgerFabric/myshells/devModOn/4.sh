# execute this as step 4.

cd ~/HyperledgerFabric/fabric
export PATH=/usr/local/go/bin:$PATH
export PATH=$(pwd)/build/bin:$PATH
export FABRIC_CFG_PATH=$(pwd)/sampleconfig

# Approve and commit the chaincode definition
peer lifecycle chaincode approveformyorg  -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --package-id mycc:1.0
peer lifecycle chaincode checkcommitreadiness -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')"
peer lifecycle chaincode commit -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --peerAddresses 127.0.0.1:7051

# Next Steps: do your own tests.
printf "\nTo use peer, remember to export these on your terminal:\n"
printf "export PATH=~/HyperledgerFabric/fabric/build/bin:\$PATH\n"
printf "export FABRIC_CFG_PATH=~/HyperledgerFabric/fabric/sampleconfig\n"
printf "Now just do your own tests!\n\n"