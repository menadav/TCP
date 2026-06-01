package world

import (
	"os"
	"fmt"
	"answer_protocol/src/constructor"
	"answer_protocol/src/models"
	"gopkg.in/yaml.v3"
)

func LoadWorld(path string)(*models.World, error){
	var yamlData models.YamlData

	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error read path", err)
		return	nil, err
	}
	err = yaml.Unmarshal(bytes, &yamlData)
	if err != nil {
		fmt.Println("Error read yaml", err)
		return	nil, err
	}
	world := constructor.NewWorld()
	for _, room := range yamlData.Rooms {
        world.Rooms[room.Id] = room
	}
    for _, item := range yamlData.Items {
        world.Items[item.ID] = item
    }
    for _, npc := range yamlData.Npcs {
        world.Npcs[npc.ID] = npc
    }
    return world, nil
}