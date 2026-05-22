package network

import (
    "answer_protocol/internal/models"
    "answer_protocol/internal/constructor"
    "bufio"
    "fmt"
    "net"
    "io"
    "os"
)

type TextProcessor func(string)

func TextClient(text string) {
    fmt.Println(text)
}

func TextServer(text string) {
    fmt.Println("\n[Server]:", text)
}


func StartScanner(source io.Reader, process TextProcessor) {
    scanner := bufio.NewScanner(source)
    for scanner.Scan() {
        process(scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        fmt.Println("ERROR scanner:", err)
    }
}

func WriteFromStdin(conn net.Conn) {
    StartScanner(os.Stdin, func(text string) {
        conn.Write([]byte(text + "\n"))
    })
}

func BroadcastMessage(conn net.Conn, hub *models.Hub) TextProcessor {
    return func(msg string) {
        player := hub.Clients[conn]
        hub.Broadcast <-"[" + player.Name + "]: " + msg
    }
}

func ReadServer(conn net.Conn, process TextProcessor) {
    StartScanner(conn, process)
}

func ReadLine(conn net.Conn) string {
    scanner := bufio.NewScanner(conn)
    scanner.Scan()
    return scanner.Text()
}

func ClientAtender(conn net.Conn, hub *models.Hub) {
    defer func() {
        hub.Unregister <- conn
        conn.Close()
    }()
    name := ReadLine(conn)
    player := constructor.NewPlayer(conn.RemoteAddr().String(), name) 
    hub.Clients[conn] = player
    ReadServer(conn, BroadcastMessage(conn, hub))
}
