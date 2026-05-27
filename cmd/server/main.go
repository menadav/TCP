package main

import (
	"answer_protocol/src/network"
	"answer_protocol/src/speakserver"
	"answer_protocol/src/constructor"
	"fmt"
	"net"
)

func main(){
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listen", err)
		return
	}
	defer listen.Close()
	hub := constructor.NewHub()
	go hub.Run()
	fmt.Println("Server ready on the port :")
	for {
		conn, err := listen.Accept()
		if err != nil{
			fmt.Println("Error Accept", err)
			continue
		}
		fmt.Println("Client connected from:", conn.RemoteAddr())
		speak.SendSuccess(conn, "hello proto=1")
		go network.ClientAtender(conn, hub)
	}
}
