package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"msgs"
	"myrsa"
	"net"
	"os"
	"path"
	sharedconfigs "sharedConfigs"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID = "Org1MSP"
	// cryptoPath   = "../../test-network/organizations/peerOrganizations/org1.example.com"
	cryptoPath = "../../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	// certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

var PRVKEY, PUBKEY, enPubkey []byte
var domains []msgs.DomainRecord

func main() {
	PRVKEY = []byte(sharedconfigs.ServerPrvkey)
	PUBKEY = []byte(sharedconfigs.ServerPubkey)
	enPubkey = []byte(sharedconfigs.EnclavePubkey)
	// 初始化leveldb
	db, err := leveldb.OpenFile(sharedconfigs.DatabasePath, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 连接gateway，照抄application-gateway-go
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()
	id := newIdentity()
	sign := newSign()
	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()
	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}
	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)
	initLedger(contract)

	// TODO for test
	testmsg := msgs.BasicMsg{
		Method: "testingConnection",
	}
	testmsg.GenSign(PRVKEY)
	teststr, _ := json.Marshal(testmsg)
	sendMsg(sharedconfigs.EnclaveAddr, string(myrsa.EncryptMsg(teststr, enPubkey)))
	//

	fmt.Println("[main] listening on", sharedconfigs.ServerPort)
	ln, err := net.Listen("tcp", sharedconfigs.ServerPort)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn, contract, db)
	}
}

func handleConn(conn net.Conn, contract *client.Contract, db *leveldb.DB) {
	defer conn.Close()
	for {
		msg, err := parseMessage(conn)
		if err != nil {
			break
		}
		handleMsg(msg, contract, db)
	}
}

