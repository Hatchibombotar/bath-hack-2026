package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var pixelScale = 1
var duckScale = 3
var duckWidth = 22

var nestSpriteBack *ebiten.Image = LoadImageFromPath("assets/nest_back.png")
var nestSpriteFront *ebiten.Image = LoadImageFromPath("assets/nest_front.png")
var sittingAssets map[string]*ebiten.Image = map[string]*ebiten.Image{
	"duck_bathHack": LoadImageFromPath(fmt.Sprintf("assets/%s/duck_sitting.png", "duck_bathHack")),
	"duck_green":    LoadImageFromPath(fmt.Sprintf("assets/%s/duck_sitting.png", "duck_green")),
}

type ServerState int

const (
	TimerSettingsState ServerState = iota
	TimerOngoingState
)

type VisiblePlayerData struct {
	DuckName  string `json:"duckname"`
	DuckSkin  string `json:"duckskin"`
	IsWorking bool   `json:"isworking"`
}

type Game struct {
	msgCh  chan []byte
	sendCh chan []byte
	ctx    context.Context
	cancel context.CancelFunc

	duck             *Duck
	cursorX, cursorY int
	hasHover         bool
	// isActionUiOpen   bool
	// actionUI         *ActionUI
	frame int

	timeRemainingOnTimer time.Duration
	timerStartTime       time.Time
	isTimerRunning       bool

	isStartingUIOpen bool
	timerLength      int
	timerDuration    time.Duration

	connectionFail bool

	State ServerState

	otherPlayerData map[int]*VisiblePlayerData
	otherPlayers    map[int]*Duck
	playerList      []int

	showUi bool
}

func indexOfPlayer(g *Game, id int) int {
	for i, otherId := range g.playerList {
		if id == otherId {
			return i
		}
	}
	g.playerList = append(g.playerList, id)
	for i, otherId := range g.playerList {
		if id == otherId {
			return i
		}
	}
	panic("bad")
}

// Update processes incoming websocket messages (non-blocking).
func (g *Game) Update() error {
	for playerId, player := range g.otherPlayerData {
		_, exists := g.otherPlayers[playerId]
		if !exists {
			panic("no duck!")
		}
		g.otherPlayers[playerId].Name = player.DuckName
		g.otherPlayers[playerId].SetSkin(player.DuckSkin)

	}
	select {
	case m := <-g.msgCh:
		// handle message: update game state based on m
		log.Printf("client got ws message: %s", m)
		g.HandleMessage(m)
	default:
		// no message this frame
	}
	g.frame += 1
	g.cursorX, g.cursorY = ebiten.CursorPosition()
	g.hasHover = false

	if g.showUi {
		switch g.State {
		case TimerSettingsState:
			UpdateTimerInputUIScreen(g)
		case TimerOngoingState:
			UpdateTimerOngoingUIScreen(g)
		}
	}

	g.duck.Update()
	g.duck.Move()

	for _, duck := range g.otherPlayers {
		duck.Update()
		duck.Move()
	}

	if g.duck.isHovered {
		g.hasHover = true
	}
	if g.duck.isHovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		g.showUi = !g.showUi
	}
	// 	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
	// 		fmt.Println("Quack")
	// 		g.isActionUiOpen = true

	// 		actionUI := &ActionUI{}
	// 		actionUI.X = g.duck.X + 20
	// 		actionUI.Y = g.duck.Y - 30
	// 		actionUI.Make()

	// 		g.actionUI = actionUI
	// 	}
	// }

	// if g.actionUI != nil {
	// 	g.actionUI.Update(g)
	// }

	if g.hasHover {
		ebiten.SetWindowMousePassthrough(false)
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetWindowMousePassthrough(true)
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	return nil
}

func getNestPosition(i int, g *Game) (int, int) {
	w, h := g.ScreenSize()
	x := float64(w) - float64(nestSpriteBack.Bounds().Size().X*duckScale) - float64(nestSpriteBack.Bounds().Size().X*duckScale*i)
	y := float64(h) - float64(nestSpriteBack.Bounds().Size().Y*duckScale) - 42

	return int(x), int(y)
}

