package main

import "encoding/json"

func (g *Game) SendMessage(message []byte) {
	if g.connectionFail {
		return
	}
	g.sendCh <- message
}

func (g *Game) SendDuckInfo() {
	data := &VisiblePlayerDataAction{Action: "update_info", PlayerData: &VisiblePlayerData{
		DuckName:  g.duck.Name,
		DuckSkin:  g.duck.GetSkin(),
		IsWorking: g.duck.isSleeping,
	}}

	message, err := json.Marshal(data)
	if err != nil {
		return
	}

	g.SendMessage(message)
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
	case "player_update":
		action := &VisiblePlayerDataAction{}
		err := json.Unmarshal(message, action)
		if err != nil {
			return err
		}
		g.otherPlayerData[action.PlayerId] = action.PlayerData

		playerIndex := indexOfPlayer(g, action.PlayerId)
		nestX, nestY := getNestPosition(1+playerIndex, g)

		_, exists := g.otherPlayers[action.PlayerId]
		if !exists {
			g.otherPlayers[action.PlayerId] = &Duck{
				isOtherDuck: true,
				X:           float64(nestX) + 6,
				Y:           float64(nestY),
				nestX:       nestX + 6,
				nestY:       nestY,
				Game:        g,
			}
			g.otherPlayers[action.PlayerId].Init()
		}

		g.otherPlayers[action.PlayerId].Name = action.PlayerData.DuckName
		g.otherPlayers[action.PlayerId].SetSkin(action.PlayerData.DuckSkin)
		g.otherPlayers[action.PlayerId].isSleeping = action.PlayerData.IsWorking

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
