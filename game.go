package main

import (
	"math/rand"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	BlockSize   = 20
	FieldHeight = 20
	FieldWidth  = 10
	TetroSize   = 4

	// Window size
	W = BlockSize * FieldWidth
	H = BlockSize * FieldHeight

	TimerPeriod = 250 // milliseconds
)

type Color struct{ R, G, B int }

type Tetro [TetroSize]struct{ x, y int }

var (
	// Position of the dropping tetromino
	posX = 0
	posY = 0

	// field[y][x] contains the color of the block with (x,y) coordinates
	// (unless < 1). Has borders comprised of -1, so that bounds checking is
	// not needed:
	// -1 -1 -1 -1
	// -1  0  0 -1
	// -1  0  0 -1
	// -1 -1 -1 -1
	field [][]int

	BinaryTetros = [][]int{
		// 0000 0
		// 0000 0
		// 0110 6
		// 0110 6
		[]int{66, 66, 66, 66},
		// 0000 0
		// 0000 0
		// 0010 2
		// 0111 7
		[]int{27, 131, 72, 232},
		// 0000 0
		// 0000 0
		// 0011 3
		// 0110 6
		[]int{36, 231, 36, 231},
		// 0000 0
		// 0000 0
		// 0110 6
		// 0011 3
		[]int{63, 132, 63, 132},
		// 0000 0
		// 0011 3
		// 0001 1
		// 0001 1
		[]int{311, 17, 223, 74},
		// 0000 0
		// 0011 3
		// 0010 1
		// 0010 1
		[]int{322, 71, 113, 47},
		// Special case since 15 can't be used
		// 1111
		[]int{1111, 9, 1111, 9},
	}
	tetro Tetro

	// Index of the dropping tetromino. Refers to its color.
	tetroIdx = 0

	// Index of the rotation (0-3)
	rotationIdx = 0

	Colors = []Color{
		Color{0, 0, 0},
		Color{170, 0, 0},
		Color{192, 192, 192},
		Color{170, 0, 170},
		Color{0, 0, 170},
		Color{0, 170, 0},
		Color{170, 85, 0},
		Color{0, 170, 170},
	}
)

func update() {
	moveDown()
	deleteCompletedLines()
}

func moveDown() {
	// Check each block in the dropping tetro
	for i := 0; i < TetroSize; i++ {
		y := tetro[i].y + posY + 1
		x := tetro[i].x + posX
		// Reached bottom of the screen or another block?
		if field[y][x] != 0 {
			// End of game?
			if posY < 2 {
				initGame()
				return
			}
			// Leave it and generate a new one
			leaveTetro()
			generateTetro()
			return
		}
	}
	posY++
}

func moveRight(dx int) {
	for i := 0; i < TetroSize; i++ {
		// Reached left/right edges?
		y := tetro[i].y + posY
		x := tetro[i].x + posX + dx
		if field[y][x] != 0 {
			// Do not move
			return
		}
	}
	posX += dx
}

func deleteCompletedLines() {
	for y := FieldHeight; y >= 1; y-- {
		deleteCompletedLine(y)
	}
}

func deleteCompletedLine(y int) {
	for x := 1; x <= FieldWidth; x++ {
		if field[y][x] == 0 {
			return
		}
	}
	// Move everything down by 1 position
	for y = y - 1; y >= 1; y-- {
		for x := 1; x <= FieldWidth; x++ {
			field[y+1][x] = field[y][x]
		}
	}
}

// Place a new tetromino on top
func generateTetro() {
	posY = 0
	posX = FieldWidth/2 - TetroSize/2
	tetroIdx = rand.Intn(len(BinaryTetros))
	rotationIdx = 0
	tetro = parseBinaryTetro(BinaryTetros[tetroIdx][0])
}

func leaveTetro() {
	for i := 0; i < TetroSize; i++ {
		x := tetro[i].x + posX
		y := tetro[i].y + posY
		// Remember the color of each block
		field[y][x] = tetroIdx + 1
	}
}

func parseBinaryTetro(t int) (res Tetro) {
	cnt := 0
	horizontal := t == 9 // special case for the horizontal line
	for i := 0; i <= 3; i++ {
		// Get ith digit of t
		p := pow(10, 3-i)
		digit := t / p
		t %= p
		// Convert the digit to binary
		for j := 3; j >= 0; j-- {
			bin := digit % 2
			digit /= 2
			if bin == 1 || (horizontal && i == TetroSize-1) {
				res[cnt].x = j
				res[cnt].y = i
				cnt++
			}
		}
	}
	return res
}

func setColor(colorIdx int) {
	c := Colors[colorIdx]
	gl.Color3ub(uint8(c.R), uint8(c.G), uint8(c.B))
}

func drawTetromino() {
	setColor(tetroIdx + 1)
	for i := 0; i < TetroSize; i++ {
		drawBlock(posY+tetro[i].y, posX+tetro[i].x)
	}
}

func drawBlock(i, j int) {
	j-- // Handle -1 borders in the [][]field
	i--
	Square(Point{j * BlockSize, i * BlockSize},
		Point{(j+1)*BlockSize - 1, i * BlockSize},
		Point{(j+1)*BlockSize - 1, (i+1)*BlockSize - 1},
		Point{j * BlockSize, (i+1)*BlockSize - 1})
}

func drawField() {
	for i := 1; i < FieldHeight+1; i++ {
		for j := 1; j < FieldWidth+1; j++ {
			if field[i][j] > 0 {
				setColor(field[i][j])
				drawBlock(i, j)
			}
		}
	}
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	drawTetromino()
	drawField()
	window.SwapBuffers()
}

func keyPress(w *glfw.Window, k glfw.Key, s int, act glfw.Action, mods glfw.ModifierKey) {
	if act != glfw.Press {
		return
	}
	// Rotation
	switch k {
	case glfw.KeyUp:
		rotationIdx++
		if rotationIdx == TetroSize {
			rotationIdx = 0
		}
		tetro = parseBinaryTetro(BinaryTetros[tetroIdx][rotationIdx])
		if posX < 0 {
			posX = 1
		}
	// Move left/right
	case glfw.KeyLeft:
		moveRight(-1)
	case glfw.KeyRight:
		moveRight(1)
	case glfw.KeyDown:
		moveDown()
	}
}

func initGame() {
	rand.Seed(time.Now().Unix())
	field = make([][]int, FieldHeight+2)
	for i := 0; i < FieldHeight+2; i++ {
		field[i] = make([]int, FieldWidth+2)
		for j := 0; j < FieldWidth+2; j++ {
			field[i][j] = 0
			if i == 0 || i == FieldHeight+1 ||
				j == 0 || j == FieldWidth+1 {
				field[i][j] = -1
			}
		}
	}
	generateTetro()
}
