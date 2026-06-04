package parse

import(
	"answer_protocol/src/speakserver"
	"answer_protocol/src/models"
	"fmt"
	"strings"
)

func parseChat(partsChat []string, player *models.Player, h *models.Hub) {
	var msg models.Message

	scopeStr := strings.ToUpper(partsChat[0])
	text := partsChat[1]
	if len(text) > 40 {
		speak.SendError(player.Conn, 203, "TEXT_TO_LONG")
		return
	}
	switch scopeStr {
	case "GLOBAL":
		msg = models.Message{
			Scope:   models.ScopeGlobal,
			Filter:  "",
			Category: "GLOBAL",
			Content: fmt.Sprintf("CHAT %s %s", player.Name, text),
		}
	case "ROOM":
		msg = models.Message{
			Scope:   models.ScopeRoom,
			Filter:  player.Room.Id,
			Category: "ROOM",
			Content: fmt.Sprintf("CHAT %s %s", player.Name, text),
		}
	case "GROUP":
		if player.Group == "" {
			speak.SendError(player.Conn, 403, "You are not in a group")
			return
		}
		msg = models.Message{
			Scope:    models.ScopeGroup,
			Filter:   player.Group,
			Category: "GROUP",
			Content:  fmt.Sprintf("CHAT %s %s", player.Name, text),
		}
	default:
		speak.SendError(player.Conn, 400, "Unknown chat scope. Use GLOBAL, ROOM, or GROUP")
		return
	}
	h.Broadcast <- msg
}
