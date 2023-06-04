package main

import (
	"encoding/json"
	"fmt"
	"msgs"
	"myrsa"
)

func main() {
	PRVKEY, PUBKEY := myrsa.GenRsaKey()
	pk := "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDYiWYP6r5aVP7NBQkxyUnr8Pny\n2lZ3NFrWZHF4PQnSvUdlQ43zxzqZ33szAmi+GD2fnAQb2PSgf/+zzFcViqTgHQx0\nHyAHXdf39cVS6REVYtm06Iy/yRIjcWSUNasAg/bD/QKMzuNWmZhTVbHFbahzXXg2\nkKRfsE6+Hr7Ncr0wrQIDAQAB\n-----END PUBLIC KEY-----\n"
	// add key
	basic_msg := msgs.BasicMsg{
		Method:    "addServerPubkey",
		SenderID:  "test",
		Content:   nil,
		Signature: nil,
	}
	add_pubkey_msg := msgs.AddServerPubkeyMsg{
		ServerID:     "test",
		ServerPubkey: PUBKEY,
	}
	basic_msg.Content, _ = json.Marshal(add_pubkey_msg)
	fmt.Println(string(basic_msg.Content))
	data, _ := json.Marshal(basic_msg)
	fmt.Println(string(data))
	fmt.Println()

	a := msgs.BasicMsg{}
	json.Unmarshal(data, &a)
	if a.Signature == nil {
		fmt.Println("sig nil.")
	} else {
		fmt.Println("bad!")
	}
	b := msgs.AddServerPubkeyMsg{}
	json.Unmarshal(a.Content, &b)
	fmt.Println(b.ServerID, b.ServerPubkey)
	fmt.Println()
	fmt.Println()
	fmt.Println(pk)
	myrsa.UnitTest()

	fmt.Println()
	fmt.Println()
	data = []byte("sdfghfvfwvwvwv")
	fmt.Println(data, len(data))
	data = myrsa.RsaEncrypt(data, PUBKEY)
	fmt.Println(data, len(data))
	data = myrsa.RsaDecrypt(data, PRVKEY)
	fmt.Println(data, len(data))

	fmt.Println([]byte{})
	fmt.Println([]byte{})
	data = []byte("sdfgfghgvhwbvuwbwvbwvbwibwvbqiiqbvqivqbdqvbyqyqbuyxbqwvbwvbwibwvbqiiqbvqivqbdqvbyqyqbuyxbqcbqybxacbyqbqybxcyuwvbwvbwibwvbqiiqbvqivqbdqvbyqyqbuyxbqcbqybxacbyqbqybxcyuwvbwvbwibwvbqiiqbvqivqbdqvbyqyqbuyxbqcbqybxacbyqbqybxcyucbqybxacbyqbqybxcyuabcqybvwdbvwqvbiudbvnbyxxbsybxyvbsdybyvv")
	fmt.Println(data, len(data))
	fmt.Println([]byte{})
	enc_result := myrsa.EncryptMsg(data, PUBKEY)
	fmt.Println(enc_result, len(enc_result))
	fmt.Println([]byte{})
	dec_res := myrsa.DecryptMsg(enc_result, PRVKEY)
	fmt.Println(dec_res, len(dec_res))
}
