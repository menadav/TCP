package main

import (
	"fnt"
	"net"
)

func main(){
	var listen net.Conn
	var err error
	var n int

	conn, err = net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error to connect", err)
		return
	}
	defer conn.Close()

	go readServer(conn)
}
