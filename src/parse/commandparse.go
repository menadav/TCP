package parse

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
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
		if command == "QUIT" {
			h.Unregister <- player
			speak.SendSuccess(player.Conn, "bye")
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
		if argument == "" {
			return
		}
		partsGroup := strings.Split(argument, " ")
		parseGroup(partsGroup, player, h)
	default:
		speak.SendError(player.Conn, 400, "Unknown command")
		return
	}
}