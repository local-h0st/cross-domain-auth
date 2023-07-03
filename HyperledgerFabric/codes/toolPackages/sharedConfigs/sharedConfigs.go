package sharedconfigs

const (
	NodeID        = "VS001"
	NodeEnclaveID = "VS001Enc"
	DatabasePath  = "./db"

	EnclaveAddr = "localhost:55555"
	EnclavePort = ":55555"
	ServerAddr  = "localhost:54321"
	ServerPort  = ":54321"

	EnclavePubkey = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCj8hW+keEOHHHLV/7BRO7I0j7a\nXAfxTvkiM8Qyex+aMQ7Ny+cavF4mWlJmdmGo9K3jHFH3LyEd2JuGPh5T0ad/O76C\nor+hX+RvgXkg0HS3MEQIwmzmNg57RSaNxzlJatXEfjpRJJ5Nc+dyA6hpYzaNj9LY\nKex5gvsGFpBMwQZyVwIDAQAB\n-----END PUBLIC KEY-----\n"
	EnclavePrvkey = "-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQCj8hW+keEOHHHLV/7BRO7I0j7aXAfxTvkiM8Qyex+aMQ7Ny+ca\nvF4mWlJmdmGo9K3jHFH3LyEd2JuGPh5T0ad/O76Cor+hX+RvgXkg0HS3MEQIwmzm\nNg57RSaNxzlJatXEfjpRJJ5Nc+dyA6hpYzaNj9LYKex5gvsGFpBMwQZyVwIDAQAB\nAoGBAKK4SZq/Qaf21X8lFIaRO4t5GcczJvL8Fkw7IxWTnNc2r+HU6slfgvcAGN73\nypCeYeSTnEsBrRXpgtun1gQNh/cqvnJU9uCpY/PVuk14vE+lYLhKkX/GAWfsmPs+\n2AUWJZeAVCJuixh9E9jnDSz+X8IWNC77cqZq8CIY/5M+nuCBAkEA1TaANVdqNBPL\n4phbC2dVFddCADWHNyaHbOrVy+bxD0x00CqyDupwvc6QMVbqpLydbbdJplZ3g6mk\nhy1HbpAJlwJBAMTYjTMSk/bLwA6D4SFGw1NyVLMOn9I6bnZzB8ryrbBdq6vbx0vV\nvGIsPNA6bKFTgUJb5DepWRMPisL02qS+dUECQQDSR0Uc1pC0uc18Nlycm5XLy5eZ\nUzF/D+3CWrzus16Ngw81+tXPZiI44E9PifQy8p6lBX6KoX6PiLDubJalkUMTAkBB\n55buwIuVl4YH1hOsBnsjFyZQhNbxleqh8cVsJ3ALmnD9qynAtCDMZa8+sDDqmoCu\nbQGtuR8/iHaW60/A1JuBAkEAgKuNrksiWi0h0KTFnassKgeaBUd2MociEK6hmKwI\nwi7kjPNHeaa1MqJMUQLhhYv33m5xuNFxIip2LTcXeJ+/5g==\n-----END RSA PRIVATE KEY-----\n"
	ServerPubkey  = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDQlXmFEiNzbO0iHjdYIUPvbWmq\nPmMJcrGVLjRUrr2HtURh9lcrGsti1r4BesFcuS+QAzBlFZsp50Ytae0snr26jnFL\nOpBGscCDLsyPrL3dlUnWGnQY5SOFjvVpAjsuc16W0TXXzdoaW6yZwX+tKd2yLkgb\ncL0alZeTI1v8lJN9YQIDAQAB\n-----END PUBLIC KEY-----\n"
	ServerPrvkey  = "-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQDQlXmFEiNzbO0iHjdYIUPvbWmqPmMJcrGVLjRUrr2HtURh9lcr\nGsti1r4BesFcuS+QAzBlFZsp50Ytae0snr26jnFLOpBGscCDLsyPrL3dlUnWGnQY\n5SOFjvVpAjsuc16W0TXXzdoaW6yZwX+tKd2yLkgbcL0alZeTI1v8lJN9YQIDAQAB\nAoGBAL7k/fk+l3lM6F3AP7CFiVI31Wu8exEricDZL4WNAuKPkA0D0dUeSaOkmvJp\nsUu2JARuFr18n6wjAMQRXMHoagQKt4zKB4Kcv6vNKNC+DLQVnd9WvujwDYChlty7\nOd+vuK5tXBIUi5FwF/QjHhOIj8EKhb228sxqXshEViuHO3dtAkEA4O6YM2cZp2Za\ngJODiRcXOjXoSkSVMqYfhARGiWRvJ0DuGAkksBAzTSxPeu0N6yEPv2Cddw2HCa1+\n24+Mqm6PPwJBAO1k0wYePhClgbTh/6sUCx5lyK+l+oTJ7H5heb0LthNX7n/B1xKv\nNr2JZtRRQ+QSRK+oTztVDSd87C0jRgITa18CQHdLU4F/nsV/rWQf2FUu3+zJhmdN\nNGvmWzSjJ93aXHFPKHeq8cBG905ov8aMTyNzJ2zyitEHZaUmVO+RlKMXe/UCQQCR\nSUxxCR849uHr/wiG/kxTvT1WaoFotV/cdPGZhkpXilA3tj1XfQ5Gb4oUVOv08E1D\nKAHdsQ7M5QJyGY1mBdaHAkAf5Y1dfFOOTQqZNkSWgm0lFfX1tfUPMD1/XnfXBgxr\n4VvDGz66MsAvh1v4qYw8GcKoGxetTWS8yzIYiBMgqD5l\n-----END RSA PRIVATE KEY-----\n"
	AdminPubkey   = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDVDTjZaYpCk5R8RNcq86rfrbxn\n6pss6d+MlodAu/UdgfstJzQK6X9gqWceCyAmm30CHdRZbC9E39E2U13uF9NX2oTc\npdC0XLRAW+rskd/jCGB0zQSfe4ktzKu0a5Cp1pwLWyZtCU7re6A2bhVeOh9JUNAD\nFBiemhYl3JIfetAVDwIDAQAB\n-----END PUBLIC KEY-----\n"
)

