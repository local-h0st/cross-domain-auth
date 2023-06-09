# Project cross-domain-auth - AlternateWorld
🎉初赛作品提交正式截止，目前结果尚未揭晓。
🎉首先庆祝第一阶段校赛顺利结束！
## OverView
这个项目是做跨域认证的，基于Hyperledger Fabric 2.5，采用Intel SGX作为Truetsed Execution Environment的硬件支持。整个项目需要编写chaincode、VS上的服务端程序、PAS上的服务端程序，device上的用户程序Client App。由于服务端程序涉及到调用智能合约，因此也属于DApp的范畴，这部分需要用到相关的go sdk开发（现采用fabric-gateway）

__TODO:__
* [ ] 合并README，梳理逻辑，整理各种Method（server源码中出现的已经整理完毕，接下来整理其他的）
* [ ] 垃圾话删一些，有用的别删，坑别删，万一以后有用
* [ ] 继续开发chaincode，FabApp4VS，命名规则最好统一
* [ ] 有必要看一看Key Concepts，以及test_network的tutorial细节，另外目录重命名了，遇到问题了再慢慢改

__References:__
* 🔰Hyperledger Fabric👉[官方文档](https://hyperledger-fabric.readthedocs.io/en/release-2.5/)。一定要看release-2.5版的！
* [fabric-samples仓库](https://github.com/hyperledger/fabric-samples)，含多项可供参考的示例代码包括chaincode和Fabric App，记得切换branch
* [Key Concepts](https://hyperledger-fabric.readthedocs.io/en/release-2.5/key_concepts.html)
* [Commands Reference](https://hyperledger-fabric.readthedocs.io/en/release-2.5/command_ref.html)
* [contract-api-go repo](https://github.com/hyperledger/fabric-contract-api-go)，内含使用`contract-api-go`编写chaincode的教程
* [fabric-gateway Guidance](https://hyperledger.github.io/fabric-gateway/)或[gateway repo](https://github.com/hyperledger/fabric-gateway/blob/main/pkg/client/)，通过gateway提供的API和chaincode交互
* [gateway readthedoc](https://hyperledger-fabric.readthedocs.io/en/release-2.5/gateway.html#writing-client-applications)
* [EGo document](https://docs.edgeless.systems/ego/knowledge/model)
* [EGo repo](https://github.com/edgelesssys/ego/tree/master)

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
### 系统主体
* [MainChaincode](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/demo)：链码部分，主要提供了和账本交互的接口
* [FabApp4VS](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/serverVS)：每一台运行chaincode的peer都需要安装此服务程序，用于和链码配套完成有关匿名身份的部分，通过fabric-gateway和链码demo交互
* [sgxInteract](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/sgxInteract)：用于和FabApp4VS配套，负责可信执行
* [toolPackages](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/toolPackages)：package包含统一的消息结构和RSA的实现

### 测试环境
* [devMod](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/shells/devModOn)：用于测试链码
* [startTestNetwork](https://github.com/local-h0st/cross-domain-auth/blob/master/HyperledgerFabric/shells/testNetworkStart)：用于测试整个系统

### 测试工具
* [genPayloadForServerVS](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/genPayloadForServerVS)：用于模拟各种设备发送信息以测试系统各个功能能否正常完成
* [sendMsg](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/tools/sendMsg)：用于向指定端口发送指定消息，通常用来向FabApp4VS发送消息，测试FabApp4VS功能是否正常
* [genJSON](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/tools/genJSON)：产生指定格式的json字符串用于测试，通过sendMsg发给FabApp4VS
* [atcc](https://github.com/local-h0st/cross-domain-auth/tree/master/HyperledgerFabric/codes/atcc)：照着教程魔改的100%能正常工作的链码，用于测试某些脚本能否正确部署这些自己开发的链码
