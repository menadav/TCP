package models

import (
	"answer_protocol/src/speakserver"
	"fmt"
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
    CombatNpc       string
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

func (p *Player) GetCombatNpc() string {
    p.mu.RLock()
    defer p.mu.RUnlock()

    return p.CombatNpc
}

func (p *Player) SetStatus(status string) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.Status = status
}

func (p *Player) SetCombatNpc(npcID string) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.CombatNpc = npcID
}

func (p *Player) SetHp(hp int){
    p.mu.Lock()
    defer p.mu.Unlock()

    p.HP = hp
}

func (p *Player) SendAsync(category string, content string) {
    if p.MsgChan != nil {
        p.MsgChan <- Message{
            Category: category,
            Content:  content,
        }
    }
}

func (p *Player) ListenMsg() {
    for msg := range p.MsgChan {
        speak.SendEvent(p.Conn, msg.Category, msg.Content) 
    }
}

func (p *Player) ApplyDamage(dmg int){
    p.mu.Lock()
    defer p.mu.Unlock()

    p.HP -= dmg
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

func (p *Player) AcceptQuest(questID string, startItem *Item) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if existing, ok := p.Quests[questID]; ok {
        if existing.Status == "in_progress" {
            return fmt.Errorf("QUEST_ALREADY_IN_PROGRESS")
        }
        if existing.Status == "completed" {
            return fmt.Errorf("QUEST_ALREADY_COMPLETED")
        }
    }
    p.Quests[questID] = &PlayerQuest{
        QuestID:  questID,
        Status:   "in_progress",
        Progress: "started",
    }
    if startItem != nil {
        p.Inventory = append(p.Inventory, startItem)
    }
    return nil
}

func (p *Player) CompleteQuest(quest *Quest, rewardItem *Item) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    pq, ok := p.Quests[quest.ID]
    if !ok || pq.Status != "in_progress" {
        return fmt.Errorf("QUEST_NOT_IN_PROGRESS")
    }
    if quest.RequiredItem != "" {
        idx := -1
        for i, item := range p.Inventory {
            if item.ID == quest.RequiredItem {
                idx = i
                break
            }
        }
        if idx == -1 {
            return fmt.Errorf("MISSING_REQUIRED_ITEM")
        }
        p.Inventory = append(p.Inventory[:idx], p.Inventory[idx+1:]...)
    }
    if rewardItem != nil {
        p.Inventory = append(p.Inventory, rewardItem)
    }
    pq.Status = "completed"
    pq.Progress = "done"
    return nil
}