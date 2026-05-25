package network

import (
	"fmt"
	"net"
)

func SendSuccess(conn net.Conn, dade string) {
	answer := fmt.Sprintf("OK %s\n", dade)
	conn.Write([]byte(answer))
}

func SendError(conn net.Conn, code int, dade string) {
	answer := fmt.Sprintf("ERR %d %s\n", code, dade)
	conn.Write([]byte(answer))
}

func SendEvent(conn net.Conn, category string, dade string) {
	answer := fmt.Sprintf("EVT %s %s\n", category, dade)
	conn.Write([]byte(answer))
}