// Draw renders a simple message count for demonstration.
func (g *Game) Draw(screen *ebiten.Image) {
	// if ebiten.IsKeyPressed(ebiten.KeySpace) {
	// 	screen.Fill(color.RGBA{255, 255, 255, 10})
	// }

	i := 0
	for range len(g.otherPlayerData) + 1 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(duckScale), float64(duckScale))
		X, Y := getNestPosition(i, g)
		op.GeoM.Translate(
			float64(X), float64(Y),
		)
		screen.DrawImage(nestSpriteBack, op)
		i++
	}

	g.duck.Draw(screen)

	for _, duck := range g.otherPlayers {
		duck.Draw(screen)
	}

	i = 0
	for range len(g.otherPlayerData) + 1 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(duckScale), float64(duckScale))
		X, Y := getNestPosition(i, g)
		op.GeoM.Translate(
			float64(X), float64(Y),
		)
		screen.DrawImage(nestSpriteFront, op)
		i++
	}

	// op = &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(float64(screen.Bounds().Dx())-float64(speechBubble.Bounds().Dx()), float64(screen.Bounds().Dy())-float64(speechBubble.Bounds().Dx()))
	// screen.DrawImage(speechBubble, op)

	if g.showUi {
		switch g.State {
		case TimerSettingsState:
			DrawTimerInputUiScreen(g, screen)
		case TimerOngoingState:
			DrawTimerOngoingUiScreen(g, screen)
		}
	}
}

// Layout returns the screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := g.ScreenSize()
	return w, h
}

func (g *Game) ScreenSize() (int, int) {
	w, h := ebiten.Monitor().Size()
	return (w / pixelScale) - 1, (h / pixelScale)
}

func main() {
	// channel to receive text messages from websocket reader
	msgCh := make(chan []byte, 64)
	sendCh := make(chan []byte, 8)

	ctx, cancel := context.WithCancel(context.Background())
	game := &Game{msgCh: msgCh, ctx: ctx, cancel: cancel, sendCh: sendCh}

	game.timerLength = 60
	game.isStartingUIOpen = true
	game.otherPlayerData = make(map[int]*VisiblePlayerData)
	game.otherPlayers = make(map[int]*Duck)

	// start websocket client goroutine
	go runWebSocketClient(ctx, "ws://localhost:8080/", msgCh, sendCh, game)

	// start ebiten main loop
	w, h := ebiten.Monitor().Size()
	// window size 1
	ebiten.SetWindowSize((w)-1, h)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowPosition(0, 0)

	nestX, nestY := getNestPosition(0, game)
	duck := &Duck{
		Game:  game,
		X:     float64(nestX) + 6,
		Y:     float64(nestY),
		nestX: nestX + 6,
		nestY: nestY,
	}

	duck.Init()

	game.duck = duck
	win, title := GetForegroundWindowInfo()
	println(win.bottom, title)
	err := ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{
		ScreenTransparent: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// cancel WS goroutine on exit
	cancel()
	// allow graceful shutdown
	time.Sleep(200 * time.Millisecond)
}
func runWebSocketClient(ctx context.Context, rawURL string, msgCh chan<- []byte, sendCh <-chan []byte, g *Game) {
	u, _ := url.Parse(rawURL)
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		g.connectionFail = true
		log.Printf("ws dial error: %v", err)
		return
	}
	defer conn.Close()

	conn.SetReadLimit(1024 * 1024)
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	pingTicker := time.NewTicker(10 * time.Second)
	defer pingTicker.Stop()

	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("ws read error: %v", err)
				return
			}
			if mt == websocket.TextMessage {
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
			select {
			case <-readDone:
			case <-time.After(1 * time.Second):
			}
			return
		case <-pingTicker.C:
			_ = conn.WriteMessage(websocket.PingMessage, nil)
		case m := <-sendCh:
			// send text message; handle write errors
			if err := conn.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
				// TODO: UNCOMMENT
				fmt.Println("tried to do a thing")
				// log.Printf("ws write error: %v", err)
				// return
			}
		}
	}
}
