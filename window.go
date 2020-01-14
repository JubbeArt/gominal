package gominal

import (
	"os"
	"os/signal"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/image/font"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

type Window struct {
	win        *glfw.Window
	userRender func(ui *UI, cols, rows int)

	colWidth  int
	rowHeight int

	fontNormal font.Face
	fontBold   font.Face

	resizeCallback func(width, height int)
	charCallback   func(char rune)
	keyCallback    func(event KeyEvent)
}

type KeyEvent struct {
	Key   glfw.Key
	Ctrl  bool
	Shift bool
}

func NewWindow() (*Window, error) {
	w := &Window{
		win:        nil,
		userRender: func(ui *UI, cols, rows int) {},
		colWidth:   12,
		rowHeight:  24,

		resizeCallback: func(width, height int) {},
		charCallback:   func(char rune) {},
		keyCallback:    func(event KeyEvent) {},
	}

	err := glfw.Init()

	if err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	width := 1200
	height := 800

	w.win, err = glfw.CreateWindow(width, height, "", nil, nil)
	if err != nil {
		return nil, err
	}

	w.win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, err
	}

	w.fontNormal, w.fontBold = loadFonts(18)

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

	return w, nil
}

func (w *Window) Run() {
	defer glfw.Terminate()
	//totalForAvr := time.Duration(0)
	//runs := 0

	interruptSignal := make(chan os.Signal)
	signal.Notify(interruptSignal, os.Interrupt)

	for !w.win.ShouldClose() {
		select {
		case <-interruptSignal:
			w.Close()
		default:
		}

		start := time.Now()

		w.win.MakeContextCurrent()
		windowWidth, windowHeight := w.win.GetSize()

		//total := time.Now()
		ui := &UI{
			win:        w.win,
			colWidth:   w.colWidth,
			rowHeight:  w.rowHeight,
			fontNormal: w.fontNormal,
			fontBold:   w.fontBold,
		}
		ui.setup()

		//tUserRender := time.Now()
		if w.userRender != nil {
			w.userRender(ui, windowWidth/w.colWidth, windowHeight/w.rowHeight)
		}
		//fmt.Println("w.userRender(ui): ", time.Now().Sub(tUserRender))

		//tUIRender := time.Now()
		img := ui.render(windowWidth, windowHeight)
		//fmt.Println("ui.render(...): ", time.Now().Sub(tUIRender))

		bounds := img.Bounds()
		imgWidth, imgHeight := bounds.Dx(), bounds.Dy()

		//tGLRender := time.Now()

		//gl.setup(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.RasterPos2f(-1, 1)
		gl.PixelZoom(1, -1)
		gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
		gl.DrawPixels(
			int32(imgWidth), int32(imgHeight),
			gl.RGBA, gl.UNSIGNED_BYTE,
			unsafe.Pointer(&img.Pix[0]))
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

func (w *Window) Rows() int {
	_, height := w.win.GetSize()
	return height / w.rowHeight
}
func (w *Window) Cols() int {
	width, _ := w.win.GetSize()
	return width / w.colWidth
}
func (w *Window) GlfwWindow() *glfw.Window                   { return w.win }
func (w *Window) Render(render func(ui *UI, cols, rows int)) { w.userRender = render }
func (w *Window) SetSize(cols, rows int)                     { w.win.SetSize(cols*w.colWidth, rows*w.rowHeight) }
func (w *Window) SetTitle(title string)                      { w.win.SetTitle(title) }
func (w *Window) SetSizeInPixel(width, height int)           { w.win.SetSize(width, height) }
func (w *Window) OnResize(handler func(width, height int))   { w.resizeCallback = handler }
func (w *Window) OnKeyDown(handler func(event KeyEvent))     { w.keyCallback = handler }
func (w *Window) OnChar(handler func(char rune))             { w.charCallback = handler }
func (w *Window) Close()                                     { w.win.SetShouldClose(true) }
