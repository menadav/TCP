package logger

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"
)

var out = log.New(os.Stdout, "", 0)

func emit(level, msg string, args ...any) {
	entry := map[string]any{
		"time":  time.Now().UTC().Format(time.RFC3339Nano),
		"level": level,
		"msg":   msg,
	}
	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		entry[key] = args[i+1]
	}
	if len(args)%2 == 1 {
		entry["arg"] = args[len(args)-1]
	}
	b, err := json.Marshal(entry)
	if err != nil {
		out.Printf(`{"level":%q,"msg":%q,"log_error":%q}`, level, msg, err.Error())
		return
	}
	out.Println(string(b))
}

func Info(msg string, args ...any) {
	emit("INFO", msg, args...)
}

func Warn(msg string, args ...any) {
	emit("WARN", msg, args...)
}

func Error(msg string, args ...any) {
	emit("ERROR", msg, args...)
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
