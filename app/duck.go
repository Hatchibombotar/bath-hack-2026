package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var duckImage *ebiten.Image = LoadImageFromPath("assets/duck.png")
var duckWalkFlipbook *ebiten.Image = LoadImageFromPath("assets/duck_walk.png")
var duckSittingImage *ebiten.Image = LoadImageFromPath("assets/duck_sitting.png")

type Duck struct {
	X, Y             int
	targetX, targetY int
	Game             *Game

	isHovered     bool //if mouse intersecting bounding box of image
	isSleeping    bool //true if in a sleeping state (in a study session)
	isHeld        bool //true if currently being dragged by mouse
	isWalking     bool
	isFacingRight bool //true if facing right, false if left
}

func (duck *Duck) Update() {
	g := duck.Game

	duck.isHovered = isPointInRect(
		g.cursorX, g.cursorY,
		g.duck.X, g.duck.Y,
		duckWidth*duckScale,
		duckWidth*duckScale,
	)

	if !duck.isHeld {
		if duck.isHovered && ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			duck.isHeld = true
		}
	} else {
		duck.isHeld = ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	}

	duck.isFacingRight = (duck.targetX - duck.X) < 0 
}

func (duck *Duck) Move() {
	// window, _ := GetForegroundWindowInfo()
	// duck.X = int(window.right)
	// duck.Y = int(window.top)

	//w, h := duck.Game.ScreenSize()
	//distanceToTarget := math.Sqrt(math.Pow(float64(duck.targetX-duck.X), 2) + math.Pow(float64(duck.targetY-duck.Y), 2))
	//fmt.Println(duck.isHeld)
	if duck.isHeld {
		duck.targetX = duck.Game.cursorX
		duck.targetY = duck.Game.cursorY

		xDistanceToTarget := float64(duck.targetX - duck.X)
		yDistanceToTarget := float64(duck.targetY - duck.Y)

		duck.X += int(xDistanceToTarget) / 4
		duck.Y += int(yDistanceToTarget) / 4
	}
}

func (duck *Duck) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	if duck.isFacingRight {
		op.GeoM.Scale(float64(duckScale), float64(duckScale))
	} else {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(duckImage.Bounds().Size().X), 0)
		op.GeoM.Scale(float64(duckScale), float64(duckScale))
	}
	op.GeoM.Translate(float64(duck.X), float64(duck.Y))

	// screen.DrawImage(duckImage, op)
	fmt.Println(duck.isWalking)
	if duck.isWalking {
		DrawSpriteFrame(screen, duckWalkFlipbook, 30, 30, (duck.Game.frame/7)%4, op)
	} else if duck.isHeld {
		screen.DrawImage(duckImage, op)
	} else {
		screen.DrawImage(duckSittingImage, op)
	}
	

	ebitenutil.DebugPrintAt(screen, "Stan Duck", duck.X, duck.Y-20)
}
