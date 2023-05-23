package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	// 解析msg_text到msg
	type message struct {
		PID              string
		P                string // // 门限签名技术的那个点，我也不知道用什么格式存储
		Tag              string // node_id
		ID_cipher        string
		PK_device2domain string // 不知道公钥是这哪个类型。另外是建议存到数据库，还是存到内存算了？
	}

	byteJSON, err := json.Marshal(message{
		"this is pid.",
		"this is p",
		"3",
		"this is cipher.",
		"this is pk.",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(byteJSON))
}
