package parse

import (
	"answer_protocol/src/network"
	"strings"
	"net"
)

func parseCommandCli(line string, player Player, h *Hub) {
	line = strings.TrimSpace(line)
	if line == ""{
		return
	}
	parts := line.SplitN(line, " ", 2)
	command := strings.ToUpper(parts[0])

	var argument string
	if len(parts) > 1 {
		argument = parts[1]
	}
	switch command {
		case "LOOK", "INVENTORY", "STATUS", "QUESTS", "WHO", "QUIT":
			if argument != "" {
				network.SendError(conn, 400, "Only command")
				return
			}
		case "MOVE":
			if argument == "" {
				network.SendError(conn, 400, "Only command")
				return
			}
		case "CHAT":
			partsChat := stringsSplitN(parts, " ", 2)
			if len(partsChat) < 2 {
				network.SendError(conn, 400, "Only command")
				return
			}
			parseChat(partsChat, player, hub)
		case "TAKE", "DROP", "TALK", "ATTACK", "QUEST":
			if argument == "" {
				network.SendError(conn, 400, "Only command")
				return
			}
		case "GROUP":
			if argument != "" {
				return
			}
		default:
			return
	}
}


func parseChat(partsChat []string, player Player, h *Hub) {
	var msg Message

	scopeStr := strings.ToUpper(partsChat[0])
	text := partsChat[1]
    switch scopeStr {
    case "GLOBAL":
        msg = Message{
            Scope:   ScopeGlobal,
            Filter:  "",
            Content: fmt.Sprintf("GLOBAL CHAT %s %s\n", player.Name, text),
        }
    case "ROOM":
        msg = Message{
            Scope:   ScopeRoom,
            Filter:  player.Room,
            Content: fmt.Sprintf("ROOM CHAT %s %s\n", player.Name, text),
        }
    case "GROUP":
        if player.Group == "" {
            network.SendError(player.Conn, 403, "You are not in a group")
            return
        }
        msg = Message{
            Scope:   ScopeGroup,
            Filter:  player.Group,
            Content: fmt.Sprintf("GROUP CHAT %s %s\n", player.Name, text),
        }
    default:
        network.SendError(player.Conn, 400, "Unknown chat scope. Use GLOBAL, ROOM, or GROUP")
        return
    }
    h.Broadcast <- msg
}