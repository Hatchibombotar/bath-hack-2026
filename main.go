package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var pixelScale = 1
var duckScale = 3
var duckWidth = 22

type Game struct {
	duck             *Duck
	cursorX, cursorY int
	isUiOpen         bool
}

func (g *Game) Update() error {
	g.cursorX, g.cursorY = ebiten.CursorPosition()
	g.duck.Update()

	if g.duck.isHovered {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			fmt.Println("Quack")
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			fmt.Println("Quack")
		}
		ebiten.SetWindowMousePassthrough(false)
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetWindowMousePassthrough(true)
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		screen.Fill(color.RGBA{255, 255, 255, 10})
	}
	g.duck.Draw(screen)
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := g.ScreenSize()
	return w, h
}

func (g *Game) ScreenSize() (int, int) {
	w, h := ebiten.Monitor().Size()
	return (w / pixelScale) - 1, (h / pixelScale)
}

func main() {
	w, h := ebiten.Monitor().Size()
	// window size 1
	ebiten.SetWindowSize((w)-1, h)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowPosition(0, 0)

	game := &Game{}

	duck := &Duck{
		Game: game,
		X:    100,
		Y:    100,
	}

	game.duck = duck

	err := ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{
		ScreenTransparent: true,
	})

	if err != nil {
		log.Fatal(err)
	}
}
