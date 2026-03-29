package main

import (
	"encoding/json"
)

func handleMessage(player *Player, message []byte) error {
	action := &GenericAction{}
	err := json.Unmarshal(message, action)
	if err != nil {
		return err
	}

	actionName := action.Action

	switch actionName {
	case "player_update":
		action := &VisiblePlayerDataAction{}
		err = json.Unmarshal(message, action)
		if err != nil {
			return err
		}
		player.visiblePlayerData = action.PlayerData
		updateAllDuckData()
	}

	return nil
}

// func broadcastMessage(data any) error {
// 	for _, player := range room.players {
// 		message, err := json.Marshal(data)
// 		if err != nil {
// 			return err
// 		}

// 		err = player.Conn.WriteMessage(websocket.TextMessage, message)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
