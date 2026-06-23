package game

import (
	"answer_protocol/src/logger"
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
)

func StartAttack(player *models.Player, target string, h *models.Hub) bool {
	room := player.Room
	if room == nil {
		speak.SendErr(player.Conn, speak.ErrNotInRoom)
		return false
	}

	var npc *models.Npc
	room.Mu.RLock()
	for _, n := range room.Npcs {
		if n.ID == target {
			npc = n
			break
		}
	}
	room.Mu.RUnlock()
	if npc == nil {
		speak.SendErr(player.Conn, speak.ErrNpcNotFound)
		return false
	}
	if !npc.IsHostile {
		speak.SendErr(player.Conn, speak.ErrNpcNotHostile)
		return false
	}
	if npc.CurrentHP <= 0 {
		speak.SendErr(player.Conn, speak.ErrTargetDefeated)
		return false
	}
	if npc.Combat {
		speak.SendErr(player.Conn, speak.ErrAlreadyInCombat)
		return false
	}
	npc.Combat = true
	player.SetStatus("combat")
	player.SetCombatNpc(npc.ID)
	player.SendAsync("COMBAT", "START_COMBAT")
	logger.Info("world change", "event", "combat_start", "player", player.Name, "npc", npc.ID, "room", room.Id)
	return true
}
