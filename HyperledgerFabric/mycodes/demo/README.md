# 开发记录

其实匿名生成阶段本质上就是调用智能合约，所以可以直接写成接口
按照原先方案，不同VS收到的message不同，这个在智能合约的实现上会出现需要为不同背书节点传入不同参数的情况
这里可以修改为传递相同的参数，每一台都比较tag和node_id，只有相同才进入enclave验证

背书节点都要运行一遍chaincode，但是不同的背书节点环境变量可以不同。通过设置环境变量，可以让tag节点和非tag节点采取不同的动作

这部分chaincode是用来参与匿名认证过程、管理pid的。应该还需要有一份chaincode，用来管理各种其他公开信息。要不这两类功能全部写在同一份chaincode源码里面吧，好像这样也行。

https://youtube.com/shorts/y0cxkflRHto?feature=share

貌似必须要保证每一个peer上运行一遍chaincode后的结果必须完全一致，否则不能达成共识

#### func (s *SmartContract) HandleMsgForPseudo(ctx contractapi.TransactionContextInterface, cipher_text string) error
该函数传参的cihper_text原文为message的JSON字符串，采用RSA方案，用node_id为tag的背书节点的公钥加密得到cipher_text


这里本来的想法就是在VS，也就是背书节点上再写一个DApp，专门用来监听connection并处理，最终调用智能合约来更新账本，后来想想这样好像不安全，因为恶意的DApp也能够调用智能合约，另外直接写成chaincode接口好像也没有什么问题？不对有问题！

仅仅在tag node上核验，如果写进chaincode上有共识性问题的，因为仅有tag node的账本会更新，其他node不会更新自己的账本，导致账本不一致

何况发给不同peer的数据也不一样，加密密钥也不一样，而chaincode只调用一次，这样显然不能解决问题

基于blockchain的意义何在？貌似仅仅是方便数据管理了。

我还是想在peer节点上单独跑一个普通的server进程，如何解决恶意DApp问题呢？智能合约被调用前先核验DApp的签名是否被官方认可，也就是只运行经过认可的DApp调用智能合约，因为运行在peer上的DApp都是官方开发的，因此可以这么做


#### func verifyIDinSGX(cipher string) bool 
情况比我预想的要好，golang也能够直接和SGX交互，不需要单独拿C写一个程序然后外部调用。

#### about Level DB
InitLedger时在公开的账本上指明了database的地址是./db，但是同样的地址在不同peer上的库文件内容是不同的，因此满足需求

### 测试数据
测试用json字符串的生成可以用../generate_json_tool/目录下的工具
```
{"PID":"this is pid.","P":"this is p","Tag":"3","ID_cipher":"this is cipher.","PK_device2domain":"this is pk."}


```