package main

import "encoding/json"

func (g *Game) SendMessage(message []byte) {
	if g.connectionFail {
		return
	}
	g.sendCh <- message
}

// Runs before update
func (g *Game) HandleMessage(message []byte) error {
	action := &GenericAction{}
	err := json.Unmarshal(message, action)
	if err != nil {
		return err
	}

	actionName := action.Action
	// process message
	switch actionName {
	case "message_friend":
		action := &MessageFriendAction{}
		//put message into action's struct
		err := json.Unmarshal(message, action)
		if err != nil {
			return err
		}

	case "duck_customisation":
		action := &VisiblePlayerDataAction{}
		err := json.Unmarshal(message, action)
		if err != nil {
			return err
		}
		g.otherPlayerData[action.PlayerId] = &action.PlayerData
	case "working_state":
		action := &VisiblePlayerDataAction{}
		err := json.Unmarshal(message, action)
		if err != nil {
			return err
		}
		g.otherPlayerData[action.PlayerId].IsWorking = action.PlayerData.IsWorking
	case "new_player":
		action := &VisiblePlayerDataAction{}
		err := json.Unmarshal(message, action)
		if err != nil {
			return err
		}
		g.otherPlayerData[action.PlayerId] = &action.PlayerData
	}

	return nil
}

func returnMessage(g *Game, data any) error {
	message, err := json.Marshal(data)
	if err != nil {
		return err
	}
	g.SendMessage(message)
	return nil
}
