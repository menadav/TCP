package models

import(
	"sync"
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
	Quest		map[string]*Quest
}

type YamlData struct {
		Rooms []*Room	`yaml:"rooms"`
		Items []*Item	`yaml:"items"`
		Npcs  []*Npc	`yaml:"npcs"`
		Quest []*Quest	`yaml:"quest"`
}