package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

func handleMessage(player *Player, message []byte) error {
	action := &GenericAction{}
	err := json.Unmarshal(message, action)
	if err != nil {
		return err
	}

	actionName := action.Action

	switch actionName {
	case "message_friend":
		action := &MessageFriendAction{}
		err := json.Unmarshal(message, action)
		if err != nil {
			return err
		}
		broadcastMessage(action)
	case "duck_customisation":
		action := &VisiblePlayerDataAction{}
		err = json.Unmarshal(message, action)
		if err != nil {
			return err
		}
		action.PlayerId = player.PlayerId
		player.visiblePlayerData.DuckName = action.PlayerData.DuckName
		player.visiblePlayerData.DuckSkin = action.PlayerData.DuckSkin
		broadcastMessage(action)
	case "working_state":
		action := &VisiblePlayerDataAction{}
		err = json.Unmarshal(message, action)
		if err != nil {
			return err
		}
		action.PlayerId = player.PlayerId
		player.visiblePlayerData.IsWorking = action.PlayerData.IsWorking
		broadcastMessage(action)
	}

	return nil
}

func broadcastMessage(data any) error {
	for _, player := range room.players {
		message, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = player.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return err
		}
	}
	return nil
}
