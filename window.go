package gominal

import (
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

type window struct {
	win        *glfw.Window
	userRender func(ui *UI)

	width  int
	height int

	resizeCallback func(width, height int)
	charCallback   func(char rune)
	keyCallback    func(event KeyEvent)
}

type KeyEvent struct {
	Key   glfw.Key
	Ctrl  bool
	Shift bool
}

func Window(render func(ui *UI)) *window {
	win := &window{
		win:            nil,
		userRender:     render,
		width:          1200,
		height:         800,
		resizeCallback: func(width, height int) {},
		charCallback:   func(char rune) {},
		keyCallback:    func(event KeyEvent) {},
	}
	return win
}

func (w *window) Size(width, height int) *window {
	w.width = width
	w.height = height
	return w
}

func (w *window) OnResize(handler func(width, height int)) *window {
	w.resizeCallback = handler
	return w
}

func (w *window) OnKeyDown(handler func(event KeyEvent)) *window {
	w.keyCallback = handler
	return w
}

func (w *window) OnChar(handler func(char rune)) *window {
	w.charCallback = handler
	return w
}

func (w *window) Run() {
	err := glfw.Init()

	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	w.win, err = glfw.CreateWindow(w.width, w.height, "Gominal", nil, nil)
	if err != nil {
		panic(err)
	}

	w.win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	w.win.SetSizeCallback(func(win *glfw.Window, width int, height int) {
		w.resizeCallback(width, height)
	})

	w.win.SetCharCallback(func(win *glfw.Window, char rune) {
		w.charCallback(char)
	})

	w.win.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			w.keyCallback(KeyEvent{
				Key:   key,
				Ctrl:  mods&glfw.ModControl == 0,
				Shift: mods&glfw.ModShift == 0,
			})
		}
	})

	ui := newUI(w.win)
	defer ui.font.Close()

	//totalForAvr := time.Duration(0)
	//runs := 0

	for !w.win.ShouldClose() {
		start := time.Now()

		w.win.MakeContextCurrent()
		windowWidth, windowHeight := w.win.GetSize()

		//total := time.Now()

		ui.clear()

		//tUserRender := time.Now()
		if w.userRender != nil {
			w.userRender(ui)
		}
		//fmt.Println("w.userRender(ui): ", time.Now().Sub(tUserRender))

		//tUIRender := time.Now()
		buffer := ui.render(windowWidth, windowHeight)
		//fmt.Println("ui.render(...): ", time.Now().Sub(tUIRender))

		bounds := buffer.Bounds()
		bufferWidth, bufferHeight := bounds.Dx(), bounds.Dy()

		//tGLRender := time.Now()

		//gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.RasterPos2f(-1, 1)
		gl.PixelZoom(1, -1)
		gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
		gl.DrawPixels(
			int32(bufferWidth), int32(bufferHeight),
			gl.RGBA, gl.UNSIGNED_BYTE,
			unsafe.Pointer(&buffer.Pix[0]))
		//gl.Flush()
		//fmt.Println("gl.Render(...): ", time.Now().Sub(tGLRender))

		w.win.SwapBuffers()

		//endTotal := time.Now().Sub(total)
		//totalForAvr += endTotal
		//runs++

		//fmt.Println("Total: ", endTotal)
		//fmt.Println("Avg: ", totalForAvr/time.Duration(runs))

		glfw.PollEvents()

		diff := time.Now().Sub(start)

		if diff < 30*time.Millisecond {
			time.Sleep(30*time.Millisecond - diff)
		}
		//fmt.Println()
	}
}
