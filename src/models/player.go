package models

import (
	"answer_protocol/src/speakserver"
	"fmt"
	"net"
	"sync"
)

type PlayerQuest struct {
    QuestID         string
    Status          string
    Progress        string
    NpcKilled       bool
    RoomVisited     bool
    ItemCollected   bool
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

func (pq *PlayerQuest) UpdateProgressString(q *Quest) {
    if pq.Status == "completed" {
        pq.Progress = "done"
        return
    }
    npcOk := q.RequiredKill == "" || pq.NpcKilled
    roomOk := q.RequiredRoom == "" || pq.RoomVisited
    itemOk := q.RequiredItem == "" || pq.ItemCollected
    if npcOk && roomOk && itemOk {
        pq.Progress = "Completed"
    } else {
        pq.Progress = "0/1"
        if pq.NpcKilled && q.RequiredKill != "" { pq.Progress += "Enemies Defeated; " }
        if pq.RoomVisited && q.RequiredRoom != "" { pq.Progress += "Location Explored; " }
        if pq.ItemCollected && q.RequiredItem != "" { pq.Progress += "Items Found; " }
        pq.Progress += ")"
    }
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
    Hand            bool
    Dmg             int
}

func (p *Player) UpdateDmg(item *Item){
    p.mu.Lock()
    defer p.mu.Unlock()

    p.Hand = false
    p.Dmg = item.Dmg
}

func (p *Player) VoidDmg(){
    p.mu.Lock()
    defer p.mu.Unlock()

    if !p.Hand{
        p.Hand = true
        p.Dmg = 5
    }
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

func (p *Player) GetDmg() int {
    p.mu.RLock()
    defer p.mu.RUnlock()

    return p.Dmg
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
func (p *Player) HandleNpcDeath(npcID string, worldQuests map[string]*Quest) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pq := range p.Quests {
		if pq.Status == "in_progress" {
			if q, ok := worldQuests[pq.QuestID]; ok {
				if q.RequiredKill == "" || q.RequiredKill != npcID {
					continue
				}
				pq.NpcKilled = true
				pq.UpdateProgressString(q)
				p.SendAsync("QUEST", fmt.Sprintf("Quest %s updated: %s", pq.QuestID, pq.Progress))
			}
		}
	}
}

func (p *Player) HandleRoomVisit(roomID string, worldQuests map[string]*Quest) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pq := range p.Quests {
		if pq.Status == "in_progress" {
			if q, ok := worldQuests[pq.QuestID]; ok {
				if q.RequiredRoom == "" || q.RequiredRoom != roomID {
					continue
				}
				pq.RoomVisited = true
				pq.UpdateProgressString(q)
				p.SendAsync("QUEST", fmt.Sprintf("Quest %s updated: %s", pq.QuestID, pq.Progress))
			}
		}
	}
}

func (p *Player) HandleItemCollection(itemID string, worldQuests map[string]*Quest) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, pq := range p.Quests {
		if pq.Status == "in_progress" {
			if q, ok := worldQuests[pq.QuestID]; ok {
				if q.RequiredItem == "" || q.RequiredItem != itemID {
					continue
                }
				pq.ItemCollected = true
				pq.UpdateProgressString(q)
				p.SendAsync("QUEST", fmt.Sprintf("Quest %s updated: %s", pq.QuestID, pq.Progress))
		    }
	    }
    }
}

func (p *Player) AcceptQuest(quest *Quest, startItem *Item) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Quests == nil {
		p.Quests = make(map[string]*PlayerQuest)
	}
	if existing, ok := p.Quests[quest.ID]; ok {
		if existing.Status == "in_progress" {
			return fmt.Errorf("QUEST_ALREADY_IN_PROGRESS")
		}
		if existing.Status == "completed" {
			return fmt.Errorf("QUEST_ALREADY_COMPLETED")
		}
	}
	pq := &PlayerQuest{
		QuestID:        quest.ID,
		Status:         "in_progress",
		Progress:       "started",
		NpcKilled:      false,
		RoomVisited:    false,
		ItemCollected:  false,
	}
	if quest.RequiredRoom != "" && p.Room != nil && p.Room.Id == quest.RequiredRoom {
		pq.RoomVisited = true
	}
	if quest.RequiredItem != "" {
		for _, item := range p.Inventory {
			if item.ID == quest.RequiredItem {
				pq.ItemCollected = true
				break
			}
		}
	}
	pq.UpdateProgressString(quest)
	p.Quests[quest.ID] = pq
	if startItem != nil {
		if p.Inventory == nil {
			p.Inventory = make([]*Item, 0)
		}
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
    if quest.RequiredKill != "" && !pq.NpcKilled {
		return fmt.Errorf("OBJECTIVE_NPC_NOT_KILLED")
	}
    if quest.RequiredRoom != "" && !pq.RoomVisited {
		return fmt.Errorf("OBJECTIVE_ROOM_NOT_VISITED")
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