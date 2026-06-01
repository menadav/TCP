package models

import (
    "net"
    "sync"
    "answer_protocol/src/speakserver"
)

type Group struct {
	Id         string
	Leader     *Player
	mu         sync.RWMutex
	Members    map[net.Conn]*Player
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