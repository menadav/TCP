package models

import (
	"answer_protocol/src/logger"
	"answer_protocol/src/speakserver"
	"fmt"
	"net"
	"sync"
)

type PlayerQuest struct {
	QuestID       string
	Status        string
	Progress      string
	NpcKilled     bool
	RoomVisited   bool
	ItemCollected bool
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
		if pq.NpcKilled && q.RequiredKill != "" {
			pq.Progress += "Enemies Defeated; "
		}
		if pq.RoomVisited && q.RequiredRoom != "" {
			pq.Progress += "Location Explored; "
		}
		if pq.ItemCollected && q.RequiredItem != "" {
			pq.Progress += "Items Found; "
		}
		pq.Progress += ")"
	}
}

type Player struct {
	Mu             sync.RWMutex
	Id             string
	Conn           net.Conn
	Name           string
	Room           *Room
	Group          string
	Inventory      []*Item
	Max_HP         int
	HP             int
	Status         string
	Quests         map[string]*PlayerQuest
	NpcDialogueIdx map[string]int
	MsgChan        chan Message
	CombatNpc      string
	Hand           bool
	Dmg            int
}

func (p *Player) UpdateDmg(item *Item) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.Hand = false
	p.Dmg = item.Dmg
}

func (p *Player) VoidDmg() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if !p.Hand {
		p.Hand = true
		p.Dmg = 5
	}
}

func (p *Player) GetName() string {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	return p.Name
}

func (p *Player) GetInventory() []*Item {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	inventory := p.Inventory
	return inventory
}

func (p *Player) GetMaxHp() int {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	return p.Max_HP
}

func (p *Player) GetHp() int {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	return p.HP
}

func (p *Player) GetDmg() int {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	return p.Dmg
}

func (p *Player) GetStatus() string {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	return p.Status
}

func (p *Player) GetCombatNpc() string {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	return p.CombatNpc
}

func (p *Player) SetStatus(status string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Status = status
}

func (p *Player) SetCombatNpc(npcID string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.CombatNpc = npcID
}

func (p *Player) SetHp(hp int) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.HP = hp
}

func (p *Player) SendAsync(category string, content string) {
	if p.MsgChan == nil {
        return
    }
	select {
		case p.MsgChan <- Message{
			Category: category,
			Content:  content,
		}:
		default:
			logger.Warn("[BROADCAST] Client buffer full, message dropped. Player: " + p.Id + " (" + p.Name + ") - Category: " + category)
		}
}

func (p *Player) ListenMsg() {
	for msg := range p.MsgChan {
		speak.SendEvent(p.Conn, msg.Category, msg.Content)
	}
}

func (p *Player) ApplyDamage(dmg int) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.HP -= dmg
}

func (p *Player) GetQuestsResponse() []PlayerQuestResponse {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

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

func (p *Player) GetCurrentRoomNpcIDs() []string {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	if p.Room == nil {
		return nil
	}

	p.Room.Mu.RLock()
	defer p.Room.Mu.RUnlock()

	var ids []string
	for _, npc := range p.Room.Npcs {
		if npc != nil {
			ids = append(ids, npc.ID)
		}
	}
	return ids
}

func (p *Player) GetCurrentRoomItemIDs() []string {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	if p.Room == nil {
		return nil
	}

	p.Room.Mu.RLock()
	defer p.Room.Mu.RUnlock()

	var ids []string
	for _, item := range p.Room.Items {
		if item != nil {
			ids = append(ids, item.ID)
		}
	}
	return ids
}

func (p *Player) GetCurrentRoomNpcIDsTalk() []string {
	p.Room.Mu.RLock()
	defer p.Room.Mu.RUnlock()

	var ids []string
	for _, npc := range p.Room.Npcs {
		if npc != nil && !npc.IsHostile {
			ids = append(ids, npc.ID)
		}
	}
	return ids
}

func (p *Player) GetCurrentRoomNpcIDsHostil() []string {
	p.Room.Mu.RLock()
	defer p.Room.Mu.RUnlock()

	var ids []string
	for _, npc := range p.Room.Npcs {
		if npc != nil && npc.IsHostile && !npc.Combat {
			ids = append(ids, npc.ID)
		}
	}
	return ids
}

func (p *Player) GetInventoryItemIDs() []string {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	ids := make([]string, 0, len(p.Inventory))

	for _, item := range p.Inventory {
		if item != nil {
			ids = append(ids, item.ID)
		}
	}
	return ids
}

func (p *Player) GetPlayerQuestsList() []PlayerQuestResponse {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	var list []PlayerQuestResponse
	for _, pq := range p.Quests {
		if pq != nil {
			list = append(list, PlayerQuestResponse{
				QuestID:  pq.QuestID,
				Status:   pq.Status,
				Progress: pq.Progress,
			})
		}
	}
	return list
}

func (p *Player) GetRoomNpcQuests() []QuestResponse {
	p.Mu.RLock()
	room := p.Room
	p.Mu.RUnlock()

	if room == nil {
		return nil
	}

	room.Mu.RLock()
	defer room.Mu.RUnlock()

	var list []QuestResponse
	for _, npc := range room.Npcs {
		if npc != nil && npc.QuestID != "" {
			list = append(list, QuestResponse{
				ID:    npc.QuestID,
				Title: "Misión: " + npc.QuestID,
			})
		}
	}

	return list
}

func (p *Player) HandleNpcDeath(npcID string, worldQuests map[string]*Quest) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

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
	p.Mu.Lock()
	defer p.Mu.Unlock()

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
	p.Mu.Lock()
	defer p.Mu.Unlock()

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

func (p *Player) AcceptQuest(quest *Quest, startItem *Item) *speak.ErrCode {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if p.Quests == nil {
		p.Quests = make(map[string]*PlayerQuest)
	}
	if existing, ok := p.Quests[quest.ID]; ok {
		if existing.Status == "in_progress" {
			return &speak.ErrQuestAlreadyActive
		}
		if existing.Status == "completed" {
			return &speak.ErrQuestAlreadyDone
		}
	}
	pq := &PlayerQuest{
		QuestID:       quest.ID,
		Status:        "in_progress",
		Progress:      "started",
		NpcKilled:     false,
		RoomVisited:   false,
		ItemCollected: false,
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

func (p *Player) CompleteQuest(quest *Quest, rewardItem *Item) *speak.ErrCode {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	pq, ok := p.Quests[quest.ID]
	if !ok || pq.Status != "in_progress" {
		return &speak.ErrQuestNotActive
	}
	if quest.RequiredKill != "" && !pq.NpcKilled {
		return &speak.ErrObjectiveIncomplete
	}
	if quest.RequiredRoom != "" && !pq.RoomVisited {
		return &speak.ErrObjectiveIncomplete
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
			return &speak.ErrMissingRequiredItem
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


func (p *Player) LeaveQuest(quest *Quest) *speak.ErrCode {
    p.Mu.Lock()
    defer p.Mu.Unlock()

    pq, ok := p.Quests[quest.ID]
    if !ok {
        return &speak.ErrQuestNotActive
    }
    if pq.Progress == "Completed" || pq.Status == "completed" {
        return &speak.ErrQuestAlreadyDone 
    }
    if quest.StartItem != "" {
        idx := -1
        for i, item := range p.Inventory {
            if item.ID == quest.StartItem {
                idx = i
                break
            }
        }
        if idx != -1 {
            p.Inventory = append(p.Inventory[:idx], p.Inventory[idx+1:]...)
        }
    }
    delete(p.Quests, quest.ID)
    return nil
}