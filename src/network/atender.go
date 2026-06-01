package network

import (
    "answer_protocol/src/models"
    "answer_protocol/src/constructor"
    "answer_protocol/src/parse"
    "bufio"
    "net"
)

func ClientAtender(conn net.Conn, hub *models.Hub) {
    scanner := bufio.NewScanner(conn)
    name_p := ""
    name_p = Authentication(scanner, conn, hub)
    if name_p == "" {
        conn.Close()
        return
    }
    player := constructor.NewPlayer(conn.RemoteAddr().String(), conn, name_p, "loc.start")
    defer func() {
        if name_p != "" {
            hub.Unregister <- player
        }
        conn.Close()
    }()
    hub.Register <- player
    processFunction := func(line string) {
		parse.ParseCommandCli(line, player, hub)
	}
    ReadServer(scanner, processFunction)
}
