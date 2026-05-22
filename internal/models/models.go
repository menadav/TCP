package models

import(
    "net"
)

type Player struct{
    Id      string
    Name    string
}

type Hub struct {
    Register   chan net.Conn  
    Unregister chan net.Conn  
    Broadcast  chan string
    Clients    map[net.Conn]*Player
}

func (h *Hub) Run(){
    for {
        select {
            case conn := <- h.Register:
                h.Clients[conn] = &Player{}
            case conn := <- h.Unregister:
                delete(h.Clients, conn)
            case msg := <- h.Broadcast:
                for conn, _ := range h.Clients{
                    conn.Write([]byte(msg + "\n"))
                }
        }
    }
}
