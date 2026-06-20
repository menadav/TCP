package game

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"strings"
)

func TakeItem(player *models.Player, query string, hub *models.Hub) {
	actualRoom := player.Room
	actualRoom.Mu.Lock()
	defer actualRoom.Mu.Unlock()

	for i, item := range actualRoom.Items {
		if strings.EqualFold(item.Name, query) || item.ID == query {
			if !item.Obtainable {
				speak.SendErr(player.Conn, speak.ErrNotObtainable)
				return
			}
			if item.Hand && !player.Hand {
				speak.SendErr(player.Conn, speak.ErrHandsFull)
				return
			}
			player.UpdateDmg(item)
			actualRoom.Items = append(actualRoom.Items[:i], actualRoom.Items[i+1:]...)
			player.Inventory = append(player.Inventory, item)
			player.HandleItemCollection(item.ID, hub.World.Quest)
			speak.SendSuccess(player.Conn, "taken="+item.ID)
			return
		}
	}
	speak.SendErr(player.Conn, speak.ErrItemNotFound)
}
