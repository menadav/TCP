package models

import (
	"answer_protocol/src/speakserver"
	"net"
    "sync"
)

type Hub struct {
    Register    chan *Player
    Unregister  chan *Player
    Broadcast   chan Message
    Clients     map[net.Conn]*Player
    Groups      map[string]*Group
    mu          sync.RWMutex
    World       *World
}

func (h *Hub) Run(){
    for {
        select {
            case player := <- h.Register:
                h.mu.Lock()
                enter := "PRESENCE ENTER " + player.Name
                for _, p := range h.Clients {
                    if p.Room.Id == player.Room.Id {
                        speak.SendEvent(p.Conn, "ROOM", enter)
                    }
                }
                h.Clients[player.Conn] = player
                h.mu.Unlock()
            case player := <- h.Unregister:
                h.mu.Lock()
                leave := "PRESENCE LEAVE " + player.Name
                for _, p := range h.Clients {
                        if p.Room.Id == player.Room.Id {
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
                h.mu.Unlock()
            case msg := <- h.Broadcast:
                h.mu.RLock()
                switch msg.Scope {
                case ScopeGlobal:
                    for _, player := range h.Clients {
                        speak.SendEvent(player.Conn, msg.Category, msg.Content)
                    }
                case ScopeRoom:
                    for _, player := range h.Clients {
                        if player.Room.Id == msg.Filter {
                            speak.SendEvent(player.Conn, msg.Category, msg.Content)
                        }
                    }
                case ScopeGroup:
                    if group, exist :=  h.Groups[msg.Filter]; exist {
                        group.Broadcast(msg)
                    }
                }
                h.mu.RUnlock()
            }
    }
}

func (h *Hub) GetOnlinePlayersNames() []string {
    h.mu.RLock()
    defer h.mu.RUnlock()
    playerNames := make([]string, 0)
    for _, player := range h.Clients {
        if player != nil {
            playerNames = append(playerNames, player.GetName())
        }
    }
    return playerNames
}
