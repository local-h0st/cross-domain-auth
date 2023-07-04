package msgs

import (
	"encoding/json"
	"myrsa"
)

// TODO 用于VS上账本公布自己的信息
type ServerRecord struct {
	NodeID        string
	ServerAddr    string
	ServerPubkey  string
	EnclavePubkey string
}

//

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
	DestDomain string
	PID        string
	CipherID   []byte
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

type VerifyResultMsg struct {
	PID        string
	DestDomain string
	Result     string
}

type DomainRecord struct {
	Domain                       string
	PasID                        string
	PasAddr                      string
	PasPubkey                    []byte
	BlacklistLastUpdateTimestamp string
	WaitQ                        []FragmentMsg // 存储正在核验但是没有写入账本的Fragment记录
}

type BlacklistRecord struct {
	Domain    string
	BlackList []string
}

type UpdateBlacklistTimestampMsg struct {
	Domain    string
	Timestamp string
}

// For serverVS
type FragmentMsg struct {
	Tag              string // node_id
	PID              string
	OrigDomain       string
	DestDomain       string
	CipherID         []byte
	PubkeyDev2Domain []byte
	Point            []byte // 门限签名技术的那个点，我也不知道用什么格式存储，此信息不应该写入账本

}
