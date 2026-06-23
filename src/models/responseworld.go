package models

import (
	"encoding/json"
)

type QuestResponse struct {
	ID    string
	Title string
}

type WhoResponse struct {
	Room   []string `json:"room"`
	Server int      `json:"server"`
}

type WorldStateResponse struct {
	RoomItems    []string              `json:"room_items"`
	RoNpcsTalk   []string              `json:"room_npcs_talk"`
	RoNpcsHostil []string              `json:"room_npcs_hostil"`
	Inventory    []string              `json:"inventory"`
	PlayerQuests []PlayerQuestResponse `json:"player_quests"`
	NpcQuests    []QuestResponse       `json:"npc_quests"`
}

type PlayerQuestResponse struct {
	QuestID  string `json:"quest_id"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
}

type StatusResponse struct {
	HP     int    `json:"hp"`
	MaxHP  int    `json:"max_hp"`
	Status string `json:"status"`
	Dmg    int    `json:"dmg"`
}

type LookResponse struct {
	Room    RoomData `json:"room"`
	Players []string `json:"players"`
	Items   []string `json:"items"`
	Npcs    []string `json:"npcs"`
}

type RoomData struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
}

func (r *Room) GetStateJSON() (string, error) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	playersList := []string{}
	for name, _ := range r.Players {
		playersList = append(playersList, name)
	}
	itemsList := []string{}
	for _, item := range r.Items {
		if item != nil {
			itemsList = append(itemsList, item.ID)
		}
	}
	npcsList := []string{}
	for _, npc := range r.Npcs {
		if npc != nil {
			npcsList = append(npcsList, npc.ID)
		}
	}
	response := LookResponse{
		Room: RoomData{
			Id:          r.Id,
			Name:        r.Name,
			Description: r.Description,
			Exits:       r.Exist,
		},
		Players: playersList,
		Items:   itemsList,
		Npcs:    npcsList,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
