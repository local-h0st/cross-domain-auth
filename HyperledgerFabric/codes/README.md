## ideas to achieve
1. 第二阶段判断pid是否valid并返回结果最好让链码来做而不是交给PASB来做
2. ID等可以如下传入：`os.Getenv("SERVERID")`
3. 关于公私钥，最好是`myrsa.GenRsaKey()`生成，最好实现运行初动态生成prvkey和pubkey，目前采取直接指定的办法连接vs和enclave。可以双方保存管理员公钥，双方分别生成公私钥匙，打印公钥，管理员用管理员私钥签名后分别发送给en和vs。enclave的key最好是在enclave内部生成。
4. enclave要用EGo跑
5. 注意防止某些全局变量并发不安全
6. msgs某些结构的字段存在冗余
7. msgs结构内可以加入随机数防止截获密文重放攻击
8. initLedger单独放一个程序就行了，这样就不用每次重启整个区块链系统了

## pre-works
在sharedConfigs中预先配置好三对公私钥、端口等信息然后编译。之后搭建区块链账本环境并部署合约demo。当将demo部署上链后，由另一个程序完成initLedger。随后启动enclave，enclave监听EnclavePort，此时enclave内容为空，enclave仅接受来自server的消息

接下来启动server，server和enclave成对出现。server先初始化leveldb，连接gateway，调用SubmitTransaction向账本写入此Node节点信息。然后向enclave发送`testingConnection`信息，之后立刻开始监听ServerPort。

enclave收到`testingConnection`消息后向server返回`testingConnection`消息，server收到后输出connection ok，然后继续监听。此时说明server和enclave双向通信没有问题。

接下来admin进程向server发送`syncDomain`消息，向server同步各个域的信息。

接下来启动域pas，pas已知节点服务器的地址，首先向server发送`updateBlacklistTimestamp`消息，即使黑名单此时为空。因为域加入必须在账本上有记录可以查询，否则会影响后续逻辑。

接下来可以进行正常的匿名认证操作了。

__pre-works应该没有大问题，有的话记得微调__

## records on the chain's ledger
索引包含三种：pid、NodeID、domainName
* pid：此记录是pid核验记录，所有字段均为原本含义。通常pid为不规则hash避免和以下两种特殊pid碰撞
* NodeID：此记录为节点服务器的信息，Valid字段是一个json字符串，格式为msgs.ServerRecord{}
* domainName：此记录是域黑名单最后更新时间戳，Valid字段是时间戳

## 消息交互格式 伪代码
### 发送端
```
bm := msgs.BasicMsg{}
bm.Method = "xxx"
bm.SenderID = "xxx"
bm.Content = json.Marshal(msgs.OtherTypeMsg{
    XXX: xxx
    ...
})  // or bm.Content = xxx if simple.
bm.GenSign(prvkey)
jsonstr, _ := json.Marshal(bm)
cipher := myrsa.EncryptMsg(jsonstr, target_pubkey)
sendMsg(target_addr, cipher)
```
### 接收端
```
// cipher received
jsonstr := myrsa.DecryptMsg(cipher)
bm := json.Unmarshal(jsonstr)
bm.VerifySign()
switch bm.Method{
    case "xxx":
        ...
    case "xxx":
        ...
    default:
        ...
}
```
按照不同Method采用不同格式对Content进行反序列化，某些简单的情况下Content甚至不是json字符串，可以直接处理

## msgs.BasicMsg.Method
### pas -> server
* server完成黑名单是否过期的检验


### server -> enclave

### enclave -> server


## 究极测试流程(不完整)
1. 部署链码
2. 启动enclave
3. 启动serverVS
4. admin向serverVS发送所有domain信息
5. 源域pas向所有serverVS分别发送fragment，所有serverVS收到后先存储fragment
6. tag serverVS如果没有目标域黑名单信息或者已经过期，则会触发黑名单同步，继续；如果黑名单无需更新，跳转10
7. serverVS临时保存此份fragment，随后向目标域pas提出黑名单同步请求
8. 目标域pas发回两次加密后的黑名单，serverVS收到后转发给enclave
9. enclave同步完成后发送ok提示给serverVS，serverVS收到ok后，将临时保存的、目标域是发来黑名单pas所在域的所有fragment全部发给enclave核验
10. 把fragment发给enclave核验
11. serverVS收到一个或多个核验结果，但是都一样，以下讨论对每一个收到的结果如何处理
12. 将每一个pid的核验结果都写入区块链账本，通知源域pas
13. 设备直接向目标域pas提出认证请求，目标域pas收到后查询账......






## 可以用来扯的话
### interactSGX

在内存中维护一个黑名单slice，只有两个功能，一个是接受加密后的ID，check是否在黑名单内，另一个就是请求同步黑名单数据

运行时先生成密钥对，这一对keys是在enclave内部生成的，安全性很高，公布公钥pk。

PAS发送给VS的消息中就包含了pk加密的ID，由于VS无法得到sk因此安全

serverVS发送来verifyID,pk(json(challenge,update_flag,domain,pk(ID)),sign)，challenge是为了确定interactSGX没有被恶意程序替代。update_flag代表是否需要更新黑名单，这一个数据来自于VS查询链上账本，如果PAS的黑名单更新，则调用chaincode更新update_flag为true，如果PAS收到了来自interactSGX的黑名单同步请求，则调用chaincode更新update_flag为updated(false),verifyID表示此操作为验证ID,sign是代表确实VS发来的请求，而不是其他的恶意程序发来的请求。这样可以避免恶意程序询问某个ID是否在黑名单上造成信息泄漏

至于为什么interSGX不采用sign的方式表明自己的身份，，好像也可以？不过这样的话服务器会多一次解密过程。那就改成签名吧。