package msgs

import (
	"encoding/json"
	"myrsa"
)

type BasicMsg struct {
	Method    string
	SenderID  string
	Content   []byte
	Signature []byte
}

func (m *BasicMsg) GenSign(prvkey []byte) {
	m.Signature = nil
	jsonbyte, _ := json.Marshal(m)
	m.Signature = myrsa.SignMsg((jsonbyte), prvkey)
}

func (m BasicMsg) VerifySign(pubkey []byte) bool {
	sign_given := m.Signature
	m.Signature = nil
	jsonbyte, _ := json.Marshal(m)
	return myrsa.VerifyMsgSig(jsonbyte, sign_given, pubkey)
}

type VerifyMsg struct {
	CipherID   []byte
	Domain     string
	SenderAddr string
	// 同步blacklist用
	UpdateFlag    bool
	DomainPasAddr string // ip:port形式，例如localhost:6666
	// DomainPasID   string // 多余了，能根据Domain查出来
	// ServerID      string // 来自哪个server
	// ServerAddr    string
}

type AddServerPubkeyMsg struct {
	// ServerID     string
	ServerPubkey []byte
}

type UpdateServerPubkeyMsg struct {
	// ServerID        string
	ServerNewPubkey []byte
	// Signature       []byte
	// 感觉senderID和serverID功能重合了，可以改一下
}

// func (m *UpdateServerPubkeyMsg) GenSign(prvkey []byte) {
// 	m.Signature = nil
// 	jsonbyte, _ := json.Marshal(m)
// 	m.Signature = myrsa.SignMsg((jsonbyte), prvkey)
// }

// func (m UpdateServerPubkeyMsg) VerifySign(pubkey []byte) bool {
// 	sign_given := m.Signature
// 	m.Signature = nil
// 	jsonbyte, _ := json.Marshal(m)
// 	return myrsa.VerifyMsgSig(jsonbyte, sign_given, pubkey)
// }
/*
type NeedPubkey struct {
	// 小丑结构
	SenderID     string
	SenderPubkey []byte
	SenderAddr   string
}
*/
type DomainInfoRecord struct {
	Domain    string
	PasID     string
	BlackList []string
}

// 结构内可以加入随机数防止截获密文重放攻击
