package main

import (
	"bytes"
	"fmt"
	"myrsa"
)

func main() {
	PRVKEY, PUBKEY := myrsa.GenRsaKey()
	fmt.Println("PRVKEY ==> ", string(bytes.Replace(PRVKEY, []byte("\n"), []byte("\\n"), -1)))
	fmt.Println("PUBKEY ==> ", string(bytes.Replace(PUBKEY, []byte("\n"), []byte("\\n"), -1)))
}
