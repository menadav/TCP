package parse

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"fmt"
	"strings"
	"unicode"
)

func parseChat(partsChat []string, player *models.Player, h *models.Hub) {
	var msg models.Message

	scopeStr := strings.ToUpper(partsChat[0])
	text := partsChat[1]
	if len(text) > 40 {
		speak.SendErr(player.Conn, speak.ErrMessageTooLong)
		return
	}
	if strings.IndexFunc(text, unicode.IsControl) >= 0 {
		speak.SendErr(player.Conn, speak.ErrControlChars)
		return
	}
	switch scopeStr {
	case "GLOBAL":
		msg = models.Message{
			Scope:    models.ScopeGlobal,
			Filter:   "",
			Category: "GLOBAL",
			Content:  fmt.Sprintf("CHAT %s %s", player.Name, text),
		}
	case "ROOM":
		msg = models.Message{
			Scope:    models.ScopeRoom,
			Filter:   player.Room.Id,
			Category: "ROOM",
			Content:  fmt.Sprintf("CHAT %s %s", player.Name, text),
		}
	case "GROUP":
		if player.Group == "" {
			speak.SendErr(player.Conn, speak.ErrNotInGroup)
			return
		}
		msg = models.Message{
			Scope:    models.ScopeGroup,
			Filter:   player.Group,
			Category: "GROUP",
			Content:  fmt.Sprintf("CHAT %s %s", player.Name, text),
		}
	default:
		speak.SendErr(player.Conn, speak.ErrInvalidArgument)
		return
	}
	h.Broadcast <- msg
}
