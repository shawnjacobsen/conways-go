package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
	screenScale  = 1
)

type Game struct {
	doubleBuffer *DoubleBuffer
	pixels []byte
}

func newPixels(width int, height int) []byte {
	return make([]byte, width*height*4)
}

func countNeighbors(grid *Grid, x, y int) int {
	// get number of neighbors
	neighborsCount := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (i == 0 && j == 0) {
				continue
			}
			x2 := (x+i+grid.width)  % grid.width
			y2 := (y+j+grid.height) % grid.height
			neighborIsAlive := *grid.getGridCell(x2, y2)
			if neighborIsAlive {
				neighborsCount++
			}
		}
	}
	return neighborsCount
}

func nextCellState(numNeighbors int, currCellState bool) bool {
	// decide if alive or not
	switch {
	case (numNeighbors < 2):
		return false
	case (numNeighbors == 2 || numNeighbors == 3) && currCellState:
		return true
	case (numNeighbors > 3):
		return false
	case (numNeighbors == 3):
		return true
	default:
		return false
	}
}

func UpdateCells(srcGrid, destGrid *Grid) {
	// init wait group to concurrently update cells
	var wg sync.WaitGroup
	
	var numRoutines int = runtime.NumCPU() * 2
	var rowsPerRoutine int = srcGrid.height / numRoutines
	wg.Add(numRoutines)
	
	// update cells
	for i := 0; i < numRoutines; i++ {
		// calculate current goroutine's responsibility
		startRow := i * rowsPerRoutine
		endRow := startRow + rowsPerRoutine
		// account for extra rows if this is the last goroutine
		if i == numRoutines - 1 {
			endRow = srcGrid.height
		}

		go func (startRow, endRow int) {
			defer wg.Done()
			for x := 0; x < srcGrid.width; x++ {
				for y := startRow; y < endRow; y++ {
					cell := srcGrid.getGridCell(x,y)
					numNeighbors := countNeighbors(srcGrid, x, y)
					*destGrid.getGridCell(x, y) = nextCellState(numNeighbors, *cell)
				}
			}
		}(startRow, endRow)
	}
	wg.Wait()
}

func (g *Game) Update() error {
	g.doubleBuffer.ApplyTransformation(UpdateCells)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// update pixels based on current game state
	for i, cell := range g.doubleBuffer.getCurrentGrid().cells {
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
	
	//  write fps to screen
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%v",ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	grid1 := newGrid(screenWidth, screenHeight)	
	grid2 := newGrid(screenWidth, screenHeight)
	randomizeGrid(grid1)
	db := newDoubleBuffer(grid1, grid2)
	game := &Game {
		doubleBuffer: db,
		pixels: newPixels(screenWidth, screenHeight),
	}
	ebiten.SetWindowSize(screenWidth*screenScale, screenHeight*screenScale)
	ebiten.SetWindowTitle("Conway's Game of Life (@Shawn Jacobsen)")
	ebiten.SetTPS(30)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}