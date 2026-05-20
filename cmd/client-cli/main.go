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

	network.ReadServer(conn, func(textCli string){
		fmt.Println("Server:", textCli)
	})
}
