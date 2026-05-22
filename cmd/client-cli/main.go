package main

import (
	"answer_protocol/internal/network"
	"fmt"
	"net"
)

func main(){
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	go network.ReadServer(conn, network.TextClient)
	network.WriteFromStdin(conn)
}
