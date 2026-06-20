package network

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type TextProcessor func(string)

func TextClient(text string) {
	fmt.Println(text)
}

func StartScanner(scanner *bufio.Scanner, process func(string), conn net.Conn) {
	for {
		conn.SetReadDeadline(time.Now().Add(3 * time.Minute))
		if !scanner.Scan() {
			break
		}
		process(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return
		}

		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			speak.SendErr(conn, speak.ErrTimeout)
			fmt.Println("Timeout de inactividad, close conexión:", conn.RemoteAddr())
			return
		}
		fmt.Println("ERROR scanner:", err)
	}
}

func WriteFromStdin(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text() + "\n"))
	}
}

func BroadcastMessage(name string, hub *models.Hub) TextProcessor {
	return func(msg string) {
		new_name := fmt.Sprintf("[%s]", name)
		hub.Broadcast <- models.Message{
			Scope:   models.ScopeGlobal,
			Content: fmt.Sprintf("EVT GLOBAL CHAT %s %s\n", new_name, msg),
		}
	}
}

func ReadLine(conn net.Conn) string {
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	return scanner.Text()
}
