package main

import "github.com/gorilla/websocket"

type Player struct {
	Conn              *websocket.Conn
	PlayerId          int
	visiblePlayerData VisiblePlayerData
}
