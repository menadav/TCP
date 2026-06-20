package main

import (
	"answer_protocol/src/network"
	"bufio"
	"fmt"
	"net"
	"os"
)

func main(){
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	go func() {
		network.StartScanner(scanner, network.TextClient, conn)
		os.Exit(0)
	}()
	network.WriteFromStdin(conn)
}