export PATH=/usr/local/go/bin/:$PATH
cd ~/HyperledgerFabric/fabric-samples/test-network
./network.sh down
./network.sh up createChannel
./network.sh deployCC -ccn basic -ccp $1 -ccl go
printf "\n\nfollow README.md to continue...\n\n"