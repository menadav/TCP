package speak

import (
	"answer_protocol/src/logger"
	"fmt"
	"net"
)

func SendSuccess(conn net.Conn, dade string) {
	answer := fmt.Sprintf("OK %s\n", dade)
	conn.Write([]byte(answer))
	logger.Info("response sent", "addr", logger.Addr(conn), "kind", "OK", "data", dade)
}

func SendEvent(conn net.Conn, category string, dade string) {
	answer := fmt.Sprintf("EVT %s %s\n", category, dade)
	conn.Write([]byte(answer))
	logger.Info("response sent", "addr", logger.Addr(conn), "kind", "EVT", "category", category, "data", dade)
}
