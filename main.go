package main

import (
	"bufio"
	"image"
	"io"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/pkg/errors"

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
		stdIn := bufio.NewReader(os.Stdin)

		for {
			line, err := stdIn.ReadBytes('\n')

			if err == io.EOF {
				break
			} else if err != nil {
				sendError(errors.WithMessage(err, "could not read line"))
			}

			err = handleRequest(line)

			if err != nil {
				sendError(err)
			}
		}
	}()

	outImg := image.NewRGBA(image.Rect(0, 0, 1920, 1080))

	logFile, _ := os.OpenFile("perf.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer logFile.Close()
	logger := log.New(logFile, "", 0)

	for !window.ShouldClose() {
		start := time.Now()
		needRedraw := false

	drawLoop:
		for {
			select {
			case req := <-drawRequests:
				req.draw(outImg)
				needRedraw = true
			case <-resized:
				needRedraw = true
			default:
				break drawLoop
			}
		}

		logger.Println("draw requests:\t", time.Now().Sub(start))
		glTime := time.Now()

		windowWidth, windowHeight := window.GetSize()

		bounds := outImg.Bounds()
		imgWidth, imgHeight := bounds.Dx(), bounds.Dy()

		gl.RasterPos2f(-1, 1)
		gl.PixelZoom(1, -1)
		gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))

		if needRedraw {
			gl.DrawPixels(
				int32(imgWidth), int32(imgHeight),
				gl.RGBA, gl.UNSIGNED_BYTE,
				unsafe.Pointer(&outImg.Pix[0]))
		}

		diff := time.Now().Sub(start)
		logger.Println("gl drawing:\t\t", time.Now().Sub(glTime))
		logger.Println("total:\t\t\t", diff)
		logger.Println()

		window.SwapBuffers()
		glfw.PollEvents()

		if diff < 30*time.Millisecond {
			time.Sleep(30*time.Millisecond - diff)
		}
	}
}
