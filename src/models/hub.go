package models

import (
    "net"
    "answer_protocol/src/speakserver"
)

type Hub struct {
    Register   chan *Player
    Unregister chan *Player
    Broadcast  chan Message
    Clients    map[net.Conn]*Player
    Groups     map[string]*Group
}

func (h *Hub) Run(){
    for {
        select {
            case player := <- h.Register:
                enter := "PRESENCE ENTER " + player.Name
                for _, p := range h.Clients {
                    if p.Room == player.Room {
                        speak.SendEvent(p.Conn, "ROOM", enter)
                    }
                }
                h.Clients[player.Conn] = player
            case player := <- h.Unregister:
                leave := "PRESENCE LEAVE " + player.Name
                for _, p := range h.Clients {
                        if p.Room == player.Room {
                            if player != p {
                                speak.SendEvent(p.Conn, "ROOM", leave)
                            }
                        }
                    }
                if player.Group != "" {
                    if group, exist := h.Groups[player.Group]; exist {
                        restantes := group.RemoveMember(player.Conn)
                        if restantes == 0 || group.Leader == player {
                            delete(h.Groups, player.Group)
                        }
                    }
                }
			    delete(h.Clients, player.Conn)
            case msg := <- h.Broadcast:
                switch msg.Scope {
                case ScopeGlobal:
                    for _, player := range h.Clients {
                        speak.SendEvent(player.Conn, msg.Category, msg.Content)
                    }
                case ScopeRoom:
                    for _, player := range h.Clients {
                        if player.Room == msg.Filter {
                            speak.SendEvent(player.Conn, msg.Category, msg.Content)
                        }
                    }
                case ScopeGroup:
                    if group, exist :=  h.Groups[msg.Filter]; exist {
                        group.Broadcast(msg)
                    }
                }
            }
    }
}

