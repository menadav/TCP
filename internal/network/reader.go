package network

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"io"
)

type TextProcessor func(string)

func TextClient(textclient string){
	fmt.Println("\n[Client]:", textclient)
}

func TextServer(textserver string){
	fmt.Println("\n[Server]:", textserver)
}

func StartScanner(source io.Reader, process TextProcessor) {
	var scanner 	*bufio.Scanner
	var err		error

	scanner = bufio.NewScanner(source)
	for scanner.Scan() {
		process(scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("ERROR scanner:", err)
	}
}

func ReadServer(conn net.Conn, process TextProcessor) {
	StartScanner(conn, process)
}

func WriteFromStdin(conn net.Conn) {
	var sendCallback 	func(string)
	var err 		error

	sendCallback = func(text string){
		_, err = conn.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("Error al enviar datos (conexión cerrada):", err)
		}
	}
	StartScanner(os.Stdin, sendCallback)
}

func ClientAtender(c net.Conn) {
	defer c.Close()
	go ReadServer(c, TextClient)
	WriteFromStdin(c)
}
