package main

import "math/rand"

type Grid struct {
	cells  []bool
	height int
	width  int
}

func newGrid(width int, height int) *Grid {
	grid := &Grid{
		cells:  make([]bool, width*height),
		height: height,
		width:  width,
	}
	return grid
}

func randomizeGrid(grid *Grid) {
	for i := range grid.cells {
		grid.cells[i] = rand.Intn(10) == 1
	}

}

func (grid *Grid) getGridCell(x int, y int) *bool {
	return &grid.cells[y*grid.width+x]
}
