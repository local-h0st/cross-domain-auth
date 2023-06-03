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
	m.Signature = myrsa.RsaSignWithSha256((jsonbyte), prvkey)
}

func (m BasicMsg) VerifySign(pubkey []byte) bool {
	sign_given := m.Signature
	m.Signature = nil
	jsonbyte, _ := json.Marshal(m)
	return myrsa.RsaVerySignWithSha256(jsonbyte, sign_given, pubkey)
}

type VerifyMsg struct {
	UpdateFlag    bool
	CipherID      []byte
	Domain        string
	ServerID      string // 来自哪个server
	ServerAddr    string
	DomainPasID   string
	DomainPasAddr string // ip:port形式，例如localhost:6666

}

type AddServerPubkeyMsg struct {
	ServerID     string
	ServerPubkey []byte
}

type UpdateServerPubkeyMsg struct {
	ServerID        string
	ServerNewPubkey []byte
	Signature       []byte
}

func (m *UpdateServerPubkeyMsg) GenSign(prvkey []byte) {
	m.Signature = nil
	jsonbyte, _ := json.Marshal(m)
	m.Signature = myrsa.RsaSignWithSha256((jsonbyte), prvkey)
}

func (m UpdateServerPubkeyMsg) VerifySign(pubkey []byte) bool {
	sign_given := m.Signature
	m.Signature = nil
	jsonbyte, _ := json.Marshal(m)
	return myrsa.RsaVerySignWithSha256(jsonbyte, sign_given, pubkey)
}

type NeedPubkey struct {
	SenderID     string
	SenderPubkey []byte
	SenderAddr   string
}

// 结构内可以加入随机数防止截获密文重放攻击
