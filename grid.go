package main

import (
	"image"
	"image/color"
)

type Grid struct {
	Data   []float32
	Width  int
	Height int
}

func (g *Grid) ValueAt(x, y int) float32 {
	return g.Data[y*g.Width+x]
}
func (g *Grid) SetValueAt(v float32, x, y int) {
	g.Data[y*g.Width+x] = v
}

// At implements image.Image.
func (g *Grid) At(x int, y int) color.Color {
	return color.Gray{
		Y: uint8(g.ValueAt(x, y) * 255),
	}
}

// Bounds implements image.Image.
func (g *Grid) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(g.Width), int(g.Height))
}

// ColorModel implements image.Image.
func (g *Grid) ColorModel() color.Model {
	return color.GrayModel
}

func NewGrid(width, height int) *Grid {
	return &Grid{
		Data:   make([]float32, width*height),
		Width:  width,
		Height: height,
	}
}

var _ image.Image = &Grid{}
