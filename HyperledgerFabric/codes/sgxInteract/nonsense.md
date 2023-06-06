测试的时候才发现RSA不支持任意长度消息，所以还需要自己写分组，因此重写了加密解密签名和验证函数，分别为`EncryptMsg()``DecryptMsg()``SignMsg``VerifySign`

sendMsg可以发送`string`或者`[]byte`，但是这两者是不一样的，发送`string`收到的`[]byte`就是`string`对应的`[]byte`，但是发送`[]byte`是真的发送一段长得和`[]byte`一模一样的`string`，，很难评。就因为这个原因我之前加密完的数据发过来直接解密不了panic

说到panic，以后真得把rsa里面的panic给catch了，要不然一个错就panic系统还跑不跑了

把消息结构和rsa分出去之后，md还要自己修改go.mod才能正常go build


另外发现一跑./sgxInteract服务器过几分钟CPU就吃满ssh挂掉，我以为是什么东西死循环了还是怎么的耗了很多资源，我寻思自己写的代码也没这问题吧，结果一查发现是Claude干的好事，它给的监听代码for里面放了个go handleConn，一整个无语，资源不被吃满才怪。这样的话serverVS里面也是这样，我还没改代码。

调试UpdateServerPubkey时发现也是解密失败，一输出map里的pubkey发现是空我才意识到basicMsg里面的SenderID功能和UpdateServerPubkeyMsg里面的ServerID功能好像重合了，前期测试阶段建议取消ServerID这个弱智项，直接以SenderID来统一，以后有需要了再改

发现生成密钥好像不是随机的。不对是随机的，只是前面几行和最后几行一样而已，中间是不一样的。

挺好笑的，对方没有自己的pubkey那是怎么给自己发的消息，不知道我写NeedPubkey时精神状态到底是怎么样的，笑死

可以写成panic就给出自己的pubkey。不对可以改一下，解密前先unmarshal，如果err就解密。以后再说


同步blacklist还没写！
还剩下Verify没验证，但是需要需要先完成同步黑名单操作，还得写个接受请求同步的程序，不过不是必须，因为可以直接send

同步blacklist完成！verify功能验证通过，先commit和push，之后把genPayload移出来