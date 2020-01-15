package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/image/font"
)

func init() {
	runtime.LockOSThread()
}

var (
	window     *glfw.Window
	colWidth   = 12
	rowHeight  = 24
	fontNormal font.Face
	fontBold   font.Face
	rows       int
	cols       int

	drawRequests = make(chan drawRequest, 100)
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
	//glfw.WindowHint(glfw.Visible, glfw.False)
	//window.Show()

	window, err = glfw.CreateWindow(400, 200, "", nil, nil)
	if err != nil {
		sendError(err)
		return
	}

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

	fontNormal, fontBold = loadFonts(18)

	window.SetSizeCallback(func(_ *glfw.Window, width int, height int) {
		newCols := width / colWidth
		newRows := height / rowHeight

		if newCols == cols && newRows == rows {
			return
		}

		cols = newCols
		rows = newRows

		sendResponse(sizeResponse{Type: "size", Rows: rows, Cols: cols, ColWidth: colWidth, RowHeight: rowHeight})
	})

	window.SetCharCallback(func(_ *glfw.Window, char rune) {
		sendResponse(charResponse{
			Type: "char",
			Rune: string(char),
		})
	})

	window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			keyName := glfw.GetKeyName(key, scanCode)

			if keyName == "" {
				keyName = keyLookup[key]
			}

			if keyName == "" {
				keyName = "unknown"
			}

			sendResponse(keyResponse{
				Type:  "key",
				Key:   keyName,
				Ctrl:  mods&glfw.ModControl != 0,
				Shift: mods&glfw.ModShift != 0,
				Alt:   mods&glfw.ModAlt != 0,
				Super: mods&glfw.ModSuper != 0,
			})
		}
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

		window.MakeContextCurrent()

		drawTime := time.Now()

	drawLoop:
		for {
			select {
			case req := <-drawRequests:
				req.draw(outImg)
			default:
				break drawLoop
			}
		}

		logger.Println("draw requests ", time.Now().Sub(drawTime))

		windowWidth, windowHeight := window.GetSize()

		bounds := outImg.Bounds()
		imgWidth, imgHeight := bounds.Dx(), bounds.Dy()

		glTime := time.Now()

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

func sendError(err error) {
	sendErrorStr(err.Error())
}

func sendErrorStr(err string) {
	sendResponse(errorResponse{
		Type:  "error",
		Error: err,
	})
}

func sendResponse(response interface{}) {
	bytes, err := json.Marshal(response)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

var keyLookup = map[glfw.Key]string{
	glfw.KeyLeftControl:  "ctrl",
	glfw.KeyRightControl: "ctrl",
	glfw.KeyLeftShift:    "shift",
	glfw.KeyRightShift:   "shift",
	glfw.KeyLeftAlt:      "alt",
	glfw.KeyRightAlt:     "alt",
	glfw.KeyLeftSuper:    "super",
	glfw.KeyRightSuper:   "super",
	glfw.KeyTab:          "tab",
	glfw.KeyEnter:        "enter",
	glfw.KeySpace:        "space",
	glfw.KeyBackspace:    "backspace",
	glfw.KeyEscape:       "escape",
	glfw.KeyLeft:         "left",
	glfw.KeyRight:        "right",
	glfw.KeyUp:           "up",
	glfw.KeyDown:         "down",
	glfw.KeyCapsLock:     "capsLock",
	glfw.KeyDelete:       "delete",
	glfw.KeyInsert:       "insert",
	glfw.KeyHome:         "home",
	glfw.KeyPageUp:       "pageUp",
	glfw.KeyPageDown:     "pageDown",
	glfw.KeyEnd:          "end",
	glfw.KeyNumLock:      "numLock",

	glfw.KeyF1:  "f1",
	glfw.KeyF2:  "f2",
	glfw.KeyF3:  "f3",
	glfw.KeyF4:  "f4",
	glfw.KeyF5:  "f5",
	glfw.KeyF6:  "f6",
	glfw.KeyF7:  "f7",
	glfw.KeyF8:  "f8",
	glfw.KeyF9:  "f9",
	glfw.KeyF10: "f10",
	glfw.KeyF11: "f11",
	glfw.KeyF12: "f12",
	glfw.KeyF13: "f13",
	glfw.KeyF14: "f14",
	glfw.KeyF15: "f15",
	glfw.KeyF16: "f16",
	glfw.KeyF17: "f17",
	glfw.KeyF18: "f18",
	glfw.KeyF19: "f19",
	glfw.KeyF20: "f20",
	glfw.KeyF21: "f21",
	glfw.KeyF22: "f22",
	glfw.KeyF23: "f23",
	glfw.KeyF24: "f24",
	glfw.KeyF25: "f25",
}
