package models

import (
    "net"
    "answer_protocol/src/network"
)

type Scope string

const (
    ScopeGlobal Scope = "GLOBAL"
    ScopeRoom   Scope = "ROOM"
    ScopeGroup  Scope = "GROUP"
)

type Message struct{
    Scope   Scope
    Filter  string
    Category
    Content string
}

type Player struct{
    Id      string
    Conn    net.Conn
    Name    string
    Room    string
    Group   string
}

type Hub struct {
    Register   chan *Player
    Unregister chan *Player
    Broadcast  chan Message
    Clients    map[net.Conn]*Player
}

func (h *Hub) Run(){
    for {
        select {
            case player := <- h.Register:
                h.Clients[player.Conn] = player
            case player := <- h.Unregister:
                delete(h.Clients, player.Conn)
            case msg := <- h.Broadcast:
                for _, player := range h.Clients{
                    switch msg.Scope {
                    case ScopeGlobal:
                        player.Conn.Write([]byte(msg.Content))
                    case ScopeRoom:
                        if player.Room == msg.Filter {
                            player.Conn.Write([]byte(msg.Content))
                        }
                    case ScopeGroup:
                        if player.Group != "" && player.Group == msg.Filter {
                            network.SendEvent(player.Conn, msg.Category, msg.Content)
                        }
                    }
            }
        }
    }
}
