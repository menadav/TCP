package game

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"fmt"
	"math/rand"
)

func Attack(player *models.Player, h *models.Hub) {
	room := player.Room
	targetID := player.GetCombatNpc()

	var npc *models.Npc
	room.Mu.RLock()
	for _, n := range room.Npcs {
		if n.ID == targetID {
			npc = n
			break
		}
	}
	room.Mu.RUnlock()
	if npc == nil {
		player.SetStatus("healthy")
		player.SetCombatNpc("")
		speak.SendErr(player.Conn, speak.ErrTargetGone)
		return
	}
	playerDmg := player.Dmg + rand.Intn(5)
	npc.CurrentHP -= playerDmg
	if npc.CurrentHP <= 0 {
		player.HandleNpcDeath(npc.ID, h.World.Quest)
		player.SetStatus("healthy")
		player.SetCombatNpc("")
		npc.Combat = false
		handleNpcDeath(player, npc, h)
		return
	}
	npcDmg := npc.AttackDmg
	player.ApplyDamage(npcDmg)
	currentRoomPlayerHp := player.GetHp()
	h.Broadcast <- models.Message{
        Scope:    models.ScopeRoom,
        Filter:   player.Room.Id,
        Category: "COMBAT",
        Content:  fmt.Sprintf("player=%s dealt=%d received=%d npc_hp=%d player_hp=%d", player.Name, playerDmg, npcDmg, npc.CurrentHP, currentRoomPlayerHp),
    }
	if currentRoomPlayerHp <= 0 {
		player.SetStatus("healthy")
		player.SetCombatNpc("")
		npc.Combat = false
		respawnPlayer(player, h)
	}
}

func Defend(player *models.Player, h *models.Hub) {
	room := player.Room
	targetID := player.GetCombatNpc()

	var npc *models.Npc
	room.Mu.RLock()
	for _, n := range room.Npcs {
		if n.ID == targetID {
			npc = n
			break
		}
	}
	room.Mu.RUnlock()
	if npc == nil {
		player.SetStatus("healthy")
		player.SetCombatNpc("")
		speak.SendErr(player.Conn, speak.ErrTargetGone)
		return
	}
	reducedDmg := npc.AttackDmg / 2
	player.ApplyDamage(reducedDmg)

	h.Broadcast <- models.Message{
        Scope:    models.ScopeRoom,
        Filter:   player.Room.Id,
        Category: "COMBAT",
        Content:  fmt.Sprintf("player=%s defended received=%d npc_hp=%d player_hp=%d", player.Name, reducedDmg, npc.CurrentHP, player.GetHp()),
    }
	if player.GetHp() <= 0 {
		player.SetStatus("healthy")
		player.SetCombatNpc("")
		npc.Combat = false
		respawnPlayer(player, h)
	}
}

func Flee(player *models.Player, h *models.Hub) {
	room := player.Room
	targetID := player.GetCombatNpc()

	var npc *models.Npc
	room.Mu.RLock()
	for _, n := range room.Npcs {
		if n.ID == targetID {
			npc = n
			break
		}
	}
	room.Mu.RUnlock()
	if player.GetStatus() != "combat" {
		speak.SendErr(player.Conn, speak.ErrNotInCombat)
		return
	}
	player.SetStatus("healthy")
	player.SetCombatNpc("")
	npc.Combat = false
	speak.SendSuccess(player.Conn, "fled")
}
