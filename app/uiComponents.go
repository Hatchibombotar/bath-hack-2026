package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	Image *ebiten.Image
	X, Y  int
	Scale int
}

func (b *Button) IsHovered(g *Game) bool {
	return isPointInRect(
		g.cursorX, g.cursorY,
		b.X, b.Y,
		b.Image.Bounds().Dx()*b.Scale,
		b.Image.Bounds().Dy()*b.Scale,
	)
}

func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(b.Scale), float64(b.Scale))
	op.GeoM.Translate(float64(b.X), float64(b.Y))

	screen.DrawImage(b.Image, op)
}
