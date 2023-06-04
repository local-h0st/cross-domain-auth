# interactSGX

在内存中维护一个黑名单slice，只有两个功能，一个是接受加密后的ID，check是否在黑名单内，另一个就是请求同步黑名单数据

运行时先生成密钥对，这一对keys是在enclave内部生成的，安全性很高，公布公钥pk。

PAS发送给VS的消息中就包含了pk加密的ID，由于VS无法得到sk因此安全

serverVS发送来verifyID,pk(json(challenge,update_flag,domain,pk(ID)),sign)，challenge是为了确定interactSGX没有被恶意程序替代。update_flag代表是否需要更新黑名单，这一个数据来自于VS查询链上账本，如果PAS的黑名单更新，则调用chaincode更新update_flag为true，如果PAS收到了来自interactSGX的黑名单同步请求，则调用chaincode更新update_flag为updated(false),verifyID表示此操作为验证ID,sign是代表确实VS发来的请求，而不是其他的恶意程序发来的请求。这样可以避免恶意程序询问某个ID是否在黑名单上造成信息泄漏

至于为什么interSGX不采用sign的方式表明自己的身份，，好像也可以？不过这样的话服务器会多一次解密过程。那就改成签名吧。

### sgxIntercat和serverVS消息交互格式
```
type basicMsg struct {
	Method     string
	Content    string
    Signature  string
}
```
每一条消息发送时：
```
basicmsg = ...
genSign(basicmsg)
jsonstr = json.Marshal(basicmsg)
ciphertext = encrypt(jsonstr) // 使用接收方的pubkey
send(ciphertext)
```
每一条消息收到时：
```
ciphertext = receive()
jsonstr = decrypt(cipher) // 使用自己的prvkey
basicmsg = json.Unmarshal(jsonstr)
verifySign(basicmsg)
// 根据不同basicmsg.Method选择不同的函数处理
switch basicmsg.Method {
    ...
}
```
Content的内容一般是json字符串，根据method的不同json结构不同，处理函数也不同
在verify ID的result中，content的内容不是json字符串，而是字符串"valid"或"invalid"
在blacklistNeedUpdate和needPubkey消息中content地址

md还要自己修改go.mod才能正常go build