// AdminPrvkey = "-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQDVDTjZaYpCk5R8RNcq86rfrbxn6pss6d+MlodAu/UdgfstJzQK\n6X9gqWceCyAmm30CHdRZbC9E39E2U13uF9NX2oTcpdC0XLRAW+rskd/jCGB0zQSf\ne4ktzKu0a5Cp1pwLWyZtCU7re6A2bhVeOh9JUNADFBiemhYl3JIfetAVDwIDAQAB\nAoGBAKjAjlT3GcJeLvC3fk7RLnl5nZAZ7cuHe8BZwsvtlNtIh3FeagRyqqgfxkOv\nwEmUQ1IX2ojx/gbp2UbUhcP/LzAmKN0CWFlpOF+XodndhPVSMYMOajua6+IWXWIO\nKk70FgFNSTKZ2x+dua6P7UgDtXCBs8EZo/RtTF+ieQ7tBSWxAkEA++TqC+9TWZWo\n6ThnfY+VJF4/xrfZL3xbogZUhMutkpJokQSVZTRsciXLhMvsrq1tFqspRUDE9BlQ\nXjFstizu5QJBANiGOnz5dM9PHVnk4ko41eOUlMHm5bIll6geYms2CDzfSUZCJ6Lu\nXusNuTUduuH8xdoJD+knLgWoLuPHMy1TQOMCQGSpkFaAp6BvTHcXEVR+Iq3L9FSn\nd+WgHsZbHT+MXarrU1pQqJsvHf9n1zMUg1sy9xtN/0orngmmbBWYTsdmoXkCQQDU\n9njScNzmBg99WjUEAZDGLV5+tKaZKIZYkcIFZviFPqyoUOsBQujS0gWm653jJiZH\nhIBEtwd6AuhTmpqIawk3AkEAnpe91sUA0SouCWIzgI2EQu4Z6pEgpqb6AHBuuTfF\nvYkZsveKy8edEi2hGzo+U6nb6S4VYnm/rOQKum1Mx5GxEQ==\n-----END RSA PRIVATE KEY-----\n"

// var node_id string = os.Getenv("SERVERID")
// selfID = os.Getenv("SERVERID")
// fmt.Println("PUBKEY ==> ", string(bytes.Replace(PUBKEY, []byte("\n"), []byte("\\n"), -1)))

// PRVKEY, PUBKEY = myrsa.GenRsaKey()
// 我在想如何实现运行初动态生成prvkey和pubkey，我现在采取直接指定的办法连接vs和enclave
// 可以双方保存管理员公钥，双方分别生成公私钥匙，打印公钥，管理员用管理员私钥签名后分别发送给en和vs

// enclave只允许和vs通信，因此同步黑名单也应该设计为发送给vs，vs再转发给enclave
// 对黑名单是否过期的判断应该在vs中完成
// enclave一定要用EGo跑

// domaininfo map 并发不安全，可能会出问题

// msgs消息结构也需要调整，第一次体会到了自己造的shit mont是多么难改，，这个以后再改吧

// 结构内可以加入随机数防止截获密文重放攻击

// serverVS处理黑名单
