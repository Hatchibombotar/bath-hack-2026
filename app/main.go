package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
)

var pixelScale = 1
var duckScale = 3
var duckWidth = 22

type ServerState int

const (
	WaitingToStart ServerState = iota
	StateRetrying
)

type Game struct {
	msgCh  chan string
	sendCh chan []byte
	ctx    context.Context
	cancel context.CancelFunc

	duck             *Duck
	cursorX, cursorY int
	hasHover         bool
	// isActionUiOpen   bool
	// actionUI         *ActionUI
	frame int

	isStartingUIOpen bool
	timerLength      int

	connectionFail bool
}

// Update processes incoming websocket messages (non-blocking).
func (g *Game) Update() error {
	select {
	case m := <-g.msgCh:
		// handle message: update game state based on m
		log.Printf("client got ws message: %s", m)
	default:
		// no message this frame
	}

	g.frame += 1
	g.cursorX, g.cursorY = ebiten.CursorPosition()
	g.hasHover = false

	UpdateUIScreen(g)
	g.duck.Update()
	g.duck.Move()

	if g.duck.isHovered {
		g.hasHover = true

		message, err := json.Marshal(&GenericAction{Action: "blah blah"})
		if err != nil {
			panic(err)
		}

		g.SendMessage(message)
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

// Draw renders a simple message count for demonstration.
func (g *Game) Draw(screen *ebiten.Image) {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		screen.Fill(color.RGBA{255, 255, 255, 10})
	}
	g.duck.Draw(screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(duckScale), float64(duckScale))
	op.GeoM.Translate(100, 100)

	// op = &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(float64(screen.Bounds().Dx())-float64(speechBubble.Bounds().Dx()), float64(screen.Bounds().Dy())-float64(speechBubble.Bounds().Dx()))
	// screen.DrawImage(speechBubble, op)

	drawUiScreen(g, screen)
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
	msgCh := make(chan string, 64)
	sendCh := make(chan []byte, 8)

	ctx, cancel := context.WithCancel(context.Background())
	game := &Game{msgCh: msgCh, ctx: ctx, cancel: cancel, sendCh: sendCh}

	game.timerLength = 60
	game.isStartingUIOpen = true

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

	duck := &Duck{
		Game: game,
		X:    200,
		Y:    200,
	}

	game.duck = duck

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
func runWebSocketClient(ctx context.Context, rawURL string, msgCh chan<- string, sendCh <-chan []byte, g *Game) {
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
				case msgCh <- string(msg):
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
