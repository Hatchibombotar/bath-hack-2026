package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ActionUI struct {
	X, Y    int
	Button0 *Button
}

var actionButtonSprite *ebiten.Image = LoadImageFromPath("assets/test_button.png")

func (ui *ActionUI) Make() {
	actionButton0 := &Button{
		Image: actionButtonSprite,
		Scale: duckScale,
		X:     ui.X - 50,
		Y:     ui.Y,
	}
	// actionButton1 := &Button{
	// 	Image: actionButtonSprite,
	// 	Scale: duckScale,
	// 	X:     ui.X,
	// 	Y:     ui.Y,
	// }

	ui.Button0 = actionButton0
}

func (ui *ActionUI) Update(g *Game) {
	if ui.Button0.IsHovered(g) {
		g.hasHover = true

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			fmt.Println("action 1")
		}
	}
}
func (ui *ActionUI) Draw(screen *ebiten.Image) {
	ui.Button0.Draw(screen)
}
