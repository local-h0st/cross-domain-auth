package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
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

var node_id string = os.Getenv("SERVERID")
var PRVKEY, PUBKEY []byte
var serving_port string = ":54321"
var db_path string = "./db"

func main() {
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

	PRVKEY, PUBKEY = myrsa.GenRsaKey()
	fmt.Println("PUBKEY ==> ", string(bytes.Replace(node_pk, []byte("\n"), []byte("\\n"), -1)))

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
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err == io.EOF {
		conn.Close()
		return nil, err
	} else if err != nil {
		conn.Close()
		return nil, err
	}
	msg := buf[:n]
	return msg, nil
}

func handleMsg(cipher []byte, contract *client.Contract, db *leveldb.DB) {
	basic_msg := basicMsg{}
	if json.Unmarshal([]byte(decryptMsg(cipher)), &basic_msg) != nil {
		fmt.Println("[handleMsg] basic msg json unmarshal failed.")
		return
	}
	switch basic_msg.Method {
	case "verifyResult":
		// HERE!!从上到下重构
	case "updateServerInfo":
		updateServerInfo(basic_msg.CipherText)
	default:
		return
	}
	// 事实上收到的是加密后的密文，需要用node_sk解密
	// decryptMsg暂时没写
	msg_text := decryptMsg(cipher)
	// 解析msg内容
	type message struct {
		PID              string
		P                string // // 门限签名技术的那个点，我也不知道用什么格式存储
		Tag              string // node_id
		ID_cipher        string
		PK_device2domain string // 不知道公钥是这哪个类型。另外是建议存到数据库，还是存到内存算了？
	}
	var msg message
	err := json.Unmarshal([]byte(msg_text), &msg)
	if err != nil {
		fmt.Println("[handleMsgForPseudo] json unmarshal failed.")
		return err
	}
	// 查询pid是否存在
	if isExist(contract, msg.PID) {
		fmt.Println("[handleMsgForPseudo] pid exists, en error was returned.")
		return fmt.Errorf("[err from HandleMsgForPseudo] pid exists.")
	}
	// 用于存储pid和点p，以json的格式
	// 这里我选择LevelDB - Fabric自带的内嵌键值存储数据库,可以直接在chaincode中使用。
	// 数据只保存在一个peer节点上,不可持久化。
	// https://www.tizi365.com/archives/411.html
	err = db.Put([]byte(msg.PID), []byte(msg_text), nil)
	if err != nil {
		return fmt.Errorf("[err from HandleMsgForPseudo] put record to db failed.")
	}
	/*
		val, err := db.Get([]byte(msg.PID), nil)
		if err != nil {
			return fmt.Errorf("[err from HandleMsgForPseudo] read record from db failed.")
		} else {
			fmt.Println("get from db: ", string(val))
		}*/
	sgxVerifyID(msg.ID_cipher)
	return nil
}

type basicMsg struct {
	Method     string
	CipherText string
}

func sgxVerifyID(cipher_id string) bool {
	// 在enclave内部解密核验是否在黑名单上，true表示合法
	// 出错了一律直接false，这里不再返回err
	// TODO

	return true
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
