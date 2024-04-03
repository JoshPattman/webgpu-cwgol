package main

import (
	_ "embed"
	"fmt"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

//go:embed cwgol.wgsl
var cwgolWGSL string

func main() {
	pixelgl.Run(run)
}

func run() {
	gridWidth, gridHeight := 1900, 1000
	// GPU STUFF
	m, err := NewManager()
	if err != nil {
		panic(err)
	}
	defer m.Release()

	fmt.Println("Created manager")

	f, err := m.NewShader("cwgol-shader", "cwgol", cwgolWGSL)
	if err != nil {
		panic(err)
	}
	defer f.Release()

	stateUpdater, err := f.NewUnaryOp(uint64(gridWidth * gridHeight))
	if err != nil {
		panic(err)
	}
	defer stateUpdater.Release()

	grid := NewGrid(gridWidth, gridHeight)
	for i := 0; i < int(grid.Width); i++ {
		for j := 0; j < int(grid.Height); j++ {
			if rand.Float32() < 0.5 {
				grid.SetValueAt(1, i, j)
			}
		}
	}

	// WINDOW STUFF
	cfg := pixelgl.WindowConfig{
		Title:  "Conway's Game of Life",
		Bounds: pixel.R(0, 0, float64(gridWidth), float64(gridHeight)),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Update()

		if win.Pressed(pixelgl.MouseButton1) {
			mousePos := win.MousePosition()
			mpx, mpy := int(mousePos.X), gridHeight-int(mousePos.Y)
			for i := mpx - 5; i < mpx+5; i++ {
				for j := mpy - 5; j < mpy+5; j++ {
					grid.SetValueAt(1, i, j)
				}
			}
		}

		workgroupSize := 16
		newData, err := stateUpdater.Do(grid.Data, [3]uint32{uint32(grid.Width/workgroupSize) + 1, uint32(grid.Height/workgroupSize) + 1, 1}, [3]uint32{uint32(grid.Width), uint32(grid.Height), 1})
		if err != nil {
			panic(err)
		}
		grid.Data = newData

		pic := pixel.PictureDataFromImage(grid)
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	}
}
