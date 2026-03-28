package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var speechBubble *ebiten.Image = LoadImageFromPath("assets/speech_bubble.png")
var plusSymbol *ebiten.Image = LoadImageFromPath("assets/plus.png")
var minusSymbol *ebiten.Image = LoadImageFromPath("assets/minus.png")
var inputBox *ebiten.Image = LoadImageFromPath("assets/input_box_full.png")

var minusButton *Button = &Button{
	Image: minusSymbol,
	Scale: duckScale,
}
var plusButton *Button = &Button{
	Image: plusSymbol,
	Scale: duckScale,
}

func UpdateUIScreen(g *Game) {
	if minusButton.IsHovered(g) || plusButton.IsHovered(g) {
		g.hasHover = true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if plusButton.IsHovered(g) {
			g.timerLength += 5
		}

		if minusButton.IsHovered(g) {
			g.timerLength -= 5
		}
	}
}

func drawUiScreen(g *Game, screen *ebiten.Image) {
	offsetX, offsetY := screen.Bounds().Size().X-200, 100

	minusButton.X = offsetX
	minusButton.Y = offsetY
	minusButton.Draw(screen)

	plusButton.X = offsetX + 100
	plusButton.Y = offsetY
	plusButton.Draw(screen)

	scoreText := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   32,
	}
	// TODO: center text properly
	content := fmt.Sprintln(g.timerLength)
	op1 := &text.DrawOptions{}
	op1.GeoM.Translate(float64(offsetX+(minusSymbol.Bounds().Size().X*duckScale)+16), float64(offsetY-10))
	op1.ColorScale.ScaleWithColor(color.Black)

	// textWidth, _ := text.Measure(content, scoreText, 1)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(duckScale), float64(duckScale))
	op.GeoM.Translate(float64(offsetX+(minusSymbol.Bounds().Dx())+26), float64(offsetY)-3)
	screen.DrawImage(inputBox, op)

	// DrawNineSlice

	text.Draw(screen, content, scoreText, op1)

	// op = &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(float64(duckScale), float64(duckScale))
	// op.GeoM.Translate(float64(offsetX+100), float64(offsetY))
	// screen.DrawImage(plusSymbol, op)

}
