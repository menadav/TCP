package models

import(
	"sync"
	"fmt"
)

type Quest struct {
    ID          string	`json:"quest_id" yaml:"id"`
    Description string	`json:"description" yaml:"description"`
    Reward      string	`json:"reward" yaml:"reward"`
    Status      string	`json:"status"`
}

type Npc struct {
	ID			string				`json:"id" yaml:"id"`
	Name		string				`json:"name" yaml:"name"`
	Description string				`json:"description" yaml:"description"`
	Dialogue	[]string			`json:"dialogue" yaml:"dialogue"`
	HP			int					`json:"hp" yaml:"hp"`
	IsHostile	bool				`json:"is_hostile" yaml:"is_hostile"`
	QuestID		string				`json:"quest_id" yaml:"quest_id"`
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

type World struct {
	Rooms		map[string]*Room
	Items		map[string]*Item
	Npcs		map[string]*Npc
	Quest		map[string]*Quest
}

type YamlData struct {
		Rooms []*Room	`yaml:"rooms"`
		Items []*Item	`yaml:"items"`
		Npcs  []*Npc	`yaml:"npcs"`
		Quest []*Quest	`yaml:"quest"`
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
    actualRoom.Mu.Lock()
    delete(actualRoom.Players, player.Name)
    actualRoom.Mu.Unlock()
    nextRoom, roomExists := w.Rooms[nextRoomID]
    if !roomExists {
        return fmt.Errorf("error grave: la sala %s no existe en el mundo", nextRoomID)
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