package models

import (
    "net"
    "sync"
)

type Scope string

const (
    ScopeGlobal Scope = "GLOBAL"
    ScopeRoom   Scope = "ROOM"
    ScopeGroup  Scope = "GROUP"
)

type Message struct {
    Scope    Scope
    Filter   string
    Category string
    Content  string
}

type Player struct {
    mu          sync.RWMutex
    Id          string
    Conn        net.Conn
    Name        string
    Room        *Room
    Group       string
    Inventory   []*Item
}

func (p *Player) GetInventory() []*Item {
    p.mu.RLock()
    defer p.mu.RUnlock()
    inventory := p.Inventory
    return inventory
}