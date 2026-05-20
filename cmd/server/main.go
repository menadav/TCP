package main

import (
	"fmt"
	"net"
)

func main(){
	var listen net.Listener
	var err error
	var user net.Conn

	listen, err = net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listen", err)
		return
	}
	defer listen.Close()
	for {
		user, err = listen.Accept()
		if err != nil{
			fmt.Println("Error Accept", err)
			continue
		}
		fmt.Println("Cliente conectado desde:", user.RemoteAddr())
		user.Write([]byte("Hola, bienvenido a The Answer Protocol\n"))
	}
}