func parseMessage(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 32768)
	n, err := conn.Read(buf)
	if err == io.EOF {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func handleMsg(cipher []byte, contract *client.Contract, db *leveldb.DB) {
	basic_msg := msgs.BasicMsg{}
	jsonstr := decryptMsg(cipher)
	if json.Unmarshal(jsonstr, &basic_msg) != nil {
		fmt.Println("[handleMsg] basic msg json unmarshal failed:", string(jsonstr))
		return
	}
	// TODO 有几个case没有核验签名，需要确认确实是PAS发过来的，PAS初始pubkey信息可以由管理员配置，所以一定可信。可以写个程序一键配置serverVS、sgxInteract等等的初始pubkey
	switch basic_msg.Method {
	// admin sig
	case "syncDomain":
		// TODO sender should be admin, not verified SenderID yet.
		if !basic_msg.VerifySign([]byte(sharedconfigs.AdminPubkey)) {
			fmt.Println("syncDomain: admin sig invalid, reject.")
			return
		}
		syncDomain(basic_msg.Content)
	case "getFragment":
		if !basic_msg.VerifySign([]byte(sharedconfigs.AdminPubkey)) {
			fmt.Println("getFragment: admin sig invalid, reject.")
			return
		}
		data, err := db.Get(basic_msg.Content, nil)
		if err != nil {
			fmt.Println("getFragment: level db get fragment failed.")
		}
		fmt.Println("fragment ==> ", string(data))
	// pas sig
	case "fragment":
		var domain_pubkey []byte
		for _, domain := range domains {
			if basic_msg.SenderID == domain.PasID {
				domain_pubkey = domain.PasPubkey
				break
			}
		}
		if domain_pubkey == nil {
			fmt.Println("fragment: failed to find PAS pubkey, may lack of domain info.")
			return
		} else {
			if !basic_msg.VerifySign(domain_pubkey) {
				fmt.Println("fragmet: pas sig invalid, reject.")
				return
			}
		}
		fragment(basic_msg.Content, contract, db)
	// enclave sig
	case "verifyResult":
		if !basic_msg.VerifySign([]byte(sharedconfigs.EnclavePubkey)) {
			fmt.Println("verifyResult: enclave sig invalid, reject.")
			return
		}
		result_msg := msgs.VerifyResultMsg{}
		if json.Unmarshal(basic_msg.Content, &result_msg) != nil {
			fmt.Println("verifyResult: json unmarshal failed.")
			return
		}
		d_ind, q_ind := -1, -1
		for k := range domains {
			if domains[k].Domain == result_msg.DestDomain {
				for j := range domains[k].WaitQ {
					if domains[k].WaitQ[j].PID == result_msg.PID {
						d_ind = k
						q_ind = j
						break
					}
				}
				break
			}
		}
		if addPseudoRecordToLedger(contract, result_msg.Result, domains[d_ind].WaitQ[q_ind]) != nil {
			fmt.Println("verifyResult: failed to put result to ledger, ")
		}
		domains[d_ind].WaitQ[q_ind] = msgs.FragmentMsg{} // 清除记录
	// no sig required
	case "queryLedger":
		queryLedger(contract)
	default:
		fmt.Println(basic_msg.Method, ": unknown method.")
		return
	}
}

func syncDomain(jsonmsg []byte) {
	sync_msg := msgs.DomainRecord{}
	if json.Unmarshal(jsonmsg, &sync_msg) != nil {
		fmt.Println("syncDomain: json unmarshal failed.")
		return
	}
	index := -1
	for k := range domains {
		if domains[k].Domain == sync_msg.Domain {
			index = k
			break
		}
	}
	if index == -1 {
		domains = append(domains, sync_msg)
	} else {
		domains[index] = sync_msg
	}
}
func fragment(jsonmsg []byte, contract *client.Contract, db *leveldb.DB) {
	fragment_msg := msgs.FragmentMsg{}
	if json.Unmarshal(jsonmsg, &fragment_msg) != nil {
		fmt.Println("[fragment] json unmarshal failed.")
		return
	}
	// 查询pid是否存在
	if isExist(contract, fragment_msg.PID) {
		fmt.Println("[fragment] pid already exists.")
		return
	}
	// 存储信息
	if db.Put([]byte(fragment_msg.PID), jsonmsg, nil) != nil {
		fmt.Printf("[fragment] put record to db failed.")
		return
	}
	if fragment_msg.Tag == sharedconfigs.NodeID {
		basic_msg := msgs.BasicMsg{
			Method:    "verifyID",
			SenderID:  sharedconfigs.NodeID,
			Content:   nil,
			Signature: nil,
		}
		basic_msg.Content, _ = json.Marshal(msgs.VerifyMsg{
			PID:        fragment_msg.PID,
			CipherID:   fragment_msg.CipherID,
			DestDomain: fragment_msg.DestDomain,
		})
		basic_msg.GenSign(PRVKEY)
		data, _ := json.Marshal(basic_msg)
		cipher := myrsa.EncryptMsg(data, []byte(sharedconfigs.EnclavePubkey))
		sendMsg(sharedconfigs.EnclaveAddr, string(cipher))
		for k := range domains {
			if domains[k].Domain == fragment_msg.DestDomain {
				domains[k].WaitQ = append(domains[k].WaitQ, fragment_msg)
				break
			}
		}
	}
}
func addPseudoRecordToLedger(contract *client.Contract, valid string, rec msgs.FragmentMsg) error {
	var pseudo_record struct {
		PID              string
		OrigDomain       string
		DestDomain       string
		PubkeyDev2Domain string
		Valid            string
		Tag              string
		Timestamp        string
	} // 需要和demo的数据结构手工同步
	pseudo_record.PID = rec.PID
	pseudo_record.OrigDomain = rec.OrigDomain
	pseudo_record.DestDomain = rec.DestDomain
	pseudo_record.PubkeyDev2Domain = string(rec.PubkeyDev2Domain)
	pseudo_record.Valid = valid
	pseudo_record.Tag = rec.Tag
	pseudo_record.Timestamp = ""

	jsonstr, _ := json.Marshal(pseudo_record)

	_, err := contract.SubmitTransaction("AddPseudoRecord", string(jsonstr))
	if err != nil {
		return err
	}
	return nil
}
func queryLedger(contract *client.Contract) {
	evaluateResult, err := contract.EvaluateTransaction("QueryAll")
	if err != nil {
		fmt.Println("[queryLedger] transcation evaluate failed.")
	}
	fmt.Println("[queryLedger] result ==>", formatJSON(evaluateResult))
}
func isExist(contract *client.Contract, pid string) bool {
	evaluateResult, err := contract.EvaluateTransaction("CheckExistance", pid)
	if err != nil {
		panic(fmt.Errorf("[isExist] failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)
	if result == "true" {
		return true
	} else {
		return false
	}
}

func decryptMsg(cipher []byte) []byte {
	// 肯定是拿自己的prvkey解密
	return myrsa.DecryptMsg(cipher, PRVKEY)
}
func encryptMsg(text []byte, pubkey []byte) []byte {
	return myrsa.EncryptMsg(text, pubkey)
}
func sendMsg(addr string, data string) {
	// 连接到指定IP和端口
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}
	// 发送消息
	n, err := fmt.Fprint(conn, data)
	if err != nil {
		fmt.Printf("send to %s failed.", addr)
	} else {
		fmt.Println("msg sent, total", n, "bytes.")
	}
}

/*----------- 直接参考application-gateway-go -----------*/
func initLedger(contract *client.Contract) {
	fmt.Printf("[InitLedger] func called.\n")

	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("[InitLedger] executed successfully.\n")
}
func newGrpcConnection() *grpc.ClientConn {
	// newGrpcConnection creates a gRPC connection to the Gateway server.
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}
func newIdentity() *identity.X509Identity {
	// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}
func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}
func newSign() identity.Sign {
	// newSign creates a function that generates a digital signature from a message digest using a private key.
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}
func formatJSON(data []byte) string {
	// Format JSON data
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

// func sendAddServerPubkeyMsg() {
// 	fmt.Println("[main] sendAddServerPubkeyMsg()")
// 	basic_msg := msgs.BasicMsg{
// 		Method:    "addServerPubkey",
// 		SenderID:  sharedconfigs.NodeID,
// 		Content:   nil,
// 		Signature: nil,
// 	}
// 	basic_msg.Content, _ = json.Marshal(msgs.AddServerPubkeyMsg{
// 		ServerPubkey: PUBKEY,
// 	})
// 	data, _ := json.Marshal(basic_msg)
// 	cipher := myrsa.EncryptMsg(data, []byte(sharedconfigs.EnclavePubkey))
// 	sendMsg(sharedconfigs.EnclaveAddr, string(cipher))
// }
