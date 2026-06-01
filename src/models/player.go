package models

import "net"

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
    Id    string
    Conn  net.Conn
    Name  string
    Room  string
    Group string
}