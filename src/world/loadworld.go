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
    for _, quest := range yamlData.Quest {
        world.Quest[quest.ID] = quest
    }
	for _, room := range world.Rooms {
        for _, itemID := range room.YamlItems {
            if itemReal, existe := world.Items[itemID]; existe {
                room.Items = append(room.Items, itemReal)
            } else {
                fmt.Printf("Advertencia: El ítem '%s' requerido en la sala '%s' no existe en la sección global de ítems.\n", itemID, room.Id)
            }
        }
        for _, npcID := range room.YamlNpcs {
            if npcReal, existe := world.Npcs[npcID]; existe {
                room.Npcs = append(room.Npcs, npcReal)
            } else {
                fmt.Printf("Advertencia: El NPC '%s' requerido en la sala '%s' no existe en la sección global de NPCs.\n", npcID, room.Id)
            }
        }
    }
    return world, nil
}