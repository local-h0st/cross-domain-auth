# server app on VS
这是跑在VS上的，不过反正都是跑在一个容器内，用Fabric SDK去调用chaincode，因此貌似跑在哪里都行。如果要限定必须在VS上跑，可以在这个程序里写检测，因为在chaincode里面检测比较麻烦。感觉限定在peer上跑更安全，虽然说不出哪里更安全。

请求调用chaincode的话智能合约会验证请求者的代码是否被篡改，因此如果想通过修改此DApp攻击的话，调用chaincode会在chaincode端失败。虽然我还没写这个核验方法。

另外本身攻击者无法左右sgx端返回结果

因为需要在peer上运行这个进程，因此需要先确定可行性，即如何找出运行chaincode的peer、如何进入容器、如何安装此进程、如何与该进程交互等等。我决定先拿atcc做个测试。

翻了一圈，发现devMod下不会新开docker容器，那么测试以后再说，先把代码写出来

有个问题，fabric-sdk-go支持到2.2版的fabric，我用的是2.5版本的，不知道会不会有兼容性问题

发现在code-server里面一键运行的话tcp监听不了，如果go build之后ssh连上去运行可以监听，不知道是因为在code-server内的原因还是单纯需要build


让Claude来写这个RSA真是写不清楚啊，无语，不如我上网copy

md，mv指令直接把原来的README给覆盖了，还好GitHub上有备份

问了Claude，说开发Fabric App有两种方式，一个是使用sdk的api，还有一个是使用gateway。在官方文档翻到了使用gateway，说用gateway更简单，那就用gateway吧。在samples仓库中找到了示例代码，有些是两种方式都有，有些是只提供了gateway的交互方式。

[官方文档Gateway](https://hyperledger-fabric.readthedocs.io/en/release-2.5/gateway.html#writing-client-applications)

This API requires Fabric v2.4 (or later) with a Gateway enabled Peer. Additional compatibility information is available in the documentation:https://hyperledger.github.io/fabric-gateway/  2.5是tested，2.44+是supported，

[Gateway仓库]https://github.com/hyperledger/fabric-gateway/blob/main/pkg/client/README.md

2.4版后的gateway需要特殊的节点支持，我有点怕test network缺少特殊节点，不过看仓库说application-gateway-go可以直接用test network跑，应该没问题吧

go get 解决一切依赖问题？

先测试一下isExist函数能否工作，以及返回的内容

```
Error: failed to normalize chaincode path: 'go list' failed with: go: downloading github.com/hyperledger/fabric-contract-api-go v1.2.1: signal: killed
```

发现部署demo失败了，看日志发现在go verdor之后莫名其妙有一句go list然后go list卡住了，但是atcc却没有这个问题，但是这俩都是我自己写的啊，难道是一个vendor了一个没有vendor吗？不应该啊

好像是因为我没写func main()...虽然这个推断和日志里面的go list没有半毛钱关系但是这应该是一个明显的错误。补上再试一次！

又戳了？看一眼代码发现log包没有及时导入，自动导入后再试一次

还是失败，提示信息一模一样，为什么啊，真的是因为我没有vendor吗？vendor之后再试一下

成功了。还真是vendor的问题，参考[这里](https://hyperledger-fabric.readthedocs.io/en/release-2.5/chaincode4ade.html)

使用peer chaincode invoke和peer chaincode query测试正常，接下来测试Fabric App的交互。我的想法是开一个容器，server放在容器里跑

容器就挂载到fortest目录下吧，把必要的证书文件之类的全部放进来。另外go build最好在code-server内搞，如果在ssh里面搞的话每次都要重新下载一遍。想那些打包chaincode的都是在ssh里面搞的，就很慢，但是没办法

docker一开始跑不起来我以为是下划线的问题，还给server_vs改名为serverVS，结果发现是fortest末尾不能有/，还有就是alpine居然没有bash...
```
cd ~/HyperledgerFabric/mycodes/serverVS/
docker run -it -v ./fortest:/root/fortest -p 7000:7051 -p 54321:54321 --name vsTest --rm alpine ash
nc -l 7000 | nc localhost 7051
```
但是这端口占用的问题一直没办法解决的，端口转发也不行，还有就是居然缺少文件？我测

文件我是复制过来的，可能不行，看看原位置文件齐不齐。network down的时候是全都没有的，起个demo看看。发现有文件了，说明是动态变化的，并且我还不能改变文件位置。

也就是说我在容器里跑的可能性也没有了咯。那么以后测试会出大问题的，现在先顶一下，问题不大。

先直接在服务器上跑。

真的坑啊，照抄的配置文件居然不对，真可恶。`certPath = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"`是错的，正确的应该是`certPath = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"`，修改完后server才顺利跑起来了。
```
./sendMsg -p 54321 -m "{"PID":"this is pid.","P":"this is p","Tag":"3","ID_cipher":"this is cipher.","PK_device2domain":"this is pk."}"
```
msg碰到空格就截断了，先不处理这个问题，先发一个没空格的试试看
```
./sendMsg -p 54321 -m "{"PID":"test","P":"thisisp","Tag":"3","ID_cipher":"thisiscipher.","PK_device2domain":"thisispk."}"
```
转义了，改成这样：
```
./sendMsg -p 54321 -m '{"PID":"test","P":"thisisp","Tag":"3","ID_cipher":"thisiscipher.","PK_device2domain":"thisispk."}'
```
返回了result为false！成功交互！
换个存在的
```
./sendMsg -p 54321 -m '{"PID":"redh3tALWAYS","P":"thisisp","Tag":"3","ID_cipher":"thisiscipher.","PK_device2domain":"thisispk."}'
```
返回结果true！
太开心了！

今天push了睡觉，明天再来复现和继续开发！