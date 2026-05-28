package models

import (
    "net"
    "sync"
    "answer_protocol/src/speakserver"
)

type Scope string

const (
    ScopeGlobal Scope = "GLOBAL"
    ScopeRoom   Scope = "ROOM"
    ScopeGroup  Scope = "GROUP"
)

type Message struct{
    Scope       Scope
    Filter      string
    Category    string
    Content     string
}

type Player struct{
    Id      string
    Conn    net.Conn
    Name    string
    Room    string
    Group   string
}
type Group struct {
	Id         string
	Leader     *Player
	mu         sync.RWMutex
	Members    map[net.Conn]*Player
}

type Hub struct {
    Register   chan *Player
    Unregister chan *Player
    Broadcast  chan Message
    Clients    map[net.Conn]*Player
    Groups     map[string]*Group
}

func (g *Group) AddMember(conn net.Conn, p *Player) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Members[conn] = p
}

func (g *Group) RemoveMember(conn net.Conn) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.Members, conn)
	return len(g.Members)
}

func (g *Group) Broadcast(msg Message) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, player := range g.Members {
		speak.SendEvent(player.Conn, msg.Category, msg.Content)
	}
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

