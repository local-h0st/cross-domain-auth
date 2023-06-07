# FabApp4VS
采用fabric-gateway进行开发，samples仓库中有示例代码。在code-server内一键运行会导致tcp无法监听，因此在code-server内go build，通过ssh进入CLI运行

### 与demo交互
先在test network上部署demo，成功后下一步先测试一下isExist函数能否工作，以及观察运行返回的结果。这次先直接在服务器跑，下回再放容器里跑（失败教训参考nonsense.md）开容器测试也是可以的，直接挂载目录就完了呗，不自己复制，唯一要解决的就是容器端口占用的问题、两个容器如何通信的问题。

asset-transfer-basic gateway版代码中const部分有误，将`certPath = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"`改为`certPath = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"`，这样才能顺利跑起来。

成功跑起来之后用sendMsg向server发送信息，观察server和chaincode的交互结果
```
./sendMsg -p 54321 -m "{"PID":"this is pid.","P":"this is p","Tag":"3","ID_cipher":"this is cipher.","PK_device2domain":"this is pk."}" // msg碰到空格就截断了，先不处理这个问题，先发一个没空格的试试看
./sendMsg -p 54321 -m "{"PID":"test","P":"thisisp","Tag":"3","ID_cipher":"thisiscipher.","PK_device2domain":"thisispk."}"   // 转义了，改成这样：
./sendMsg -p 54321 -m '{"PID":"test","P":"thisisp","Tag":"3","ID_cipher":"thisiscipher.","PK_device2domain":"thisispk."}'   // 返回了result为false！成功交互！
换个pid存在的
./sendMsg -p 54321 -m '{"PID":"redh3tALWAYS","P":"thisisp","Tag":"3","ID_cipher":"thisiscipher.","PK_device2domain":"thisispk."}'   // 返回结果true！太开心了！今天push了睡觉，明天再来复现和继续开发！
```

## 一些问题
这是跑在VS上的，不过反正都是跑在一个容器内，用Fabric SDK去调用chaincode，因此貌似跑在哪里都行。如果要限定必须在VS上跑，可以在这个程序里写检测，因为在chaincode里面检测比较麻烦。感觉限定在peer上跑更安全，虽然说不出哪里更安全。请求调用chaincode的话智能合约会验证请求者的代码是否被篡改，因此如果想通过修改此DApp攻击的话，调用chaincode会在chaincode端失败。虽然我还没写这个核验方法。另外本身攻击者无法左右sgx端返回结果

刚才看了一下发现服务器不支持sgx，但是Intel官方提供了simulate mod，翻了一圈没找着开启方法，另外原生sgx仅仅支持c语言。好在找到了ego!
EGo感觉能封神！

// “Intel的威胁模型最开始就不对，不应该隔离管理员”

单开一个App管理enclave，管理员配置的肯定是对的，因此最初会有一对公私钥对，每次要询问enclave时先用公钥加密一个挑战过去，如果被恶意程序替换了那么恶意程序一定不能解开挑战，serverVS没有收到解决的挑战就不认可就拒绝发送pid认证或不认可结果。可以挑战和ID一起发，app返回挑战和true/false

破案了，没有空格截断，只是我没打引号罢了

	// 这里我选择LevelDB - Fabric自带的内嵌键值存储数据库,可以直接在chaincode中使用。
	// 数据只保存在一个peer节点上,不可持久化。
	// https://www.tizi365.com/archives/411.html

		// 	val, err := db.Get([]byte(msg.PID), nil)
