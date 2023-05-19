# 开发记录

其实匿名生成阶段本质上就是调用智能合约，所以可以直接写成接口
按照原先方案，不同VS收到的message不同，这个在智能合约的实现上会出现需要为不同背书节点传入不同参数的情况
这里可以修改为传递相同的参数，每一台都比较tag和node_id，只有相同才进入enclave验证

https://youtube.com/shorts/y0cxkflRHto?feature=share

#### func (s *SmartContract) HandleMsgForPseudo(ctx contractapi.TransactionContextInterface, cipher_text string) error
该函数传参的cihper_text原文为message的JSON字符串，采用RSA方案，用node_id为tag的背书节点的公钥加密得到cipher_text


#### func verifyIDinSGX(cipher string) bool 
情况比我预想的要好，golang也能够直接和SGX交互，不需要单独拿C写一个程序然后外部调用。