package main

import (
	"encoding/json"
	"fmt"
	"msgs"
	"net"
)

const AdminPrvkey = "-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQDVDTjZaYpCk5R8RNcq86rfrbxn6pss6d+MlodAu/UdgfstJzQK\n6X9gqWceCyAmm30CHdRZbC9E39E2U13uF9NX2oTcpdC0XLRAW+rskd/jCGB0zQSf\ne4ktzKu0a5Cp1pwLWyZtCU7re6A2bhVeOh9JUNADFBiemhYl3JIfetAVDwIDAQAB\nAoGBAKjAjlT3GcJeLvC3fk7RLnl5nZAZ7cuHe8BZwsvtlNtIh3FeagRyqqgfxkOv\nwEmUQ1IX2ojx/gbp2UbUhcP/LzAmKN0CWFlpOF+XodndhPVSMYMOajua6+IWXWIO\nKk70FgFNSTKZ2x+dua6P7UgDtXCBs8EZo/RtTF+ieQ7tBSWxAkEA++TqC+9TWZWo\n6ThnfY+VJF4/xrfZL3xbogZUhMutkpJokQSVZTRsciXLhMvsrq1tFqspRUDE9BlQ\nXjFstizu5QJBANiGOnz5dM9PHVnk4ko41eOUlMHm5bIll6geYms2CDzfSUZCJ6Lu\nXusNuTUduuH8xdoJD+knLgWoLuPHMy1TQOMCQGSpkFaAp6BvTHcXEVR+Iq3L9FSn\nd+WgHsZbHT+MXarrU1pQqJsvHf9n1zMUg1sy9xtN/0orngmmbBWYTsdmoXkCQQDU\n9njScNzmBg99WjUEAZDGLV5+tKaZKIZYkcIFZviFPqyoUOsBQujS0gWm653jJiZH\nhIBEtwd6AuhTmpqIawk3AkEAnpe91sUA0SouCWIzgI2EQu4Z6pEgpqb6AHBuuTfF\nvYkZsveKy8edEi2hGzo+U6nb6S4VYnm/rOQKum1Mx5GxEQ==\n-----END RSA PRIVATE KEY-----\n"
const TargetServerAddr = "localhost:54321"

func main() {

}

func syncDomain() {
	basic_msg := msgs.BasicMsg{
		Method:    "syncDomain",
		SenderID:  "admin",
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.DomainRecord{
		Domain:    "domainCalifornia",
		PasID:     "none",
		PasAddr:   "none",
		PasPubkey: nil,
		WaitQ:     nil,
	})
	basic_msg.GenSign([]byte(AdminPrvkey))
	jsonstr, _ := json.Marshal(basic_msg)
	sendMsg(TargetServerAddr, string(jsonstr))
}

func sendMsg(addr string, data string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}
	n, err := fmt.Fprint(conn, data)
	if err != nil {
		fmt.Printf("send to %s failed.", addr)
	} else {
		fmt.Println(n, "bytes sent.")
	}
}
