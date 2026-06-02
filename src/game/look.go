package game

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"encoding/json"
)

func LookRoom(player *models.Player, h *models.Hub) {
	room := player.Room
	if room == nil {
		speak.SendError(player.Conn, 400, "NOT_ROOM")
		return
	}
	roomJSON, err := room.GetStateJSON()
	if err != nil {
		speak.SendError(player.Conn, 500, "INTERNAL_ERROR")
		return
	}
	speak.SendSuccess(player.Conn, roomJSON)
}

func LookInventory(player *models.Player) {
	inventory := player.GetInventory()
	bytesJSON, err := json.Marshal(inventory)
	if err != nil {
		speak.SendError(player.Conn, 500, "INTERNAL_ERROR")
		return
	speak.SendSuccess(player.Conn, string(bytesJSON))
	}
}

func LookStatus(player *models.Player){
	
}