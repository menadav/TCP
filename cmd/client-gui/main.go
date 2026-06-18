package main

import (
	"answer_protocol/src/gui"
	"answer_protocol/src/models"
	"net"
	"fmt"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Error Dial", err)
		return
	}
	defer conn.Close()

	player := &models.Player{
		Conn: conn,
		Name: "Player1",
	}

	gui.Start(a, player)
}