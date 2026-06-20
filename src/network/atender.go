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
    start := hub.World.Rooms["start"]
    player := constructor.NewPlayer(conn.RemoteAddr().String(), conn, name_p, start)
    if start != nil {
        start.Mu.Lock()  
        if start.Players == nil {
            start.Players = make(map[string]*models.Player)
        }
        start.Players[player.Name] = player
        start.Mu.Unlock()
    }
    defer func() {
        if name_p != "" {
            if player.Room != nil {
                player.Room.Mu.Lock()
                delete(player.Room.Players, player.Name)
                player.Room.Mu.Unlock()
            }
            hub.Unregister <- player
        }
        close(player.MsgChan)
        conn.Close()
    }()
    hub.Register <- player
    go player.ListenMsg()
    processFunction := func(line string) {
		parse.ParseCommandCli(line, player, hub)
	}
    StartScanner(scanner, processFunction, conn)
}