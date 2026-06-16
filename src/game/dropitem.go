package game

import (
    "answer_protocol/src/models"
    "answer_protocol/src/speakserver"
    "strings"
)

func DropItem(player *models.Player, query string) {
    for i, item := range player.Inventory {
        if strings.EqualFold(item.Name, query) || item.ID == query {
            player.Inventory = append(player.Inventory[:i], player.Inventory[i+1:]...)
            if item.Hand{
                player.VoidDmg()
            }
            room := player.Room
            room.Mu.Lock()
            room.Items = append(room.Items, item)
            room.Mu.Unlock()
            speak.SendSuccess(player.Conn, "dropped="+item.ID)
            return
        }
    }
    speak.SendError(player.Conn, 404, "ITEM_NOT_FOUND")
}