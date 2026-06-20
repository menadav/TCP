package speak

import (
	"fmt"
	"net"
)

func SendSuccess(conn net.Conn, dade string) {
	answer := fmt.Sprintf("OK %s\n", dade)
	conn.Write([]byte(answer))
}

func SendEvent(conn net.Conn, category string, dade string) {
	answer := fmt.Sprintf("EVT %s %s\n", category, dade)
	conn.Write([]byte(answer))
}
