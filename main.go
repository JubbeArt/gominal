package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/pkg/errors"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

const (
	colWidth  = 12
	rowHeight = 24
)

var (
	cols int
	rows int
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

	windowState := getWindowState()

	win, err := glfw.CreateWindow(windowState.Width, windowState.Height, "", nil, nil)
	if err != nil {
		sendError(err)
		return
	}

	if windowState.X != -1 {
		win.SetPos(windowState.X, windowState.Y)
	}

	win.Show()
	win.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		sendError(err)
		return
	}

	loadFonts(18)
	setupCallbacks(win)

	fmt.Println("RUNNING")

	quit := make(chan struct{})

	go func() {
		stdIn := bufio.NewReader(os.Stdin)

		for {
			line, err := stdIn.ReadBytes('\n')

			if err == io.EOF {
				quit <- struct{}{}
			} else if err != nil {
				sendError(errors.WithMessage(err, "could not read line"))
			}

			err = handleRequest(win, line)

			if err != nil {
				sendError(err)
			}
		}
	}()

	//outImg := image.NewRGBA(image.Rect(0, 0, 1920, 1080))

	logFile, _ := os.OpenFile("perf.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer logFile.Close()
	logger := log.New(logFile, "", 0)

	for !win.ShouldClose() {
		start := time.Now()
		//needRedraw := false

		select {
		case <-quit:
			win.SetShouldClose(true)
		default:

		}

		//drawLoop:
		//	for {
		//		select {
		//		case req := <-drawRequests:
		//			req.draw(outImg)
		//			needRedraw = true
		//		case <-resized:
		//			needRedraw = true
		//		default:
		//			break drawLoop
		//		}
		//	}

		logger.Println("draw requests:\t", time.Now().Sub(start))
		glTime := time.Now()

		windowWidth, windowHeight := win.GetSize()

		//bounds := outImg.Bounds()
		//imgWidth, imgHeight := bounds.Dx(), bounds.Dy()

		gl.RasterPos2f(-1, 1)
		gl.PixelZoom(1, -1)
		gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))

		//if needRedraw {
		//	gl.DrawPixels(
		//		int32(imgWidth), int32(imgHeight),
		//		gl.RGBA, gl.UNSIGNED_BYTE,
		//		unsafe.Pointer(&outImg.Pix[0]))
		//}

		diff := time.Now().Sub(start)
		logger.Println("gl drawing:\t\t", time.Now().Sub(glTime))
		logger.Println("total:\t\t\t", diff)
		logger.Println()

		win.SwapBuffers()
		glfw.PollEvents()

		if diff < 30*time.Millisecond {
			time.Sleep(30*time.Millisecond - diff)
		}
	}
}

func setupCallbacks(win *glfw.Window) {
	win.SetSizeCallback(sizeCallback)
	win.SetCharCallback(charCallback)
	win.SetKeyCallback(keyCallback)
	win.SetMouseButtonCallback(mouseClickCallback)
	win.SetCursorPosCallback(mouseMoveCallback)

	win.SetCloseCallback(func(_ *glfw.Window) {
		state := WindowState{}
		state.Width, state.Height = win.GetSize()
		state.X, state.Y = win.GetPos()
		state.save()
	})
}
