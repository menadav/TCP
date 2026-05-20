package network

import (
	"net"
	"fmt"
	"bufio"
)

func ReadServer(conn net.Conn, process func(string)){
	var scanner 	*bufio.Scanner
	var err		error
	var text	string

	scanner = bufio.NewScanner(conn)
	for scanner.Scan() {
		text = scanner.Text()
		process(text)
	}

	err = scanner.Err()
	if err != nil {
		fmt.Println("ERROR scanner", err)
		return
	}
}
