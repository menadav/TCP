package parse

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"fmt"
	"strings"
)

func ParseCommandCli(line string, player *models.Player, h *models.Hub) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	parts := strings.SplitN(line, " ", 2)
	command := strings.ToUpper(parts[0])
	var argument string
	if len(parts) > 1 {
		argument = parts[1]
	}
	switch command {
	case "LOOK", "INVENTORY", "STATUS", "QUESTS", "WHO", "QUIT":
		if argument != "" {
			speak.SendError(player.Conn, 400, "Only command, no arguments allowed")
			return
		}
	case "MOVE":
		if argument == "" {
			speak.SendError(player.Conn, 400, "Move requires a destination")
			return
		}
	case "CHAT":
		if argument == "" {
			speak.SendError(player.Conn, 400, "Chat requires a scope and a message")
			return
		}
		partsChat := strings.SplitN(argument, " ", 2)
		if len(partsChat) < 2 {
			speak.SendError(player.Conn, 400, "Chat format invalid. Use: CHAT <SCOPE> <MESSAGE>")
			return
		}
		parseChat(partsChat, player, h)
	case "TAKE", "DROP", "TALK", "ATTACK", "QUEST":
		if argument == "" {
			speak.SendError(player.Conn, 400, "This command requires an argument")
			return
		}
	case "GROUP":
		if argument != "" {
			return
		}
	default:
		speak.SendError(player.Conn, 400, "Unknown command")
		return
	}
}

func parseChat(partsChat []string, player *models.Player, h *models.Hub) {
	var msg models.Message

	scopeStr := strings.ToUpper(partsChat[0])
	text := partsChat[1]
	switch scopeStr {
	case "GLOBAL":
		msg = models.Message{
			Scope:   models.ScopeGlobal,
			Filter:  "",
			Content: fmt.Sprintf("GLOBAL CHAT %s %s\n", player.Name, text),
		}
	case "ROOM":
		msg = models.Message{
			Scope:   models.ScopeRoom,
			Filter:  player.Room,
			Content: fmt.Sprintf("ROOM CHAT %s %s\n", player.Name, text),
		}
	case "GROUP":
		if player.Group == "" {
			speak.SendError(player.Conn, 403, "You are not in a group")
			return
		}
		msg = models.Message{
			Scope:    models.ScopeGroup,
			Filter:   player.Group,
			Category: "CHAT",
			Content:  fmt.Sprintf("GROUP CHAT %s %s\n", player.Name, text),
		}
	default:
		speak.SendError(player.Conn, 400, "Unknown chat scope. Use GLOBAL, ROOM, or GROUP")
		return
	}
	h.Broadcast <- msg
}