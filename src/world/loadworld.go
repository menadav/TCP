package world

import (
	"answer_protocol/src/constructor"
	"answer_protocol/src/logger"
	"answer_protocol/src/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
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

	if err := validateWorld(world, yamlData); err != nil {
		logger.Error("world validation failed", "path", path, "err", err)
		return nil, err
	}

	for _, room := range world.Rooms {
		for _, itemID := range room.YamlItems {
			room.Items = append(room.Items, world.Items[itemID])
		}
		for _, npcID := range room.YamlNpcs {
			room.Npcs = append(room.Npcs, world.Npcs[npcID])
		}
	}
	return world, nil
}

func validateWorld(world *models.World, data models.YamlData) error {
	var problems []string

	if _, ok := world.Rooms["start"]; !ok {
		problems = append(problems, "missing required spawn room 'start'")
	}

	for _, room := range data.Rooms {
		for dir, target := range room.Exist {
			if _, ok := world.Rooms[target]; !ok {
				problems = append(problems, fmt.Sprintf("room %q exit %s points to unknown room %q", room.Id, dir, target))
			}
		}
		for _, itemID := range room.YamlItems {
			if _, ok := world.Items[itemID]; !ok {
				problems = append(problems, fmt.Sprintf("room %q references unknown item %q", room.Id, itemID))
			}
		}
		for _, npcID := range room.YamlNpcs {
			if _, ok := world.Npcs[npcID]; !ok {
				problems = append(problems, fmt.Sprintf("room %q references unknown npc %q", room.Id, npcID))
			}
		}
	}

	for _, quest := range data.Quest {
		if quest.StartItem != "" {
			if _, ok := world.Items[quest.StartItem]; !ok {
				problems = append(problems, fmt.Sprintf("quest %q start_item %q is not a known item", quest.ID, quest.StartItem))
			}
		}
		if quest.RequiredItem != "" {
			if _, ok := world.Items[quest.RequiredItem]; !ok {
				problems = append(problems, fmt.Sprintf("quest %q required_item %q is not a known item", quest.ID, quest.RequiredItem))
			}
		}
		if quest.Reward != "" {
			if _, ok := world.Items[quest.Reward]; !ok {
				problems = append(problems, fmt.Sprintf("quest %q reward %q is not a known item", quest.ID, quest.Reward))
			}
		}
		if quest.TargetNpc != "" {
			if _, ok := world.Npcs[quest.TargetNpc]; !ok {
				problems = append(problems, fmt.Sprintf("quest %q target_npc %q is not a known npc", quest.ID, quest.TargetNpc))
			}
		}
		if quest.RequiredKill != "" {
			if _, ok := world.Npcs[quest.RequiredKill]; !ok {
				problems = append(problems, fmt.Sprintf("quest %q required_kill %q is not a known npc", quest.ID, quest.RequiredKill))
			}
		}
		if quest.RequiredRoom != "" {
			if _, ok := world.Rooms[quest.RequiredRoom]; !ok {
				problems = append(problems, fmt.Sprintf("quest %q required_room %q is not a known room", quest.ID, quest.RequiredRoom))
			}
		}
	}

	if len(problems) > 0 {
		return fmt.Errorf("world validation failed:\n  - %s", strings.Join(problems, "\n  - "))
	}
	return nil
}
