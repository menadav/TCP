package network

import (
	"answer_protocol/src/logger"
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
		conn.SetReadDeadline(time.Now().Add(50 * time.Second))
		if !scanner.Scan() {
			break
		}
		process(scanner.Text())
	}
	err := scanner.Err()
	if err := scanner.Err(); err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return
		}

		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			speak.SendErr(conn, speak.ErrTimeout)
			logger.Warn("connection timeout", "addr", logger.Addr(conn))
			return
		}
		logger.Error("scanner error", "addr", logger.Addr(conn), "err", err)
	}
}

func WriteFromStdin(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text() + "\n"))
	}
	fmt.Println("ERR 106 CONTROL_D")
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
