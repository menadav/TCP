package network

import (
    "answer_protocol/src/models"
    "answer_protocol/src/constructor"
    "answer_protocol/src/parse"
    "answer_protocol/src/speakserver"
    "bufio"
    "net"
    "strings"
    "regexp"
    "time"
)

func Authentication(scanner *bufio.Scanner, conn net.Conn) string {
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
        if len(name_list) != 2 {
            speak.SendError(conn, 400, "malformed command CONNECT <name>")
            continue
        }
        name_clean := strings.ToUpper(name_list[0])
        if  name_clean == "CONNECT" {
            name_p := name_list[1]
            if len(name_p) > 12 {
                speak.SendError(conn, 400, "max 12 characters")
                continue
            } else if len(name_p) < 3{
                speak.SendError(conn, 400, "min 3 characters")
                continue
            } else if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(name_p){
                speak.SendError(conn, 400, "only letters")
                continue
            }
            conn.SetReadDeadline(time.Time{})
            speak.SendSuccess(conn, "connected")
            return name_p
        } else {
            speak.SendError(conn, 400, "malformed_command")
        }
    }
}

func ClientAtender(conn net.Conn, hub *models.Hub) {
    scanner := bufio.NewScanner(conn)
    name_p := ""
    name_p = Authentication(scanner, conn)
    if name_p == "" {
        conn.Close()
        return
    }
    player := constructor.NewPlayer(conn.RemoteAddr().String(), conn, name_p, "loc.start")
    defer func() {
        if name_p != "" {
            hub.Unregister <- player
        }
        conn.Close()
    }()
    hub.Register <- player
    processFunction := func(line string) {
		parse.ParseCommandCli(line, player, hub)
	}
    ReadServer(scanner, processFunction)
}
