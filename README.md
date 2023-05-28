# Project cross-domain-auth "AlternateWorld"

## TODO & done
* [ ] 继续写chaincode，Fabric App for verification server。Client App还没开始写
* [ ] 有必要看一看Key Concepts，以及test_network的tutorial细节

* [x] 整理项目的markdown
* [x] 去command reference看看peer chaincode invoke和peer chaincode query
* [x] 启动test_network也写成脚本，完善readme的test_network部分
* [x] 先拿atcc的chaincode部署在测试网络上
* [x] 自己写chaincode(atcc)测试，数据用my favorite songs
* [x] 重装fabric-samples
* [x] 重装服务器并恢复开发环境


## OverView of the Proj

这个项目是做跨域认证的，基于Hyperledger Fabric 2.5，采用Intel SGX作为Truetsed Execution Environment的硬件支持。

🎉首先庆祝第一阶段顺利结束！接下来就是搭环境写代码的实现阶段了。整个项目需要编写chaincode、VS上的服务端程序、PAS上的服务端程序，device上的用户程序。由于服务端程序涉及到调用智能合约，因此也属于DApp的范畴，这部分需要用到相关的go sdk开发（现采用fabric-gateway）

🔰Hyperledger Fabric👉[官方文档](https://hyperledger-fabric.readthedocs.io/en/release-2.5/)。一定要看release-2.5版的
* [fabric-samples仓库](https://github.com/hyperledger/fabric-samples)，含多项可供参考的示例代码包括chaincode和Fabric App，记得切换branch
* [Key Concepts](https://hyperledger-fabric.readthedocs.io/en/release-2.5/key_concepts.html)
* [Commands Reference](https://hyperledger-fabric.readthedocs.io/en/release-2.5/command_ref.html)
* [contract-api-go repo](https://github.com/hyperledger/fabric-contract-api-go)，内含使用`contract-api-go`编写chaincode的教程
* [fabric-gateway Guidance](https://hyperledger.github.io/fabric-gateway/)，通过gateway提供的API和chaincode交互

*建议别看任何的中文文档，会变得不幸...英文文档会更加up-to-date，直接看也会少很多坑！*

```
// 以下目录经过重命名，很多脚本跑起来可能会出问题，不过易于发现
~/HyperledgerFabric/mycodes/ ==> ~/HyperledgerFabric/codes/
~/HyperledgerFabric/myshells/ ==> ~/HyperledgerFabric/shells/
~/HyperledgerFabric/codes/send_msg_tool/ ==> ~/HyperledgerFabric/codes/sendMsg/
~/HyperledgerFabric/codes/gengerate_json_tool/ ==> ~/HyperledgerFabric/codes/genJSON/

// 以下目录经过合并
~/HyperledgerFabric/codes/genJSON
~/HyperledgerFabric/codes/sendMsg
合并入
~/HyperledgerFabric/codes/tools/
```

## 项目结构
每个README同目录下都有一个用于记录中间过程的nonsense.md
#### [MainChaincode](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/demo)
链码部分，主要提供了和账本交互的接口
#### [FabApp4VS](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/serverVS)
每一台运行chaincode的peer都需要安装此服务程序，用于和链码配套完成有关匿名身份的部分，通过fabric-gateway和链码demo交互
#### [sendMsg](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/sendMsg)
用于向指定端口发送指定消息，通常用来向FabApp4VS发送包含pid等内容的信息，以测试FabApp4VS功能是否正常
#### [TestChaincode](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/atcc)
照着教程魔改的100%能正常工作的链码，用于测试某些脚本能否正常部署这些自己开发的链码

### devMod & test_network
* [How to start devMod](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/shells/devModOn)
* [How to deploy your on test network](https://github.com/local-h0st/cross-domain-auth/blob/master/HyperledgerFabric/shells/testNetworkStart)