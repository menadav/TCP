package main

import (
	"answer_protocol/internal/network"
	"fmt"
	"net"
)

func main(){
	var conn net.Conn
	var err error

	conn, err = net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Error to connect client CLI:", err)
		return
	}
	defer conn.Close()

	go network.ReadServer(conn, network.TextServer)
	network.WriteFromStdin(conn)
}
