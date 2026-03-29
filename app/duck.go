package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Duck struct {
	Name     string
	PlayerId int

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
	duck.SetSkin("duck_green")
	duck.PlayerId = 0
}

func (duck *Duck) SetSkin(skinName string) {
	for i, name := range duck.Skins {
		if name == skinName {
			duck.Skin = i
			duck.Assets["duck"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck.png", skinName))
			duck.Assets["duck_walk"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_walk.png", skinName))
			duck.Assets["duck_sitting"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_sitting.png", skinName))
			duck.Assets["duck_sleeping"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_sleeping.png", skinName))
			duck.Assets["duck_flying"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_flying.png", skinName))
			return
		}
	}
}

func (duck *Duck) NextSkin() {
	duck.Skin += 1
	duck.Skin = duck.Skin % len(duck.Skins)
	newSkin := duck.Skins[duck.Skin]

	duck.Assets["duck"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck.png", newSkin))
	duck.Assets["duck_walk"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_walk.png", newSkin))
	duck.Assets["duck_sitting"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_sitting.png", newSkin))
	duck.Assets["duck_sleeping"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_sleeping.png", newSkin))
	duck.Assets["duck_flying"] = LoadImageFromPath(fmt.Sprintf("assets/%s/duck_flying.png", newSkin))
	//fmt.Println(*duck.Game.otherPlayerData[0])
	//message, err := json.Marshal(&VisiblePlayerDataAction{Action: "duck_customisation", PlayerId: duck.PlayerId, PlayerData: *duck.Game.otherPlayerData[duck.PlayerId]})
	//fmt.Println("SENDING MESSAGE")
	//if err == nil {
	//	duck.Game.SendMessage(message)
	//}
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

	if duck.isHeld || duck.isFlying {
		duck.isFacingRight = (duck.targetX - duck.X) < 0
	}

	if duck.isSleeping {
		duck.isWalking = false
		duck.isFlying = false
	} else if duck.isHeld {
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			duck.isHeld = false
			duck.lastTimestamp = g.frame
			duck.waitTime = rand.IntN(200) + 50
			//could make array of motion booleans that can be set to false via one func call to prevent repetition
			duck.isWalking = false
			duck.isFlying = false
			duck.takeBreak = true
		}
	} else if duck.isHovered && ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		duck.isHeld = true
		duck.isWalking = false
		duck.isFlying = false
	} else if g.frame-duck.lastTimestamp > duck.waitTime {
		duck.takeBreak = !duck.takeBreak
		duck.targetX = duck.X + rand.IntN(1000) - 500
		duck.targetY = duck.Y + rand.IntN(1000) - 500
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
		duck.isFlying = ((duck.waitTime % 4) >= 2) // ratio of last nums represents probability^-1
		//goingUp := ((duck.waitTime % 8) >= 4)
		if duck.isFlying {
			duck.isWalking = false
			xDistanceToTarget := float64(duck.targetX - duck.X)
			yDistanceToTarget := float64(duck.targetY - duck.Y)

			duck.X += int(xDistanceToTarget) / 80
			duck.Y += int(yDistanceToTarget) / 80
		} else {
			duck.isWalking = true
			duck.isFacingRight = ((duck.waitTime % 2) == 0)
			if duck.isFacingRight {
				duck.X -= 1
			} else {
				duck.X += 1
			}
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
	} else if duck.isFlying {
		DrawSpriteFrame(screen, duck.Assets["duck_flying"], 30, 30, (duck.Game.frame/7)%4, op)
	} else {
		screen.DrawImage(duck.Assets["duck_sitting"], op)
	}

	ebitenutil.DebugPrintAt(screen, duck.Name, duck.X, duck.Y-20)
}
