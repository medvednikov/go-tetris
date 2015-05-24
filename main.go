package main

import (
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var window *glfw.Window

type Point struct{ X, Y int }

func init() {
	runtime.LockOSThread()
}

func main() {
	initGame()
	// Create the window
	var err error
	err = glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	window, err = glfw.CreateWindow(W, H, "go-tetris", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetKeyCallback(keyPress)
	if err := gl.Init(); err != nil {
		panic(err)
	}
	// Timer
	ticker := time.NewTicker(time.Millisecond * TimerPeriod)
	go func() {
		for range ticker.C {
			//fmt.Println("tick ", t)
			update()
		}
	}()
	// Init OpenGL
	gl.Ortho(0, W, H, 0, -1, 1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	//gl.ClearColor(20, 20, 0, 0)
	gl.ClearColor(255, 255, 255, 0)
	gl.LineWidth(1)
	gl.Color3f(1, 0, 0)
	for !window.ShouldClose() {
		drawScene()
		glfw.PollEvents()
	}
}

func glVertex(x, y int) {
	gl.Vertex2i(int32(x), int32(y))
}

func Line(x1, y1, x2, y2 int) {
	gl.Begin(gl.LINES)
	glVertex(x1, y1)
	glVertex(x2, y2)
	gl.End()
}

func Square(p1, p2, p3, p4 Point) {
	gl.Begin(gl.POLYGON)
	glVertex(p1.X, p1.Y)
	glVertex(p2.X, p2.Y)
	glVertex(p3.X, p3.Y)
	glVertex(p4.X, p4.Y)
	gl.End()
}

func pow(a, n int) int {
	res := 1
	for i := 0; i < n; i++ {
		res *= a
	}
	return res
}
