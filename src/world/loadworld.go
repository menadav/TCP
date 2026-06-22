package world

import (
	"answer_protocol/src/constructor"
	"answer_protocol/src/logger"
	"answer_protocol/src/models"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadWorld(path string) (*models.World, error) {
	var yamlData models.YamlData

	bytes, err := os.ReadFile(path)
	if err != nil {
		logger.Error("world read failed", "path", path, "err", err)
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &yamlData)
	if err != nil {
		logger.Error("world parse failed", "path", path, "err", err)
		return nil, err
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
	for _, npc := range world.Npcs {
		npc.CurrentHP = npc.MaxHP
	}
	for _, quest := range yamlData.Quest {
		world.Quest[quest.ID] = quest
	}
	for _, room := range world.Rooms {
		for _, itemID := range room.YamlItems {
			if itemReal, existe := world.Items[itemID]; existe {
				room.Items = append(room.Items, itemReal)
			} else {
				logger.Warn("world load: item not found", "item", itemID, "room", room.Id)
			}
		}
		for _, npcID := range room.YamlNpcs {
			if npcReal, existe := world.Npcs[npcID]; existe {
				room.Npcs = append(room.Npcs, npcReal)
			} else {
				logger.Warn("world load: npc not found", "npc", npcID, "room", room.Id)
			}
		}
	}
	return world, nil
}
