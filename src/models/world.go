package models

import(
	"sync"
	"encoding/json"
)

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


type Room struct{
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

type World struct{
	Rooms		map[string]*Room
	Items		map[string]*Item
	Npcs		map[string]*Npc
}

type YamlData struct {
		Rooms []*Room `yaml:"rooms"`
		Items []*Item `yaml:"items"`
		Npcs  []*Npc  `yaml:"npcs"`
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