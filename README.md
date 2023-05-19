# Project cross-domain-auth

这个项目是做跨域认证的，基于Hyperledger Fabric 2.5，采用Intel SGX作为Truetsed Execution Environment的硬件支持

## devMod
为了方便测试链码，Hyperledger官方给出了[devMod](https://hyperledger-fabric.readthedocs.io/en/release-2.5/peer-chaincode-devmode.html)。根据教程一条条在CLI里面敲命令太麻烦了，因此我写了4个自动化脚本，放在~/HyperledgerFabric/myshells/devModOn目录下。同时在～下写了dev.sh，能够方便地调用那四个shell，要开启devMod，请按照以下步骤：

```
// 新建一个shell窗口
./dev.sh 1
// 新建一个shell窗口
./dev.sh 2
// 新建一个shell窗口
./dev.sh 3 "your chaincode path"    // 例如：./dev.sh 3 ~/HyperledgerFabric/mycodes/atcc
// 新建一个shell窗口
./dev.sh 4
// 在第四个shell中按照提示export环境变量，随后即可开始测试链码
```

官方教程部署的链码位于～/HyperledgerFabric/fabric/integration/chaincode/simple/cmd，我整合命令后的sh能够跑官方的chaincode，随后用完自己写的chaincode测试，也就是atcc，能够正常工作。

我测试atcc用的命令如下：
```
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["InitLedger"]}' --isInit
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["QueryAllRecords"]}'
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["RecordExists","#1"]}'
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["ChangeRating","#3","99"]}'
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["QueryRecord","#3"]}'
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["DeleteRecord","#3"]}'
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["QueryAllRecords"]}'
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["AddRecord","#0","payphone","Maroon5","10"]}'
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["QueryAllRecords"]}'
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["UpdateRecord","#0","Moves Like Jagger","Maroon 5","999"]}'
CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["QueryRecord","#0"]}'


```