## 流程

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





## 究极测试开始
先start test network部署demo，部署真有点慢
以非enclave模式启动sgxInteract
启动serverVS
输出信息暂时正常
serverVS输出的pubkey和sgxInteract输出的pubkey一致
测试不应该向sgxInteract发送任何消息
运行genPayloadForServerVS
serverVS两条消息都json unmarshal failed了
但是没有panic，说明不是公钥的问题
其他全停了，chaincode先跑着无所谓

第二次测试重新按照上面的顺序来
部署chaincode直接卡死
太有所谓了一定要先停了chaincode
serverVS输出[addPseudoRecordToLedger] result ==> []
没问题肯定是nil，demo里面就这么写的
好怪不应该有问题的，可能是当时没有停chaincode的缘故
全停了重开试试
没用，再查询试试
也没用
addPseudoRecordToLedger不执行啊
算了下次再说今天先睡觉

我靠我知道问题了！我去gateway-go的示例代码看了一下，写入Ledger应该是Submit而不是Evaluate，我就是chaincode怎么会写错呢，账本不修改只能serverVS出问题！【来自iPhone 半夜在GitHub上翻看自己代码发现不对劲然后在iPhone的Safari上登陆codeserver来这个markdown里面记录想法然后心满意足地睡觉明天准备毛概期末考试】

成功查到记录！
先commit and push
但是不显示pid也就是key，需要稍稍调一下输出

管理员在预设阶段自己需要生成一堆密钥，用于溯源时的身份认证。
溯源时管理员向serverVS发送消息，用自己的私钥签名，serverVS收到消息就核验签名，如果通过那么就返回fragment

第二阶段判断pid是否valid并返回结果最好让链码来做而不是交给PASB来做

发现true id解码不对，会不会是`[]byte`外面再套一个`[]byte`的原因？
破案了，不知道脑子怎么想的在serverVS里面又对CipherID加密了一回，，无语

///////////
nice EGo 能用！
之前的genPayloadForServer VS移到tools里面去了
打算写个pas，再写个admin
先push了再说