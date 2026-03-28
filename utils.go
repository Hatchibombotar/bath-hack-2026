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
