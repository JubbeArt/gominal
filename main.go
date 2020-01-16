package main

import (
	"bufio"
	"image"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

var (
	window    *glfw.Window
	colWidth  = 12
	rowHeight = 24
	rows      int
	cols      int

	drawRequests = make(chan drawRequest, 100)
	resized      = make(chan struct{}, 10)
)

func main() {
	err := glfw.Init()

	if err != nil {
		sendError(err)
		return
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Visible, glfw.False)

	windowState := savedWindowState()

	window, err = glfw.CreateWindow(windowState.Width, windowState.Height, "", nil, nil)
	if err != nil {
		sendError(err)
		return
	}

	if windowState.X != -1 {
		window.SetPos(windowState.X, windowState.Y)
	}

	window.Show()

	{
		width, height := window.GetSize()
		cols = width / colWidth
		rows = height / rowHeight
	}

	sendResponse(sizeResponse{Type: "size", Rows: rows, Cols: cols, ColWidth: colWidth, RowHeight: rowHeight})

	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		sendError(err)
		return
	}

	loadFonts(18)

	window.SetSizeCallback(func(_ *glfw.Window, width int, height int) {
		sizeCallback(width, height)
		resized <- struct{}{}
	})

	window.SetCharCallback(func(_ *glfw.Window, char rune) {
		charCallback(char)
	})

	window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
		keyCallback(key, scanCode, action, mods)
	})

	window.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		mouseCallback(button, action, mods)
	})

	window.SetCloseCallback(func(_ *glfw.Window) {
		state := WindowState{}
		state.Width, state.Height = window.GetSize()
		state.X, state.Y = window.GetPos()
		state.save()
	})

	go func() {
		stdIn := bufio.NewScanner(os.Stdin)

		for stdIn.Scan() {
			line := stdIn.Bytes()
			err := handleRequest(line)

			if err != nil {
				sendError(err)
			}
		}

		window.SetShouldClose(true)
	}()

	outImg := image.NewRGBA(image.Rect(0, 0, 1920, 1080))

	logFile, _ := os.OpenFile("perf.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	logger := log.New(logFile, "", 0)
	defer logFile.Close()

	for !window.ShouldClose() {
		start := time.Now()
		// prolly not needed here
		window.MakeContextCurrent()
		drawTime := time.Now()
		needRedraw := false

	drawLoop:
		for {
			select {
			case req := <-drawRequests:
				req.draw(outImg)
				needRedraw = true
			case <-resized:
				_ = true
			//	needRedraw = true
			default:
				break drawLoop
			}
		}

		logger.Println("draw requests ", time.Now().Sub(drawTime))

		windowWidth, windowHeight := window.GetSize()

		bounds := outImg.Bounds()
		imgWidth, imgHeight := bounds.Dx(), bounds.Dy()

		glTime := time.Now()

		if needRedraw {

		}

		// has to be called every time?
		gl.RasterPos2f(-1, 1)
		gl.PixelZoom(1, -1)
		gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
		gl.DrawPixels(
			int32(imgWidth), int32(imgHeight),
			gl.RGBA, gl.UNSIGNED_BYTE,
			unsafe.Pointer(&outImg.Pix[0]))

		logger.Println("gl drawing ", time.Now().Sub(glTime))

		window.SwapBuffers()
		glfw.PollEvents()

		diff := time.Now().Sub(start)

		logger.Println("total ", diff)
		logger.Println()

		if diff < 30*time.Millisecond {
			time.Sleep(30*time.Millisecond - diff)
		}
	}
}
