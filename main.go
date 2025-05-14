package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
)

type bitset uint64

type Cell struct {
	width, height int32
	bitset        bitset
}

const (
	gridSize     = 128
	baseWidth    = 1000
	baseHeight   = 700
	minCellSize  = 5
	maxCellSize  = 8
	sidebarRatio = 0.3
)

var (
	selectedX, selectedY = -1, -1
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowAlwaysRun)
	rl.InitWindow(baseWidth, baseHeight, "Agents")
	defer rl.CloseWindow()
	rl.SetTargetFPS(30)

	grid := initializeGrid()

	for !rl.WindowShouldClose() {
		draw(grid)
	}
}

func initializeGrid() [gridSize][gridSize]Cell {
	var grid [gridSize][gridSize]Cell
	grid[10][10].bitset = 1 << 1
	grid[20][15].bitset = 1 << 2
	grid[50][30].bitset = 1 << 3
	grid[10][5].bitset = 1<<1 | 1<<2 | 1<<3
	return grid
}

func draw(grid [gridSize][gridSize]Cell) {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.Black)

	sw := int32(rl.GetScreenWidth())
	sh := int32(rl.GetScreenHeight())

	sidebarW := int32(float32(sw) * sidebarRatio)

	dynCell := int32((sw - sidebarW) / gridSize)
	dynCell = max(minCellSize, min(dynCell, maxCellSize))

	gridPixels := dynCell * gridSize
	offsetX := (sw - sidebarW - gridPixels) / 2
	offsetY := max((sh-gridPixels)/2, 0)

	handleMouse(offsetX, offsetY, dynCell)
	drawGrid(grid, offsetX, offsetY, dynCell)
	drawSidebar(grid, sidebarW)

	rl.DrawFPS(10, sh-30)
}

func handleMouse(offsetX, offsetY, cellSize int32) {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		x, y := rl.GetMouseX(), rl.GetMouseY()
		gx := (x - offsetX) / cellSize
		gy := (y - offsetY) / cellSize
		if gx >= 0 && gx < gridSize && gy >= 0 && gy < gridSize {
			selectedX, selectedY = int(gx), int(gy)
		}
	}
}

func drawGrid(grid [gridSize][gridSize]Cell, ox, oy, size int32) {
	for y := range gridSize {
		for x := range gridSize {
			cell := grid[y][x]
			col := rl.DarkGray
			if cell.bitset != 0 {
				col = rl.Green
			}
			if x == selectedX && y == selectedY {
				col = rl.Yellow
			}
			rl.DrawRectangle(ox+int32(x)*size, oy+int32(y)*size, size-1, size-1, col)
		}
	}
}

func drawSidebar(grid [gridSize][gridSize]Cell, sidebarW int32) {
	sw := int32(rl.GetScreenWidth())
	tsh := int32(rl.GetScreenHeight())
	rl.DrawRectangle(sw-sidebarW, 0, sidebarW, tsh, rl.RayWhite)
	rl.DrawText("Cell Info", sw-sidebarW+10, 20, 22, rl.Black)

	extY := int32(60)
	if selectedX >= 0 && selectedY >= 0 {
		cell := grid[selectedY][selectedX]
		info := fmt.Sprintf("X: %d\nY: %d\nBitset: %x", selectedX, selectedY, cell.bitset)
		rl.DrawText(info, sw-sidebarW+10, extY, 20, rl.Black)
	} else {
		rl.DrawText("Click a cell to view info", sw-sidebarW+10, extY, 20, rl.Gray)
	}
}
