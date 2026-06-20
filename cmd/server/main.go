package main

import (
	"answer_protocol/src/constructor"
	"answer_protocol/src/network"
	"answer_protocol/src/speakserver"
	"answer_protocol/src/world"
	"fmt"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listen", err)
		return
	}
	defer listen.Close()
	data, err := world.LoadWorld("data.yaml")
	if err != nil {
		fmt.Println("Error load world", err)
		return
	}
	hub := constructor.NewHub(data)
	go hub.Run()
	fmt.Println("Server ready on the port :")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error Accept", err)
			continue
		}
		fmt.Println("Client connected from:", conn.RemoteAddr())
		go network.ClientAtender(conn, hub)
		speak.SendSuccess(conn, "hello proto=1")
	}
}
