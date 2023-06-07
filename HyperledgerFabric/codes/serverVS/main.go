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

const (
	db_path       = "./db"
	serving_port  = ":54321"
	selfAddr      = "localhost:54321"
	enclaveAddr   = "localhost:55555"
	enclavePubkey = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCj8hW+keEOHHHLV/7BRO7I0j7a\nXAfxTvkiM8Qyex+aMQ7Ny+cavF4mWlJmdmGo9K3jHFH3LyEd2JuGPh5T0ad/O76C\nor+hX+RvgXkg0HS3MEQIwmzmNg57RSaNxzlJatXEfjpRJJ5Nc+dyA6hpYzaNj9LY\nKex5gvsGFpBMwQZyVwIDAQAB\n-----END PUBLIC KEY-----\n"
)

// var node_id string = os.Getenv("SERVERID")
var node_id string = "serverVS001"
var PRVKEY, PUBKEY []byte

func sendAddServerPubkeyMsg() {
	fmt.Println("[main] sendAddServerPubkeyMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "addServerPubkey",
		SenderID:  node_id,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.AddServerPubkeyMsg{
		ServerPubkey: PUBKEY,
	})
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(enclavePubkey))
	sendMsg(enclaveAddr, string(cipher))
}

func main() {
	PRVKEY, PUBKEY = myrsa.GenRsaKey()
	fmt.Println("PUBKEY ==> ", string(bytes.Replace(PUBKEY, []byte("\n"), []byte("\\n"), -1)))

	// 方便测试就直接指定了
	PRVKEY = []byte("-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQDQlXmFEiNzbO0iHjdYIUPvbWmqPmMJcrGVLjRUrr2HtURh9lcr\nGsti1r4BesFcuS+QAzBlFZsp50Ytae0snr26jnFLOpBGscCDLsyPrL3dlUnWGnQY\n5SOFjvVpAjsuc16W0TXXzdoaW6yZwX+tKd2yLkgbcL0alZeTI1v8lJN9YQIDAQAB\nAoGBAL7k/fk+l3lM6F3AP7CFiVI31Wu8exEricDZL4WNAuKPkA0D0dUeSaOkmvJp\nsUu2JARuFr18n6wjAMQRXMHoagQKt4zKB4Kcv6vNKNC+DLQVnd9WvujwDYChlty7\nOd+vuK5tXBIUi5FwF/QjHhOIj8EKhb228sxqXshEViuHO3dtAkEA4O6YM2cZp2Za\ngJODiRcXOjXoSkSVMqYfhARGiWRvJ0DuGAkksBAzTSxPeu0N6yEPv2Cddw2HCa1+\n24+Mqm6PPwJBAO1k0wYePhClgbTh/6sUCx5lyK+l+oTJ7H5heb0LthNX7n/B1xKv\nNr2JZtRRQ+QSRK+oTztVDSd87C0jRgITa18CQHdLU4F/nsV/rWQf2FUu3+zJhmdN\nNGvmWzSjJ93aXHFPKHeq8cBG905ov8aMTyNzJ2zyitEHZaUmVO+RlKMXe/UCQQCR\nSUxxCR849uHr/wiG/kxTvT1WaoFotV/cdPGZhkpXilA3tj1XfQ5Gb4oUVOv08E1D\nKAHdsQ7M5QJyGY1mBdaHAkAf5Y1dfFOOTQqZNkSWgm0lFfX1tfUPMD1/XnfXBgxr\n4VvDGz66MsAvh1v4qYw8GcKoGxetTWS8yzIYiBMgqD5l\n-----END RSA PRIVATE KEY-----\n")
	PUBKEY = []byte("-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDQlXmFEiNzbO0iHjdYIUPvbWmq\nPmMJcrGVLjRUrr2HtURh9lcrGsti1r4BesFcuS+QAzBlFZsp50Ytae0snr26jnFL\nOpBGscCDLsyPrL3dlUnWGnQY5SOFjvVpAjsuc16W0TXXzdoaW6yZwX+tKd2yLkgb\ncL0alZeTI1v8lJN9YQIDAQAB\n-----END PUBLIC KEY-----\n")
	sendAddServerPubkeyMsg()
	// 初始化leveldb
	db, err := leveldb.OpenFile(db_path, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 连接gateway，照抄application-gateway-go
	fmt.Println("[main] connecting to gateway...")
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

	fmt.Println("[main] listening on", serving_port)
	ln, err := net.Listen("tcp", serving_port)
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
	// TODO 没有核验签名，需要确认确实是PAS发过来的，PAS初始pubkey信息可以由管理员配置，所以一定可信。可以写个程序一键配置serverVS、sgxInteract等等的初始pubkey
	switch basic_msg.Method {
	case "fragment":
		fragment(basic_msg.Content, contract, db)
	case "verifyResult":
		result_msg := msgs.VerifyResultMsg{}
		if json.Unmarshal(basic_msg.Content, &result_msg) != nil {
			fmt.Println("verifyResult json unmarshal failed.")
			return
		}
		addPseudoRecordToLedger(contract, result_msg.PID, result_msg.Result, string(result_msg.PubkeyDeviceToDomain))
	case "queryLedger":
		queryLedger(contract)
	default:
		return
	}

}

func fragment(jsonmsg []byte, contract *client.Contract, db *leveldb.DB) {
	fmt.Println("[fragment] exec...")
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
	if fragment_msg.Tag == node_id {
		basic_msg := msgs.BasicMsg{
			Method:    "verifyID",
			SenderID:  node_id,
			Content:   nil,
			Signature: nil,
		}
		basic_msg.Content, _ = json.Marshal(msgs.VerifyMsg{
			PID:                  fragment_msg.PID,
			PubkeyDeviceToDomain: fragment_msg.PubkeyDeviceToDomain,
			CipherID:             myrsa.EncryptMsg(fragment_msg.CipherID, []byte(enclavePubkey)),
			Domain:               "domainA",
			SenderAddr:           selfAddr,
			UpdateFlag:           false,
			DomainPasAddr:        "",
		})
		basic_msg.GenSign(PRVKEY)
		data, _ := json.Marshal(basic_msg)
		cipher := myrsa.EncryptMsg(data, []byte(enclavePubkey))
		sendMsg(enclaveAddr, string(cipher))
	}
}
func addPseudoRecordToLedger(contract *client.Contract, pid, valid, pubkey_device_to_domain string) error {
	_, err := contract.EvaluateTransaction("AddPseudoRecord", pid, valid, pubkey_device_to_domain)
	if err != nil {
		return err
	}
	// fmt.Println("[addPseudoRecordToLedger] result ==>", evaluateResult)
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

/*----------- 直接抄application-gateway-go -----------*/
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
