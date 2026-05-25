package main

import (
	"answer_protocol/internal/network"
	"answer_protocol/internal/constructor"
	"fmt"
	"net"
)



func main(){
	listen, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println("Error listen", err)
		return
	}
	defer listen.Close()
	hub := constructor.NewHub()
	go hub.Run()
	fmt.Println("Servidor listo en el puerto :8080")
	for {
		conn, err := listen.Accept()
		if err != nil{
			fmt.Println("Error Accept", err)
			continue
		}
		fmt.Println("Client connected from:", conn.RemoteAddr())
		network.SendSuccess(conn, "hello proto=1")
		go network.ClientAtender(conn, hub)
	}
}
