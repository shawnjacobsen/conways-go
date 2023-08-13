package main

type DoubleBuffer struct {
	currentGrid *Grid
	nextGrid    *Grid
}

func newDoubleBuffer(currentGrid, nextGrid *Grid) *DoubleBuffer {
	return &DoubleBuffer{
		currentGrid: currentGrid,
		nextGrid:    nextGrid,
	}
}

func (db *DoubleBuffer) Swap() {
	db.currentGrid, db.nextGrid = db.nextGrid, db.currentGrid
}

func (db *DoubleBuffer) ApplyTransformation(operation func(src, dest *Grid)) {
	// apply operation
	operation(db.currentGrid, db.nextGrid)

	// swap grids
	db.Swap()
}

func (db *DoubleBuffer) getCurrentGrid() *Grid {
	return db.currentGrid
}