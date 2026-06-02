package models

import(
	"sync"
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
	Players		map[string]bool		`json:"players" yaml:"-"`
	Items		[]string			`json:"items" yaml:"items"`
	Npcs		[]string			`json:"npcs" yaml:"npcs"`
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