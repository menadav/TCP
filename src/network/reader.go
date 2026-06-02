package network

import (
    "answer_protocol/src/models"
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

type TextProcessor func(string)

func TextClient(text string) {
    fmt.Println(text)
}

func StartScanner(scanner *bufio.Scanner, process TextProcessor) {
    for scanner.Scan() {
        process(scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        if strings.Contains(err.Error(), "use of closed network connection") {
            return
        }
        fmt.Println("ERROR scanner:", err)
    }
}

func WriteFromStdin(conn net.Conn) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        conn.Write([]byte(scanner.Text() + "\n"))
    }
}

func BroadcastMessage(name string, hub *models.Hub) TextProcessor {
    return func(msg string) {
        new_name := fmt.Sprintf("[%s]", name)
        hub.Broadcast <- models.Message{
			Scope:   models.ScopeGlobal,
			Content: fmt.Sprintf("EVT GLOBAL CHAT %s %s\n", new_name, msg),
		}
    }
}

func ReadServer(scanner *bufio.Scanner, process TextProcessor) {
    StartScanner(scanner, process)
}

func ReadLine(conn net.Conn) string {
    scanner := bufio.NewScanner(conn)
    scanner.Scan()
    return scanner.Text()
}
