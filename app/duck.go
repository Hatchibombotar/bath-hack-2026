package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Duck struct {
	Name string

	Skins  []string
	Skin   int
	Assets map[string]*ebiten.Image

	X, Y             int
	targetX, targetY int
	Game             *Game

	isHovered     bool //if mouse intersecting bounding box of image
	isSleeping    bool //true if in a sleeping state (in a study session)
	isHeld        bool //true if currently being dragged by mouse
	isWalking     bool
	isFacingRight bool //true if facing right, false if left
}

func (duck *Duck) Init() {
	duck.Assets = make(map[string]*ebiten.Image)
	duck.Skins = []string{"duck_bathHack","duck_green"}
}

func (duck *Duck) NextSkin() {
	// if duck.Assets == nil {
	// 	duck.Assets = make(map[string]*ebiten.Image)
	// 	duck.Skins = []string{"duck_bathHack","duck_green"}
	// }
	duck.Skin += 1
	duck.Skin = duck.Skin % len(duck.Skins)
	newSkin := duck.Skins[duck.Skin]
	
	duck.Assets["duck"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck.png", newSkin))
	duck.Assets["duck_walk"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_walk.png", newSkin))
	duck.Assets["duck_sitting"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_sitting.png", newSkin))
}

func (duck *Duck) Update() {
	duck.isWalking = true
	// if duck.Skins == nil {
	// 	duck.NextSkin()
	// }
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		duck.NextSkin()
	}

	g := duck.Game

	duck.isHovered = isPointInRect(
		g.cursorX, g.cursorY,
		g.duck.X, g.duck.Y,
		duckWidth*duckScale,
		duckWidth*duckScale,
	)
	duck.isFacingRight = (duck.targetX - duck.X) < 0

	if !duck.isHeld {
		if duck.isHovered && ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			duck.isHeld = true
		}
	} else {
		duck.isHeld = ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	}
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
		op.GeoM.Translate(float64(duck.Assets["duck"].Bounds().Size().X), 0)
		op.GeoM.Scale(float64(duckScale), float64(duckScale))
	}
	op.GeoM.Translate(float64(duck.X), float64(duck.Y))

	// screen.DrawImage(duckImage, op)
	//fmt.Println(duck.isWalking)
	if duck.isWalking {
		DrawSpriteFrame(screen, duck.Assets["duck_walk"], 30, 30, (duck.Game.frame/7)%4, op)
	} else if duck.isHeld {
		screen.DrawImage(duck.Assets["duck"], op)
	} else {
		screen.DrawImage(duck.Assets["duck_sitting"], op)
	}

	ebitenutil.DebugPrintAt(screen, duck.Name, duck.X, duck.Y-20)
}
