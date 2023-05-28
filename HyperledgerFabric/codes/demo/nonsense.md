# nonsense in demo


### 关于peer chaincode
之前一直奇怪`peer chaincode invoke`和`peer chaincode query`有什么区别，直到我查阅了Cmd Ref的[这里](https://hyperledger-fabric.readthedocs.io/en/release-2.5/commands/peerchaincode.html)，里面提到`invoke`会try to commit the endorsed transaction to the network，但是`query`不会generate transaction.

这和后来fabric-gateway里面query和submitTransaction（好像是这个）差不多

### misc(in reversed order)
编写Fabric App for verification server时，sdk需要和chaincode交互。我的demo chaincode之前是在devMod上跑的，没有任何的peer，sdk交互起来好像有点问题，sdk和chaincode交互的需要一个真实的网络而不是devMod。因此我尝试将atcc部署在test_network上，先看看能不能正常交互，反正tutorial说可以拿来部署其他的chaincode。之前部署直接导致服务器卡死，重装才解决，心有余悸，这次也是做好了比较充分的备份才敢第二次尝试。好在这次成功了！

链码直接采用`contractapi`，而不是`shim`包，因为据官方文档说shim更加初级，有可能会有奇奇怪怪的问题。
在`~/HyperledgerFabric/mycodes/demo`目录下存放的是链码的源代码，目前只是写了一个大致的框架。其他的服务端程序尚未开始开发。demo目录以后想起来了再改个名，比如改成demo_chaincode之类的