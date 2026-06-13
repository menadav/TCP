package models

import (
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
                h.innerNotifyRoomEnter(player, player.Room.Id)
                h.Clients[player.Conn] = player
                h.mu.Unlock()
            case player := <- h.Unregister:
                h.mu.Lock()
                h.innerNotifyRoomLeave(player, player.Room.Id)
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
                        player.MsgChan <- msg
                    }
                case ScopeRoom:
                    for _, player := range h.Clients {
                        if player.Room.Id == msg.Filter {
                            player.MsgChan <- msg
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

func (h *Hub) NotifyRoomLeave(player *Player, roomID string) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    h.innerNotifyRoomLeave(player, roomID)
}

func (h *Hub) NotifyRoomEnter(player *Player, roomID string) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    h.innerNotifyRoomEnter(player, roomID)
}

func (h *Hub) innerNotifyRoomLeave(player *Player, roomID string) {
    leaveMsg := "PRESENCE LEAVE " + player.Name
    for _, p := range h.Clients {
        if p.Room.Id == roomID && p != player {
            p.MsgChan <- Message{
                Scope:    ScopeRoom,
                Category: "ROOM",
                Content:  leaveMsg,
            }
        }
    }
}

func (h *Hub) innerNotifyRoomEnter(player *Player, roomID string) {
    enterMsg := "PRESENCE ENTER " + player.Name
    for _, p := range h.Clients {
        if p.Room.Id == roomID && p != player {
            p.MsgChan <- Message{
                Scope:    ScopeRoom,
                Category: "ROOM",
                Content:  enterMsg,
            }
        }
    }
}