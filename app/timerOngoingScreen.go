package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var pauseSymbol *ebiten.Image = LoadImageFromPath("assets/pause_button.png")
var stopSymbol *ebiten.Image = LoadImageFromPath("assets/stop_button.png")
var inputBoxLong *ebiten.Image = LoadImageFromPath("assets/input_box_long.png")

var pauseButton *Button = &Button{
	Image: pauseSymbol,
	Scale: duckScale,
}

var stopButton *Button = &Button{
	Image: stopSymbol,
	Scale: duckScale,
}
var resumeButton *Button = &Button{
	Image: playSymbol,
	Scale: duckScale,
}

func UpdateTimerOngoingUIScreen(g *Game) {
	if pauseButton.IsHovered(g) || stopButton.IsHovered(g) || resumeButton.IsHovered(g) {
		g.hasHover = true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if pauseButton.IsHovered(g) {
			g.isTimerRunning = false
			g.timerDuration = g.timeRemainingOnTimer
		}
		if resumeButton.IsHovered(g) && !g.isTimerRunning {
			g.isTimerRunning = true
			g.timerStartTime = time.Now()
		}
		if stopButton.IsHovered(g) {
			g.State = TimerSettingsState
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
	}

	if g.isTimerRunning {
		g.timeRemainingOnTimer = time.Duration(g.timerDuration) - time.Since(g.timerStartTime)
	}
}

func DrawTimerOngoingUiScreen(g *Game, screen *ebiten.Image) {
	offsetX, offsetY := screen.Bounds().Size().X-200, 200

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(duckScale), float64(duckScale))
	op.GeoM.Translate(float64(offsetX-42), float64(offsetY)-42)
	screen.DrawImage(speechBubble, op)

	resumeButton.X = offsetX + 46 - 36
	resumeButton.Y = offsetY + 36
	resumeButton.Draw(screen)

	pauseButton.X = offsetX + 46
	pauseButton.Y = offsetY + 36
	pauseButton.Draw(screen)

	stopButton.X = offsetX + 46 + 36
	stopButton.Y = offsetY + 36
	stopButton.Draw(screen)

	scoreText := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   32,
	}
	// TODO: center text properly
	content := fmt.Sprint(int(g.timeRemainingOnTimer.Minutes()), ":")
	content += fmt.Sprintf("%02d", int(g.timeRemainingOnTimer.Seconds()-float64(60*int(g.timeRemainingOnTimer.Minutes()))))

	// textWidth, _ := text.Measure(content, scoreText, 1)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(duckScale), float64(duckScale))
	op.GeoM.Translate(float64(offsetX+(9)), float64(offsetY)-3)
	screen.DrawImage(inputBoxLong, op)

	// DrawNineSlice

	op1 := &text.DrawOptions{}
	op1.GeoM.Translate(float64(offsetX+20), float64(offsetY-10))
	op1.ColorScale.ScaleWithColor(color.Black)

	text.Draw(screen, content, scoreText, op1)

	// op = &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(float64(duckScale), float64(duckScale))
	// op.GeoM.Translate(float64(offsetX+100), float64(offsetY))
	// screen.DrawImage(plusSymbol, op)

}
