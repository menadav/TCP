package game

import (
	"answer_protocol/src/models"
	"fmt"
)

func respawnPlayer(player *models.Player, h *models.Hub) {
	maxHp := player.Max_HP
	hp := maxHp / 2

	player.SetHp(hp)
	err := h.World.UpdatePlayerRoom(player, "start")
	if err != nil {
		fmt.Printf("[ERROR] There are not respawn %v\n", err)
		return
	}
	player.SendAsync("COMBAT", "¡You have defeat in combat!")
	player.SendAsync("INFO", "You have respawn"+player.Room.Name)
	ShowRoom(player, h)
}

func handleNpcDeath(player *models.Player, npc *models.Npc, h *models.Hub) {
	room := player.Room
	if room == nil {
		return
	}

	room.Mu.Lock()
	idx := -1
	for i, n := range room.Npcs {
		if n.ID == npc.ID {
			idx = i
			break
		}
	}
	if idx != -1 {
		room.Npcs = append(room.Npcs[:idx], room.Npcs[idx+1:]...)
	}
	room.Mu.Unlock()
	player.SendAsync("COMBAT", fmt.Sprintf("VICTORY %s has been defeated!", npc.Name))
	h.Broadcast <- models.Message{
		Scope:    models.ScopeRoom,
		Filter:   room.Id,
		Category: "ROOM",
		Content:  fmt.Sprintf("¡%s have defeat to %s!", player.GetName(), npc.Name),
	}
}
