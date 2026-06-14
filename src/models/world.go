package models

import(
	"sync"
	"fmt"
)

type Quest struct {
    ID           string `json:"quest_id" yaml:"id"`
    Description  string `json:"description" yaml:"description"`
    Type         string `json:"type" yaml:"type"`
    StartItem    string `json:"start_item" yaml:"start_item"`
    RequiredItem string `json:"required_item" yaml:"required_item"`
    TargetNpc    string `json:"target_npc" yaml:"target_npc"`
    Reward       string `json:"reward" yaml:"reward"`
    Status       string `json:"status"`
}

type Npc struct {
	ID			string				`json:"id" yaml:"id"`
	Name		string				`json:"name" yaml:"name"`
	Description string				`json:"description" yaml:"description"`
	Dialogue	[]string			`json:"dialogue" yaml:"dialogue"`
	MaxHP	    int					`json:"hp" yaml:"hp"`
    CurrentHP   int                 `json:"-" yaml:"-"`
    DialogueIdx int                 `json:"-" yaml:"-"`
    AttackDmg   int                 `json:"attack_dmg" yaml:"attack_dmg"`
	IsHostile	bool				`json:"is_hostile" yaml:"is_hostile"`
	QuestID		string				`json:"quest_id" yaml:"quest_id"`
    Combat      bool                `json:"combat" yaml:"combat"`
}

type Item struct {
	ID          string				`json:"id" yaml:"id"`
    Name        string				`json:"name" yaml:"name"`
    Description string 				`json:"description" yaml:"description"`
    Obtainable  bool				`json:"obtainable" yaml:"obtainable"`
} 


type Room struct {
	Mu			sync.RWMutex		`json:"-" yaml:"-"`
	Id			string				`json:"id" yaml:"id"`
	Name		string				`json:"name" yaml:"name"`
	Description	string				`json:"description" yaml:"description"`
	Exist		map[string]string	`json:"exist" yaml:"exist"`
	Players		map[string]*Player	`json:"players" yaml:"-"`
	Items		[]*Item				`json:"items" yaml:"-"`
	Npcs		[]*Npc				`json:"npcs" yaml:"-"`
	YamlItems   []string            `json:"-" yaml:"items"`
    YamlNpcs    []string            `json:"-" yaml:"npcs"`
}

type YamlData struct {
		Rooms []*Room	`yaml:"rooms"`
		Items []*Item	`yaml:"items"`
		Npcs  []*Npc	`yaml:"npcs"`
		Quest []*Quest	`yaml:"quest"`
}

type World struct {
	Rooms		map[string]*Room
	Items		map[string]*Item
	Npcs		map[string]*Npc
	Quest		map[string]*Quest
}

func (w *World) UpdatePlayerRoom(player *Player, targetRoomID string) error {
    nextRoom, roomExists := w.Rooms[targetRoomID]
    if !roomExists {
        return fmt.Errorf("error grave: la sala de destino '%s' no existe en el mundo", targetRoomID)
    }
    actualRoom := player.Room
    if actualRoom != nil {
        actualRoom.Mu.Lock()
        delete(actualRoom.Players, player.Name)
        actualRoom.Mu.Unlock()
    }
    nextRoom.Mu.Lock()
    if nextRoom.Players == nil {
        nextRoom.Players = make(map[string]*Player)
    }
    nextRoom.Players[player.Name] = player
    nextRoom.Mu.Unlock()
    player.Room = nextRoom
    return nil
}


func (w *World) MovePlayer(player *Player, direction string) error {
    actualRoom := player.Room
    if actualRoom == nil {
        return fmt.Errorf("el jugador no está en ninguna sala")
    }

    nextRoomID, exists := actualRoom.Exist[direction]
    if !exists {
        return fmt.Errorf("no hay salida hacia el %s", direction)
    }

    return w.UpdatePlayerRoom(player, nextRoomID)
}
