# Project cross-domain-auth "AlternateWorld"

## TODO & done
* [ ] 继续写chaincode，Fabric App for verification server。Client App还没开始写
* [ ] 有必要看一看Key Concepts，以及test_network的tutorial细节

[没什么用的碎碎念](https://github.com/local-h0st/cross-domain-auth/blob/master/records.md)

* [x] 去command reference看看peer chaincode invoke和peer chaincode query
* [x] 启动test_network也写成脚本，完善readme的test_network部分
* [x] 先拿atcc的chaincode部署在测试网络上
* [x] 自己写chaincode(atcc)测试，数据用my favorite songs
* [x] 重装fabric-samples
* [x] 重装服务器并恢复开发环境


## OverView of the Proj

这个项目是做跨域认证的，基于Hyperledger Fabric 2.5，采用Intel SGX作为Truetsed Execution Environment的硬件支持。

🎉首先庆祝第一阶段顺利结束！接下来就是搭环境写代码的实现阶段了。

🔰Hyperledger Fabric👉[官方文档](https://hyperledger-fabric.readthedocs.io/en/release-2.5/)  （一定要看release-2.5版的，版本不一样冲突的太多了）

[Key Concepts](https://hyperledger-fabric.readthedocs.io/en/release-2.5/key_concepts.html)

[Commands Reference](https://hyperledger-fabric.readthedocs.io/en/release-2.5/command_ref.html)

[contract-api-go仓库](https://github.com/hyperledger/fabric-contract-api-go)，内含使用`contract-api-go`编写chaincode的教程

[fabric-sdk-go仓库](https://github.com/hyperledger/fabric-sdk-go)，开发Fabric App必看

[fabric-samples仓库](https://github.com/hyperledger/fabric-samples)，含多项可供参考的示例代码包括chaincode和Fabric App，记得切换branch

还有一个github.io的Fabric[中文文档](https://hyperledger.github.io/)（欸好像不是这个网址），不过看着好像没什么用

*建议别看任何的中文文档，会变得不幸...直接看英文文档会更加新，也会少很多坑*

## 项目结构
* [chaincode](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/mycodes/demo)
* [Fabric app for verification server](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/mycodes/server_vs)

整个项目需要编写chaincode、VS上的服务端程序、PAS上的服务端程序，device上的用户程序。
由于服务端程序涉及到调用智能合约，因此也属于DApp的范畴，这部分需要用到相关的go sdk开发
链码直接采用`contractapi`，而不是`shim`包，因为据官方文档说shim更加初级，有可能会有奇奇怪怪的问题。
在`~/HyperledgerFabric/mycodes/demo`目录下存放的是链码的源代码，目前只是写了一个大致的框架。其他的服务端程序尚未开始开发。demo目录以后想起来了再改个名，比如改成demo_chaincode之类的

## test_network
[deploy your chaincode on test_network](https://github.com/local-h0st/cross-domain-auth/blob/master/HyperledgerFabric/myshells/testNetworkStart)

编写Fabric App for verification server时，sdk需要和chaincode交互。我的demo chaincode之前是在devMod上跑的，没有任何的peer，sdk交互起来好像有点问题，sdk和chaincode交互的需要一个真实的网络而不是devMod。因此我尝试将atcc部署在test_network上，先看看能不能正常交互，反正tutorial说可以拿来部署其他的chaincode。之前部署直接导致服务器卡死，重装才解决，心有余悸，这次也是做好了比较充分的备份才敢第二次尝试。好在这次成功了！

### 关于peer chaincode
之前一直奇怪`peer chaincode invoke`和`peer chaincode query`有什么区别，直到我查阅了Cmd Ref的[这里](https://hyperledger-fabric.readthedocs.io/en/release-2.5/commands/peerchaincode.html)，里面提到`invoke`会try to commit the endorsed transaction to the network，但是`query`不会generate transaction.

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

官方教程部署的链码位于`～/HyperledgerFabric/fabric/integration/chaincode/simple/cmd`，我整合命令后的sh能够跑官方的chaincode，随后用我自己写的chaincode测试，也就是atcc，能够正常工作。

开启devMod部署atcc链码后，测试链码功能用的命令如下：
```
// -c, --ctor string =>  Constructor message for the chaincode in JSON format (default "{}")

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