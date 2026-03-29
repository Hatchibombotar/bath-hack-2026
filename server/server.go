package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type ServerRoom struct {
	players []*Player
}

var room *ServerRoom = &ServerRoom{}
var playerId int = 0

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	player := &Player{
		Conn:     conn,
		PlayerId: playerId,
		visiblePlayerData: &VisiblePlayerData{
			DuckName:  "duck",
			DuckSkin:  "duck_green",
			IsWorking: false,
		},
	}
	playerId += 1
	// broadcastMessage(&VisiblePlayerDataAction{Action: "new_player", PlayerId: player.PlayerId, PlayerData: player.visiblePlayerData})
	// for _, player := range room.players {
	// 	message, err := json.Marshal(&VisiblePlayerDataAction{Action: "1", PlayerData: player.visiblePlayerData})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	err = player.Conn.WriteMessage(websocket.TextMessage, message)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	room.players = append(room.players, player)
	updateAllDuckData()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\n", message)

		handleMessage(player, message)

		// if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		// 	fmt.Println("Error writing message:", err)
		// 	break
		// }
	}
}

func updateAllDuckData() {
	for _, playerToBeSentTo := range room.players {
		for _, player := range room.players {
			if player == playerToBeSentTo {
				continue
			}
			data := &VisiblePlayerDataAction{Action: "new_player", PlayerId: player.PlayerId, PlayerData: player.visiblePlayerData}

			message, err := json.Marshal(data)
			if err != nil {
				continue
			}

			err = playerToBeSentTo.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				continue
			}
		}
	}
}

func main() {
	http.HandleFunc("/", wsHandler)

	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
