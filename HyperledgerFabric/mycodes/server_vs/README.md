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