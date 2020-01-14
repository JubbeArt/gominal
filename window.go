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

type Window struct {
	win        *glfw.Window
	userRender func(ui *UI)

	width     int
	height    int
	colWidth  int
	rowHeight int

	resizeCallback func(width, height int)
	charCallback   func(char rune)
	keyCallback    func(event KeyEvent)
}

type KeyEvent struct {
	Key   glfw.Key
	Ctrl  bool
	Shift bool
}

func NewWindow(render func(ui *UI)) *Window {
	win := &Window{
		win:        nil,
		userRender: render,
		width:      1200,
		height:     800,
		colWidth:   12,
		rowHeight:  24,

		resizeCallback: func(width, height int) {},
		charCallback:   func(char rune) {},
		keyCallback:    func(event KeyEvent) {},
	}
	return win
}

func (w *Window) Run() {
	err := glfw.Init()

	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	w.win, err = glfw.CreateWindow(w.width, w.height, "", nil, nil)
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

	font, err := loadFont(18)

	if err != nil {
		panic("could not find a suitable font")
	}
	defer font.Close()

	//totalForAvr := time.Duration(0)
	//runs := 0

	for !w.win.ShouldClose() {
		start := time.Now()

		w.win.MakeContextCurrent()
		windowWidth, windowHeight := w.win.GetSize()

		//total := time.Now()
		ui := &UI{
			win:       w.win,
			colWidth:  w.colWidth,
			rowHeight: w.rowHeight,
			font:      font,
		}
		ui.setup()

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

		//gl.setup(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
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

func (w *Window) SetTitle(title string) *Window {
	w.SetTitle(title)
	return w
}

func (w *Window) Size(cols, rows int) *Window {
	w.width = cols * w.colWidth
	w.height = rows * w.rowHeight
	return w
}

func (w *Window) SizeInPixel(width, height int) *Window {
	w.width = width
	w.height = height
	return w
}

func (w *Window) OnResize(handler func(width, height int)) *Window {
	w.resizeCallback = handler
	return w
}

func (w *Window) OnKeyDown(handler func(event KeyEvent)) *Window {
	w.keyCallback = handler
	return w
}

func (w *Window) OnChar(handler func(char rune)) *Window {
	w.charCallback = handler
	return w
}
