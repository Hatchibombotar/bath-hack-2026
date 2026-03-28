package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var pixelScale = 1
var duckScale = 3
var duckWidth = 22

// Game holds game state and a channel for incoming WS messages.
type Game struct {
	msgCh  chan string
	ctx    context.Context
	cancel context.CancelFunc

	duck             *Duck
	cursorX, cursorY int
	hasHover         bool
	isActionUiOpen   bool
	actionUI         *ActionUI
	frame            int
}

// Update processes incoming websocket messages (non-blocking).
func (g *Game) Update() error {
	select {
	case m := <-g.msgCh:
		// handle message: update game state based on m
		log.Printf("got ws message: %s", m)
	default:
		// no message this frame
	}

	g.frame += 1
	g.cursorX, g.cursorY = ebiten.CursorPosition()
	g.hasHover = false

	g.duck.Update()
	if !g.isActionUiOpen {
		g.duck.Move()
	}

	if g.duck.isHovered {
		g.hasHover = true
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
			fmt.Println("Quack")
			g.isActionUiOpen = true

			actionUI := &ActionUI{}
			actionUI.X = g.duck.X + 20
			actionUI.Y = g.duck.Y - 30
			actionUI.Make()

			g.actionUI = actionUI
		}
	}

	if g.actionUI != nil {
		g.actionUI.Update(g)
	}

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
	ebitenutil.DebugPrint(screen, "WebSocket messages received (check logs)")
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		screen.Fill(color.RGBA{255, 255, 255, 10})
	}
	g.duck.Draw(screen)
	ebitenutil.DebugPrint(screen, "Hello, World!")

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(duckScale), float64(duckScale))
	op.GeoM.Translate(100, 100)

	if g.isActionUiOpen {
		g.actionUI.Draw(screen)
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
	msgCh := make(chan string, 64)

	ctx, cancel := context.WithCancel(context.Background())
	game := &Game{msgCh: msgCh, ctx: ctx, cancel: cancel}

	// start websocket client goroutine
	go runWebSocketClient(ctx, "ws://localhost:8080/", msgCh)

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

// runWebSocketClient connects and reads messages, forwarding text messages to msgCh.
func runWebSocketClient(ctx context.Context, rawURL string, msgCh chan<- string) {
	u, _ := url.Parse(rawURL)
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("ws dial error: %v", err)
		return
	}
	defer conn.Close()

	// ensure read deadlines and pong handler if desired
	conn.SetReadLimit(1024 * 1024)
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// writer: periodic ping
	pingTicker := time.NewTicker(10 * time.Second)
	defer pingTicker.Stop()

	readDone := make(chan struct{})
	// reader goroutine
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
					// drop if channel full
				}
			} else {
				// handle other types if needed
			}
		}
	}()

	// main loop: send pings and watch for context cancel
	for {
		select {
		case <-ctx.Done():
			// close connection gracefully
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
			// wait for reader to finish
			select {
			case <-readDone:
			case <-time.After(1 * time.Second):
			}
			return
		case <-pingTicker.C:
			_ = conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}
