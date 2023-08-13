package main

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)


type Grid struct {
	cells []bool
	height int
	width int
}

func newGrid(width int, height int) *Grid {
	grid := &Grid {
		cells: make([]bool, width * height),
		height: height,
		width: width, 
	}
	for i := range grid.cells {
		grid.cells[i] = rand.Intn(50) == 1
	}

	return grid
}

func (grid *Grid) getGridCell(x int, y int) *bool {
	return &grid.cells[y * grid.width + x]
}

const (
	screenWidth  = 1920 / 4
	screenHeight = 1080 / 4
)

type Game struct {
	grid *Grid
	pixels []byte
}

func newPixels(width int, height int) []byte {
	return make([]byte, width*height*4)
}

func (g *Game) Update() error {
	// update cells
	newGrid := newGrid(g.grid.width, g.grid.height)
	for x := 0; x < g.grid.width; x++ {
		for y := 0; y < g.grid.height; y++ {
			cell := g.grid.getGridCell(x,y)
			// get number of neighbors
			neighborsCount := 0
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					if (i == 0 && j == 0) {
						continue
					}
					x2 := (x+i+g.grid.width)  % g.grid.width
					y2 := (y+j+g.grid.height) % g.grid.height
					neighborIsAlive := *g.grid.getGridCell(x2, y2)
					if neighborIsAlive {
						neighborsCount++
					}
				}
			}

			// decide if alive or not
			switch {
			case (neighborsCount < 2):
				*newGrid.getGridCell(x,y) = false
			case (neighborsCount == 2 || neighborsCount == 3) && *cell:
				*newGrid.getGridCell(x,y) = true
			case (neighborsCount > 3):
				*newGrid.getGridCell(x,y) = false
			case (neighborsCount == 3):
				*newGrid.getGridCell(x,y) = true
			}
		}
	}
	g.grid = newGrid
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// update pixels based on current game state
	for i, cell := range g.grid.cells {
		if cell {
			g.pixels[4*i+0] = 0xFF
			g.pixels[4*i+1] = 0xFF
			g.pixels[4*i+2] = 0xFF
			g.pixels[4*i+3] = 0xFF
		} else {
			g.pixels[4*i+0] = 0
			g.pixels[4*i+1] = 0
			g.pixels[4*i+2] = 0
			g.pixels[4*i+3] = 0
		}
	}

	// draw pixel
	screen.WritePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	
	game := &Game {
		grid: newGrid(screenWidth, screenHeight),
		pixels: newPixels(screenWidth, screenHeight),
	}
	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Conway's Game of Life (@Shawn Jacobsen)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}