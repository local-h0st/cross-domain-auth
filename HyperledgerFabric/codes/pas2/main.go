package main

import (
	"encoding/json"
	"fmt"
	"io"
	"msgs"
	"myrsa"
	"net"
)

const (
	pasPrvkey  = "-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQCzpzmS912X0O19pFUmPLOQaLzf4MLv/e60o17LHOjzeFJwR2gP\nnYbVqhJ3bKxf3nO9wvHXO8H/6BaxID+KDddaQVViURcW2KUNNseN52jfliDYExU9\ntnHzPdAmjtZEbt5izmQn1z3jEXFYleZVzpCoRT969RyqXjtFleXdntHOfQIDAQAB\nAoGBAJaotm+5YpPeckvbdE0MuslwDHTzWIdKvNRf7S8In5MOZJQkTfBKerjUV4gv\nap87PnT09zs4sgiZ6e3AzYhI8z8kKhOiwQd95LRLL5xxPd42gixLE2RyBqg5764n\nNkbKMqYtBgNEPKsPWplIzSHRc0gVRD+1BG2janxRaykZWpABAkEA62AsIdRP97Jf\n4jKhjOsGIL5tkLwPl+O6LHDd5N2zg8VmZrmD7QHo+y5w+3ry1QJ2pk5j8xwTi9aw\nlu4UVzpKPQJBAMNlHWHEtECCWtRgkd+ynfc1Y4Q2hwE2+ol5CFCqxO+Ziw0PQBNf\nmqLxG8GXO8e/aWbXDosJvlvuiFhwnrQpGUECQQCjnwlOwv6MG82Xusae5UovPPGB\naZoVZlMnTZaS4KNH+NOEmXXiLi+9XL1htEhWVw4P8fJ9L4lO7oF3ii1NrdGpAkB3\n8MixbBKNircAqOrCWx1WUyJsVSBYIYx2+KGfCsRqo2DUumjFu6jrnn9APXpHqfqk\nUxytQmTkf66YQ0FYK+ZBAkEAz6IXgM7BOO1EaiTi/mckpU28HOHlRAaqgSS/+Q34\nuwHLJM8KXTnImoXuw/wxUu73oME/v+9xaVpJK0ulO8+23g==\n-----END RSA PRIVATE KEY-----\n"
	pasPubkey  = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCzpzmS912X0O19pFUmPLOQaLzf\n4MLv/e60o17LHOjzeFJwR2gPnYbVqhJ3bKxf3nO9wvHXO8H/6BaxID+KDddaQVVi\nURcW2KUNNseN52jfliDYExU9tnHzPdAmjtZEbt5izmQn1z3jEXFYleZVzpCoRT96\n9RyqXjtFleXdntHOfQIDAQAB\n-----END PUBLIC KEY-----\n"
	pasID      = "pasHang"
	domainName = "domainHang"
	pasAddr    = "localhost:51235"
	pasPort    = ":51235"
)

type serverInfo struct {
	NodeID        string
	ServerAddr    string
	ServerPubkey  string
	EnclavePubkey string
}

var serverCache []serverInfo = []serverInfo{serverInfo{
	NodeID:        "VS001",
	ServerAddr:    "localhost:54321",
	EnclavePubkey: "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCTRELjLgbesbJ/ge8qbgb8/m6u\nz5D9MhOUawKwuZFmkruvJKF8xtH0SAHao4Rxni62j2/fZ7Wjvcn9pzRQmMgreVNA\nj2FP67tgSUEGTDEd79J6Cu9mbBDf3u3NMURpR14dKekwx2pRpJuvXtIZdJyIdAnz\nk6moUmNt7TXRbqnmPQIDAQAB\n-----END PUBLIC KEY-----\n",
	ServerPubkey:  "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCU+Q+LZjRCJyWxFMpUeTmdHvzl\nRyNCAsaWxUFI1W9mXywxcItmLmGLEG/tYVzYkjdgjqldvYZfi6q7JNdIYaM7EmXr\n5sgDWbfrH6eeXPe+dTu5dGTSvZ2K3Kvlo17kDKp/H4LcDcG3YpvAPQp/nII9ReWI\n56iq8mohZxb1eMgxPwIDAQAB\n-----END PUBLIC KEY-----\n",
}} //  最开始有先验知识，域管理员在加入时可以得知

func main() {
	ln, err := net.Listen("tcp", pasPort)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			for {
				msg, err := func(conn net.Conn) ([]byte, error) {
					buf := make([]byte, 32768)
					n, err := conn.Read(buf)
					if err == io.EOF {
						return nil, err
					} else if err != nil {
						return nil, err
					}
					msg := buf[:n]
					return msg, nil
				}(conn)
				if err != nil {
					break
				}
				handleMsg(msg)
			}
		}(conn)
	}
}

func handleMsg(cipher []byte) {
	basic_msg := msgs.BasicMsg{}
	if json.Unmarshal(myrsa.DecryptMsg(cipher, []byte(pasPrvkey)), &basic_msg) != nil {
		panic("damn it.(1)")
	}
	index := -1
	for k := range serverCache {
		if serverCache[k].NodeID == basic_msg.SenderID {
			index = k
		}
	}
	if index == -1 {
		panic("damn it.(2)")
	}
	if !basic_msg.VerifySign([]byte(serverCache[index].ServerPubkey)) {
		panic("damn it.(3)")
	}
	switch basic_msg.Method {
	case "requireSyncBlacklist":
		// TODO .

	default:
		fmt.Println(basic_msg.Method, ": unknown method.")
	}
}
