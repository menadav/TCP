package network

import (
    "answer_protocol/internal/models"
    "answer_protocol/internal/constructor"
    "bufio"
    "net"
    "strings"
)

func Authentication(scanner *bufio.Scanner, conn net.Conn) string {
    for {
        if !scanner.Scan() {
            return ""
        }
        name := scanner.Text()
        name_list := strings.Fields(name)
        if len(name_list) != 2 {
            conn.Write([]byte("[S] ERR malformed_command\n"))
            continue
        }
        if name_list[0] == "CONNECT" {
            name_p := name_list[1]
            conn.Write([]byte("[S] OK connected\n"))
            return name_p
        } else {
            conn.Write([]byte("[S] ERR unauthorized\n"))
        }
    }
}

func ClientAtender(conn net.Conn, hub *models.Hub) {
    scanner := bufio.NewScanner(conn)

    defer func() {
        hub.Unregister <- conn
        conn.Close()
    }()
    name_p := Authentication(scanner, conn)
    if name_p == "" {
        return
    }
    hub.Register <- conn
    player := constructor.NewPlayer(conn.RemoteAddr().String(), name_p)
    hub.Clients[conn] = player
    processFunction := BroadcastMessage(name_p, hub)
    ReadServer(scanner, processFunction)
}