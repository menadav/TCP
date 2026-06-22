package logger

import (
	"log/slog"
	"net"
	"os"
)

var log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

func Addr(conn net.Conn) string {
	if conn == nil {
		return ""
	}
	if a := conn.RemoteAddr(); a != nil {
		return a.String()
	}
	return ""
}
