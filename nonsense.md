# nonsense in project root ~
这里主要是记录一些开发时的中间过程（倒着更新）

## useless repos & links
* [fabric-sdk-go仓库](https://github.com/hyperledger/fabric-sdk-go)是用的传统sdk，对gateway方式来说参考意义不大。
* Fabric[中文文档](https://hyperledger.github.io/)（欸好像不是这个网址），github.io的，不过看着好像没什么用

## misc
让Claude来写这个RSA真是写不清楚啊，无语，不如我上网copy。md，mv指令直接把原来的README给覆盖了，还好GitHub上有备份

根据这个文档[devMod](https://hyperledger-fabric.readthedocs.io/en/release-2.5/peer-chaincode-devmode.html)，按照步骤写了4个shell，放在devModOn目录下，并且把示例代码替换成了自己的chaincode，能开起来，但是init失败，可能是语法问题。在[这里](https://hyperledger-fabric.readthedocs.io/en/release-2.5/commands/peerchaincode.html)这里的一个角落找到了contract-api版本的peer chaincode -c格式，不过好像暂时没啥用？不对啊替换成自己的chaincode了凭什么用原来的peer chaincode invoke能正常工作啊。可能是因为之前go build出来的simpleChaincode没删，后来换成我的之后go build出问题了，没能成功build覆盖原来的。【图片】（哈哈哈哈这里没有图，反正就是InitLedger后QueryRecord成功，阿里云盘有一份备份，时间是2023.5.19）成功了呗！接下来：修改了3.sh，每次都进入我的chaincode目录下go build，之后mv到fabric目录下。可以把1.sh中rm simpleChaincode写到3.sh中，再添加一个功能，用于指定我自己写的chaincode的目录，最后全部推送至Github！

Fabric不同版本的问题：1.4链码安装在peer上，2.0链码容器独立运行。在release-2.5的官方文档写的是：Developers can use the network to test their smart contracts and applications.

官方文档一定要看2.5版本的，版本不一样冲突的太多了