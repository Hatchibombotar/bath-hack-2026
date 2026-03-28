package main

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var duckImage *ebiten.Image = LoadImageFromPath("assets/duck.png")

type Duck struct {
	X, Y             int
	targetX, targetY int
	Game             *Game

	isHovered bool
}

func (duck *Duck) Update() {
	g := duck.Game
	w, h := duck.Game.ScreenSize()
	distanceToTarget := math.Sqrt(math.Pow(float64(duck.targetX-duck.X), 2) + math.Pow(float64(duck.targetY-duck.Y), 2))
	if distanceToTarget < 3 {
		duck.targetX = int(float64(w) * rand.Float64())
		duck.targetY = int(float64(h) * rand.Float64())
	}

	duck.isHovered = isPointInRect(
		g.cursorX, g.cursorY,
		g.duck.X, g.duck.Y,
		duckWidth*duckScale,
		duckWidth*duckScale,
	)

	if !duck.isHovered {
		if duck.targetX > duck.X {
			duck.X += 1
		} else if duck.targetX < duck.X {
			duck.X -= 1
		}
		if duck.targetY > duck.Y {
			duck.Y += 1
		} else if duck.targetY < duck.Y {
			duck.Y -= 1
		}
	}
}
func (duck *Duck) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(duckScale), float64(duckScale))
	op.GeoM.Translate(float64(duck.X), float64(duck.Y))
	screen.DrawImage(duckImage, op)

	ebitenutil.DebugPrintAt(screen, "Stan Duck", duck.X, duck.Y-20)
}
