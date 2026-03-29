package main

import (
	"fmt"
	"math/rand/v2"

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
	nestX, nestY     int

	isHovered     bool //if mouse intersecting bounding box of image
	isSleeping    bool //true if in a sleeping state (in a study session)
	isHeld        bool //true if currently being dragged by mouse
	isWalking     bool
	isFlying      bool
	isFacingRight bool //true if facing right, false if left
	takeBreak     bool

	lastTimestamp int
	waitTime      int
}

func (duck *Duck) Init() {
	duck.Assets = make(map[string]*ebiten.Image)
	duck.Skins = []string{"duck_bathHack", "duck_green"}
	duck.Name = "big stan"
	duck.NextSkin()
}

func (duck *Duck) NextSkin() {
	duck.Skin += 1
	duck.Skin = duck.Skin % len(duck.Skins)
	newSkin := duck.Skins[duck.Skin]

	duck.Assets["duck"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck.png", newSkin))
	duck.Assets["duck_walk"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_walk.png", newSkin))
	duck.Assets["duck_sitting"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_sitting.png", newSkin))
	duck.Assets["duck_sleeping"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_sleeping.png", newSkin))
}

func (duck *Duck) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		duck.NextSkin()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		duck.isSleeping = !duck.isSleeping
	}

	g := duck.Game

	duck.isHovered = isPointInRect(
		g.cursorX, g.cursorY,
		g.duck.X, g.duck.Y,
		duckWidth*duckScale,
		duckWidth*duckScale,
	)

	if duck.isSleeping {
		duck.isWalking = false
		duck.isFlying = false
	} else if duck.isHeld {
		duck.isFacingRight = (duck.targetX - duck.X) < 0
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			duck.isHeld = false
			duck.lastTimestamp = g.frame
			duck.waitTime = rand.IntN(200) + 50
		}
	} else if duck.isHovered && ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		duck.isHeld = true
		//could make array of motion booleans that can be set to false via one func call to prevent repetition
		duck.isWalking = false
		duck.isFlying = false
	} else if g.frame-duck.lastTimestamp > duck.waitTime {
		duck.takeBreak = !duck.takeBreak
		if duck.takeBreak {
			duck.isWalking = false
			duck.isFlying = false
		}
		duck.lastTimestamp = g.frame
		duck.waitTime = rand.IntN(400) + 100
	}
}

func (duck *Duck) Move() {
	if duck.isHeld {
		duck.targetX = duck.Game.cursorX - (int(duck.Assets["duck"].Bounds().Size().X)/2)*duckScale
		duck.targetY = duck.Game.cursorY - (int(duck.Assets["duck"].Bounds().Size().Y)/2)*duckScale

		xDistanceToTarget := float64(duck.targetX - duck.X)
		yDistanceToTarget := float64(duck.targetY - duck.Y)

		duck.X += int(xDistanceToTarget) / 8
		duck.Y += int(yDistanceToTarget) / 8
	} else if !duck.isSleeping && !duck.takeBreak {
		duck.isFacingRight = ((duck.waitTime % 2) == 0)
		duck.isFlying = ((duck.waitTime % 4) >= 2) // ratio of last nums represents probability
		goingUp := ((duck.waitTime % 8) >= 4)
		if duck.isFacingRight {
			duck.X -= 1
		} else {
			duck.X += 1
		}
		if duck.isFlying {
			duck.isFlying = true
			duck.isWalking = false
			if goingUp {
				duck.Y += 1
			} else {
				duck.Y -= 1
			}
		} else {
			duck.isFlying = false
			duck.isWalking = true
		}
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

	if duck.isSleeping {
		screen.DrawImage(duck.Assets["duck_sleeping"], op)
	} else if duck.isHeld {
		screen.DrawImage(duck.Assets["duck"], op)
	} else if duck.isWalking {
		DrawSpriteFrame(screen, duck.Assets["duck_walk"], 30, 30, (duck.Game.frame/7)%4, op)
	} else {
		screen.DrawImage(duck.Assets["duck_sitting"], op)
	}

	ebitenutil.DebugPrintAt(screen, duck.Name, duck.X, duck.Y-20)
}
