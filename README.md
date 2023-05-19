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

官方教程部署的链码位于～/HyperledgerFabric/fabric/integration/chaincode/simple/cmd，我整合命令后的sh能够跑官方的chaincode，随后用我自己写的chaincode测试，也就是atcc，能够正常工作。

开启devMod部署atcc链码后，测试链码功能用的命令如下：
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

## 项目结构
整个项目应该是需要编写链码、VS上的服务端程序、PAS上的服务端程序，device上的用户程序。
由于服务端程序涉及到调用智能合约，因此也属于DApp的范畴，这部分需要用到相关的go sdk开发
链码直接采用contractapi，而不是shim包，因为据官方文档说shim更加初级，有可能会有奇奇怪怪的问题。
在~/HyperledgerFabric/mycodes/demo目录下存放的是链码的源代码，目前只是写了一个大致的框架。其他的服务端程序尚未开始开发。demo目录以后想起来了再改个名，比如改成demo_chaincode啥的

*建议别看中文文档，会变得不幸...直接看英文文档会更加新，也会少很多坑*