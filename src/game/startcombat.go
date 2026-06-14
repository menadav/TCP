package game

import(
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
)

func StartAttack(player *models.Player, target string, h *models.Hub) bool {
    room := player.Room
    if room == nil {
        speak.SendError(player.Conn, 400, "NOT_IN_ROOM")
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
        speak.SendError(player.Conn, 404, "NPC_NOT_FOUND")
        return false
    }
    if !npc.IsHostile {
        speak.SendError(player.Conn, 403, "NPC_NOT_HOSTILE")
        return false
    }
    if npc.CurrentHP <= 0 {
        speak.SendError(player.Conn, 403, "NPC_ALREADY_DEFEATED")
        return false
    }
    if npc.Combat{
        speak.SendError(player.Conn, 403, "NPC_ALREADY_COMBAT")
        return false
    }
    npc.Combat = true
    player.SetStatus("combat")
    player.SetCombatNpc(npc.ID)
    player.SendAsync("COMBAT", "START_COMBAT")
    return true
}