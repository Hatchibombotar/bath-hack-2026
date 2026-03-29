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

	X, Y             float64
	targetX, targetY float64
	Game             *Game
	nestX, nestY     int

	isAtNest bool

	isHovered  bool //if mouse intersecting bounding box of image
	isSleeping bool //true if in a sleeping state (in a study session)
	isHeld     bool //true if currently being dragged by mouse

	isMoving  bool
	isWalking bool
	isFlying  bool

	isFacingRight bool //true if facing right, false if left

	isOtherDuck bool // true if this is not the client's duck

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

func (duck *Duck) GetSkin() string {
	return duck.Skins[duck.Skin]
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

	duck.Game.SendDuckInfo()
}

func (duck *Duck) Update() {
	if !duck.isOtherDuck {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			duck.NextSkin()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			duck.isSleeping = !duck.isSleeping
		}

		g := duck.Game

		duck.isAtNest = (float64(duck.nestX+duck.nestY) - duck.X - duck.Y) < 6

		duck.isHovered = isPointInRect(
			g.cursorX, g.cursorY,
			int(g.duck.X), int(g.duck.Y),
			duckWidth*duckScale,
			duckWidth*duckScale,
		)

		if duck.isHeld || duck.isMoving {
			duck.isFacingRight = (duck.targetX - duck.X) < 0
		}

		if duck.isSleeping {
			duck.isFacingRight = false
			duck.isMoving = false
			duck.isWalking = false
			duck.isFlying = false
		} else if duck.isHeld {
			if !ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
				duck.isHeld = false
				duck.lastTimestamp = g.frame
				duck.waitTime = rand.IntN(200) + 50
				//could make array of motion booleans that can be set to false via one func call to prevent repetition
				duck.isWalking = false
				duck.isFlying = false
			}
		} else if duck.isHovered && ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
			duck.isHeld = true
			duck.isWalking = false
			duck.isFlying = false
		} else if g.frame-duck.lastTimestamp > duck.waitTime {
			duck.isMoving = !duck.isMoving
			maxX, maxY := g.ScreenSize()
			if duck.isMoving {
				duck.targetX = float64(rand.IntN(maxX))
				if rand.IntN(2) != 1 {
					duck.isFlying = true
					duck.isWalking = false
					duck.targetY = float64(rand.IntN(maxY - 35))
				} else {
					duck.isWalking = true
					duck.isFlying = false
					duck.targetY = duck.Y
				}
			} else {
				duck.isFlying = false
				duck.isWalking = false
			}

			duck.lastTimestamp = g.frame
			duck.waitTime = rand.IntN(400) + 100
		}
	} else {
		if duck.isHeld || duck.isMoving {
			duck.isFacingRight = (duck.targetX - duck.X) < 0
		}
		duck.isAtNest = (float64(duck.nestX+duck.nestY) - duck.X - duck.Y) < 6

		if duck.isSleeping {
			duck.targetY = float64(duck.nestY) + 50
		} else {
			duck.targetY = float64(duck.nestY)
		}
		//do client duck behaviour
	}
}

func (duck *Duck) Move() {
	if duck.isSleeping {
		xDistanceToTarget := float64(duck.nestX) - duck.X
		yDistanceToTarget := float64(duck.nestY) - duck.Y

		duck.X += float64(xDistanceToTarget) / 30
		duck.Y += float64(yDistanceToTarget) / 30

		if duck.isAtNest {
			duck.isFacingRight = true
		} else {
			duck.isFacingRight = (float64(duck.nestX) - duck.X) < 0
		}

		//fmt.Println(xDistanceToTarget + yDistanceToTarget)
	} else if duck.isHeld {
		duck.targetX = float64(duck.Game.cursorX - (int(duck.Assets["duck"].Bounds().Size().X)/2)*duckScale)
		duck.targetY = float64(duck.Game.cursorY - (int(duck.Assets["duck"].Bounds().Size().Y)/2)*duckScale)

		xDistanceToTarget := float64(duck.targetX - duck.X)
		yDistanceToTarget := float64(duck.targetY - duck.Y)

		duck.X += float64(xDistanceToTarget) / 8
		duck.Y += float64(yDistanceToTarget) / 8
	} else if duck.isMoving {
		xDistanceToTarget := float64(duck.targetX - duck.X)
		yDistanceToTarget := float64(duck.targetY - duck.Y)
		if duck.isWalking {
			if duck.targetX > duck.X {
				duck.X += 1.0
			} else {
				duck.X -= 1.0
			}
		} else if duck.isFlying {
			duck.X += float64(xDistanceToTarget) / 80
			duck.Y += float64(yDistanceToTarget) / 80
		}
		if xDistanceToTarget < 3 && yDistanceToTarget < 3 {
			duck.isMoving = false
			duck.isFlying = false
			duck.isWalking = false
			duck.lastTimestamp = duck.Game.frame
			duck.waitTime = rand.IntN(400) + 100
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
		if duck.isAtNest {
			DrawSpriteFrame(screen, duck.Assets["duck_sleeping"], 30, 30, (duck.Game.frame/20)%3, op)
		} else {
			DrawSpriteFrame(screen, duck.Assets["duck_flying"], 30, 30, (duck.Game.frame/7)%4, op)
		}
	} else if duck.isHeld {
		screen.DrawImage(duck.Assets["duck"], op)
	} else if duck.isWalking {
		DrawSpriteFrame(screen, duck.Assets["duck_walk"], 30, 30, (duck.Game.frame/7)%4, op)
	} else if duck.isFlying {
		DrawSpriteFrame(screen, duck.Assets["duck_flying"], 30, 30, (duck.Game.frame/7)%4, op)
	} else {
		screen.DrawImage(duck.Assets["duck_sitting"], op)
	}

	ebitenutil.DebugPrintAt(screen, duck.Name, int(duck.X), int(duck.Y-20))
}
