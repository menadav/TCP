package network

import(
	"answer_protocol/src/models"
    "bufio"
    "net"
	"answer_protocol/src/speakserver"
    "answer_protocol/src/utils"
    "strings"
    "regexp"
    "time"
)

func Authentication(scanner *bufio.Scanner, conn net.Conn, h *models.Hub) string {
    for {
        conn.SetReadDeadline(time.Now().Add(30 * time.Second))
        if !scanner.Scan() {
            err := scanner.Err()
            if err != nil {
                if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                    speak.SendError(conn, 408, "CONNECTION_TIMEOUT")
                }
            }
            return ""
        }
        name := scanner.Text()
        name_list := strings.Fields(name)
        if name == "REQ" {
            continue
        }
        if len(name_list) != 2 {
            speak.SendError(conn, 400, "MALFORMED_COMMAND_CONNECT <name>")
            continue
        }
        name_clean := strings.ToUpper(name_list[0])
        if  name_clean == "CONNECT" {
            name_p := name_list[1]
            if len(name_p) > 12 {
                speak.SendError(conn, 203, "MAX_12_CHARACTERS")
                continue
            } else if len(name_p) < 3{
                speak.SendError(conn, 400, "MIN_3_CHARACTERS")
                continue
            } else if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(name_p){
                speak.SendError(conn, 202, "ONLY_LETTERS")
                continue
            } else if utils.ExistName(h.Clients, name_p) {
                speak.SendError(conn, 201, "NAME_IN_USE")
                continue
            }
            conn.SetReadDeadline(time.Time{})
            speak.SendSuccess(conn, "connected")
            return name_p
        } else {
            speak.SendError(conn, 202, "MALFORMED_COMMAND")
        }
    }
}