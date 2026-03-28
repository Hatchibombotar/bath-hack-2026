package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func UpdateTimerOngoingUIScreen(g *Game) {
	if minusButton.IsHovered(g) || plusButton.IsHovered(g) || playButton.IsHovered(g) {
		g.hasHover = true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if plusButton.IsHovered(g) {
			g.timerLength += 5
		}

		if minusButton.IsHovered(g) {
			g.timerLength -= 5
		}

		g.timerLength = int(math.Max(5, float64(g.timerLength)))
		g.timerLength = int(math.Min(95, float64(g.timerLength)))
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if playButton.IsHovered(g) {
			g.State = TimerOngoingState
			fmt.Println("eek")
		}
	}
}

func DrawTimerOngoingUiScreen(g *Game, screen *ebiten.Image) {
	offsetX, offsetY := screen.Bounds().Size().X-200, 200

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(duckScale), float64(duckScale))
	op.GeoM.Translate(float64(offsetX-42), float64(offsetY)-42)
	screen.DrawImage(speechBubble, op)

	minusButton.X = offsetX
	minusButton.Y = offsetY
	minusButton.Draw(screen)

	plusButton.X = offsetX + 96
	plusButton.Y = offsetY
	plusButton.Draw(screen)

	playButton.X = offsetX + 46
	playButton.Y = offsetY + 36
	playButton.Draw(screen)

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

	op = &ebiten.DrawImageOptions{}
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
