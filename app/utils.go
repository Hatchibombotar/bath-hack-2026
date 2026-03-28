package main

import (
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/**
var emb embed.FS

func LoadImageFromPath(path string) *ebiten.Image {
	file, err := emb.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	sheet := ebiten.NewImageFromImage(img)

	return sheet
}

func isPointInRect(px, py, rx, ry, rw, rh int) bool {
	// normalize width/height to handle negative values
	x0, x1 := rx, rx+rw
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	y0, y1 := ry, ry+rh
	if y0 > y1 {
		y0, y1 = y1, y0
	}

	// check inclusive on left/top, exclusive on right/bottom
	return px >= x0 && px < x1 && py >= y0 && py < y1
}

func DrawSpriteFrame(screen *ebiten.Image, spritesheet *ebiten.Image, frameWidth, frameHeight, frameIndex int, op *ebiten.DrawImageOptions) {
	sx := frameIndex * frameWidth
	sr := image.Rect(sx, 0, sx+frameWidth, frameHeight)

	sub := spritesheet.SubImage(sr).(*ebiten.Image)

	screen.DrawImage(sub, op)
}

func ReadFileBytes(path string) []byte {
	bytes, err := emb.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return bytes
}

// AI
func DrawNineSlice(screen *ebiten.Image, img *ebiten.Image, size int, x, y, scale float64) {
	sb := img.Bounds()
	sw := sb.Dx()
	sh := sb.Dy()

	// Source patch rectangles (in source image coordinates)
	// Left, center, right widths
	lw := size
	rw := size
	cw := sw - lw - rw
	if cw < 0 {
		cw = 0
	}
	// Top, middle, bottom heights
	th := size
	bh := size
	mh := sh - th - bh
	if mh < 0 {
		mh = 0
	}

	// Target total size in pixels (before scale). If user wants to define a target width/height
	// different from source, we can compute target widths/heights by scale; here scale scales whole result.
	// Compute target patch sizes (scaled later by `scale`)
	// Corners keep original size; edges/center will stretch to fill available space.
	tlw := float64(lw)
	trw := float64(rw)
	tth := float64(th)
	tbh := float64(bh)
	// For demonstration we scale to maintain same overall source size then apply scale.
	// If you want a different final width/height, compute targetW/targetH and adjust cw/mh accordingly.
	targetCW := float64(cw)
	targetMH := float64(mh)

	// Convenience: function to draw a source rect to a destination rect (in screen coords)
	drawPatch := func(sx0, sy0, sx1, sy1 int, dx0, dy0, dx1, dy1 float64) {
		if sx1 <= sx0 || sy1 <= sy0 || dx1 <= dx0 || dy1 <= dy0 {
			return
		}
		src := img.SubImage(image.Rect(sx0, sy0, sx1, sy1)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		// Destination rectangle in device pixels: we map source (0,0)-(w,h) to destination rect using GeoM
		sw := float64(sx1 - sx0)
		sh := float64(sy1 - sy0)
		scaleX := (dx1 - dx0) / sw
		scaleY := (dy1 - dy0) / sh
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(dx0, dy0)
		screen.DrawImage(src, op)
	}

	// Compute unscaled destination positions
	// left x, center x, right x
	dx0 := x
	dx1 := x + tlw*scale
	dx2 := dx1 + targetCW*scale
	dx3 := dx2 + trw*scale
	// top y, middle y, bottom y
	dy0 := y
	dy1 := y + tth*scale
	dy2 := dy1 + targetMH*scale
	dy3 := dy2 + tbh*scale

	// Source coordinates:
	sx0 := 0
	sx1 := lw
	sx2 := lw + cw
	sx3 := sw

	sy0 := 0
	sy1 := th
	sy2 := th + mh
	sy3 := sh

	// Top-left
	drawPatch(sx0, sy0, sx1, sy1, dx0, dy0, dx1, dy1)
	// Top-center
	drawPatch(sx1, sy0, sx2, sy1, dx1, dy0, dx2, dy1)
	// Top-right
	drawPatch(sx2, sy0, sx3, sy1, dx2, dy0, dx3, dy1)

	// Middle-left
	drawPatch(sx0, sy1, sx1, sy2, dx0, dy1, dx1, dy2)
	// Middle-center
	drawPatch(sx1, sy1, sx2, sy2, dx1, dy1, dx2, dy2)
	// Middle-right
	drawPatch(sx2, sy1, sx3, sy2, dx2, dy1, dx3, dy2)

	// Bottom-left
	drawPatch(sx0, sy2, sx1, sy3, dx0, dy2, dx1, dy3)
	// Bottom-center
	drawPatch(sx1, sy2, sx2, sy3, dx1, dy2, dx2, dy3)
	// Bottom-right
	drawPatch(sx2, sy2, sx3, sy3, dx2, dy2, dx3, dy3)
}
