

问了Claude，说开发Fabric App有两种方式，一个是使用sdk的api，还有一个是使用gateway。在官方文档翻到了使用gateway，说用gateway更简单，那就用gateway吧。在samples仓库中找到了示例代码，有些是两种方式都有，有些是只提供了gateway的交互方式。

[官方文档Gateway](https://hyperledger-fabric.readthedocs.io/en/release-2.5/gateway.html#writing-client-applications)

This API requires Fabric v2.4 (or later) with a Gateway enabled Peer. Additional compatibility information is available in the documentation:https://hyperledger.github.io/fabric-gateway/  2.5是tested，2.44+是supported，

[Gateway仓库]https://github.com/hyperledger/fabric-gateway/blob/main/pkg/client/README.md

2.4版后的gateway需要特殊的节点支持，我有点怕test network缺少特殊节点，不过看仓库说application-gateway-go可以直接用test network跑，应该没问题吧

采用fabric-gateway进行开发。在某个角落找到了一句话：This API requires Fabric v2.4 (or later) with a Gateway enabled Peer，并且有一张表格显示2.5是tested，2.44+是supported。在samples仓库中有示例代码，有些是sdk和gateway两种方式都有，有些是只提供了gateway的交互方式。



测试Fabric App的交互。我的想法是开一个容器，server放在容器里跑

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