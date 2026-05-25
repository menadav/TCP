package network

import (
    "answer_protocol/internal/models"
    "answer_protocol/internal/constructor"
    "bufio"
    "net"
    "strings"
    "regexp"
    "time"
)

func Authentication(scanner *bufio.Scanner, conn net.Conn) string {
    for {
        conn.SetReadDeadline(time.Now().Add(30 * time.Second))
        if !scanner.Scan() {
            err := scanner.Err()
            if err != nil {
                if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                    SendError(conn, 408, "CONNECTION_TIMEOUT")
                }
            }
            return ""
        }
        name := scanner.Text()
        name_list := strings.Fields(name)
        if len(name_list) != 2 {
            SendError(conn, 400, "malformed command CONNECT <name>")
            continue
        }
        if name_list[0] == "CONNECT" {
            name_p := name_list[1]
            if len(name_p) > 12 {
                SendError(conn, 400, "max 12 characters")
                continue
            } else if len(name_p) < 3{
                SendError(conn, 400, "min 3 characters")
            } else if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(name_p){
                SendError(conn, 400, "only letters")
                continue
            }
            SendSuccess(conn, "connected")
            return name_p
        } else {
            SendError(conn, 400, "malformed_command")
        }
    }
}

func ClientAtender(conn net.Conn, hub *models.Hub) {
    scanner := bufio.NewScanner(conn)
    name_p := ""

    defer func() {
        if name_p == "" {
            hub.Unregister <- conn
        }
        conn.Close()
    }()
    name_p = Authentication(scanner, conn)
    if name_p == "" {
        return
    }
    hub.Register <- conn
    player := constructor.NewPlayer(conn.RemoteAddr().String(), name_p)
    hub.Clients[conn] = player
    processFunction := BroadcastMessage(name_p, hub)
    ReadServer(scanner, processFunction)
}
