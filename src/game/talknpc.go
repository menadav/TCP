package game

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
)

func TalkNpc(player *models.Player, npcID string) {
	actualRoom := player.Room
	if actualRoom == nil {
		speak.SendErr(player.Conn, speak.ErrNotInRoom)
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
		speak.SendErr(player.Conn, speak.ErrNpcNotFound)
		return
	}
	if npc.IsHostile {
		speak.SendErr(player.Conn, speak.ErrNpcHostile)
		return
	}
	if len(npc.Dialogue) == 0 {
		speak.SendErr(player.Conn, speak.ErrNpcNoDialogue)
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
