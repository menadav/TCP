package game

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"encoding/json"
)

func ShowRoom(player *models.Player, h *models.Hub) {
	room := player.Room
	if room == nil {
		speak.SendError(player.Conn, 400, "NOT_ROOM")
		return
	}
	roomJSON, err := room.GetStateJSON()
	if err != nil {
		speak.SendError(player.Conn, 900, "INTERNAL_ERROR")
		return
	}
	speak.SendSuccess(player.Conn, roomJSON)
}

func ShowInventory(player *models.Player) {
	inventory := player.GetInventory()
	bytesJSON, err := json.Marshal(inventory)
	if err != nil {
		speak.SendError(player.Conn, 900, "INTERNAL_ERROR")
		return
	}
	speak.SendSuccess(player.Conn, string(bytesJSON))
}


func ShowStatus(player *models.Player){
	status := models.StatusResponse {
		HP:		player.GetHp(),
		MaxHP:	player.GetMaxHp(),
		Status:	player.GetStatus(),
		Dmg:	player.GetDmg(),
	}
	bytesJSON, err := json.Marshal(status)
	if err != nil {
		speak.SendError(player.Conn, 900, "INTERNAL_ERROR")
		return
	}
	speak.SendSuccess(player.Conn, string(bytesJSON))
}

func ShowQuest(player *models.Player) {
	questsList := player.GetQuestsResponse()
    bytesJSON, err := json.Marshal(questsList)
    if err != nil {
        speak.SendError(player.Conn, 900, "INTERNAL_ERROR")
        return
    }
    speak.SendSuccess(player.Conn, string(bytesJSON))
}

func ShowWho(player *models.Player, h *models.Hub) {
	nameList := h.GetOnlinePlayersNames()
	bytesJSON, err := json.Marshal(nameList)
    if err != nil {
        speak.SendError(player.Conn, 900, "INTERNAL_ERROR")
        return
    }
    speak.SendSuccess(player.Conn, "players=" + string(bytesJSON))
}