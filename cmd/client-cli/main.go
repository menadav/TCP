package main

import (
	"answer_protocol/internal/network"
	"bufio"
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
	scanner := bufio.NewScanner(conn)
	go network.ReadServer(scanner, network.TextClient)
	network.WriteFromStdin(conn)
}
