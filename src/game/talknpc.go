package game

import (
    "answer_protocol/src/models"
    "answer_protocol/src/speakserver"
)

func TalkNpc(player *models.Player, npcID string){
	actualRoom := player.Room
	if actualRoom == nil {
        speak.SendError(player.Conn, 400, "NOT_IN_ROOM")
        return
    }
	var npc *models.Npc
    actualRoom.Mu.RLock()
    for _, n := range actualRoom.Npcs {
        if n.ID == npcID {
            npc = n
            break
        }
    }
    actualRoom.Mu.RUnlock()
	if npc == nil {
        speak.SendError(player.Conn, 404, "NPC_NOT_FOUND")
        return
    }
    if npc.IsHostile {
        speak.SendError(player.Conn, 403, "NPC_IS_HOSTILE")
        return
    }
    if len(npc.Dialogue) == 0 {
        speak.SendError(player.Conn, 404, "NPC_HAS_NO_DIALOGUE")
        return
    }
	if player.NpcDialogueIdx == nil {
        player.NpcDialogueIdx = make(map[string]int)
    }
    idx := player.NpcDialogueIdx[npcID]
    line := npc.Dialogue[idx]
    player.NpcDialogueIdx[npcID] = (idx + 1) % len(npc.Dialogue)
    player.MsgChan <- models.Message{
		Scope:    models.ScopeRoom,
		Category: "NPC",
		Content:  npc.Name + " " + line,
	}
}