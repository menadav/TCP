package network

import (
	"answer_protocol/src/models"
	"answer_protocol/src/speakserver"
	"answer_protocol/src/utils"
	"bufio"
	"net"
	"regexp"
	"strings"
	"time"
)

func Authentication(scanner *bufio.Scanner, conn net.Conn, h *models.Hub) string {
	for {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if !scanner.Scan() {
			err := scanner.Err()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					speak.SendErr(conn, speak.ErrTimeout)
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
			speak.SendErr(conn, speak.ErrMalformedCommand)
			continue
		}
		name_clean := strings.ToUpper(name_list[0])
		if name_clean == "CONNECT" {
			name_p := name_list[1]
			if len(name_p) > 12 {
				speak.SendErr(conn, speak.ErrNameTooLong)
				continue
			} else if len(name_p) < 3 {
				speak.SendErr(conn, speak.ErrNameTooShort)
				continue
			} else if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(name_p) {
				speak.SendErr(conn, speak.ErrNameInvalid)
				continue
			} else if utils.ExistName(h.Clients, name_p) {
				speak.SendErr(conn, speak.ErrNameInUse)
				continue
			}
			conn.SetReadDeadline(time.Time{})
			speak.SendSuccess(conn, "connected")
			return name_p
		} else {
			speak.SendErr(conn, speak.ErrMalformedCommand)
		}
	}
}
