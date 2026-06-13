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

type Message struct {
    Scope    Scope
    Filter   string
    Category string
    Content  string
}
type PlayerQuest struct {
    QuestID  string
    Status   string
    Progress string
}

func (pq *PlayerQuest) GetQuestID() string {
    return pq.QuestID
}

func (pq *PlayerQuest) GetStatus() string {
    return pq.Status
}

func (pq *PlayerQuest) GetProgress() string {
    return pq.Progress
}

type Player struct {
    mu              sync.RWMutex
    Id              string
    Conn            net.Conn
    Name            string
    Room            *Room
    Group           string
    Inventory       []*Item
    Max_HP          int
    HP              int
    Status          string
    Quests          map[string]*PlayerQuest
    NpcDialogueIdx  map[string]int
    MsgChan         chan Message         
}

func (p *Player) GetName() string {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.Name
}

func (p *Player) GetInventory() []*Item {
    p.mu.RLock()
    defer p.mu.RUnlock()
    inventory := p.Inventory
    return inventory
}

func (p *Player) GetMaxHp() int {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.Max_HP
}

func (p *Player) GetHp() int {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.HP
}

func (p *Player) GetStatus() string {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.Status
}

func (p *Player) ListenMsg() {
    for msg := range p.MsgChan {
        speak.SendEvent(p.Conn, msg.Category, msg.Content) 
    }
}

func (p *Player) GetQuestsResponse() []PlayerQuestResponse {
    p.mu.RLock()
    defer p.mu.RUnlock()
    responseList := make([]PlayerQuestResponse, 0)
    for _, pQuest := range p.Quests {
        if pQuest != nil {
            questResp := PlayerQuestResponse{
                QuestID:  pQuest.GetQuestID(),
                Status:   pQuest.GetStatus(),
                Progress: pQuest.GetProgress(),
            }
            responseList = append(responseList, questResp)
        }
    }
    return responseList
}
