package gominal

import (
	"image"
	"image/color"
	"image/draw"
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type UI struct {
	runes   [][]rune
	images  [][]image.Image
	isImage [][]bool
	colors  [][]color.Color

	colWidth  int
	rowHeight int

	cols int
	rows int

	font font.Face

	win *glfw.Window
}

func newUI(win *glfw.Window, colWidth, rowHeight int) *UI {
	ui := &UI{
		win:       win,
		colWidth:  colWidth,
		rowHeight: rowHeight,
	}

	var err error
	ui.font, err = loadFont(18)

	if err != nil {
		panic("could not find a suitable font")
	}

	return ui
}

func (ui *UI) clear() {
	width, height := ui.win.GetSize()

	ui.cols = width / ui.colWidth
	ui.rows = height / ui.rowHeight

	ui.runes = make([][]rune, ui.rows)
	ui.colors = make([][]color.Color, ui.rows)
	ui.isImage = make([][]bool, ui.rows)
	ui.images = make([][]image.Image, ui.rows)

	for row := range ui.runes {
		ui.runes[row] = make([]rune, ui.cols)
		ui.colors[row] = make([]color.Color, ui.cols)
		ui.isImage[row] = make([]bool, ui.cols)
		ui.images[row] = make([]image.Image, ui.cols)
	}
}

func (ui *UI) Rune(char rune, col, row int) {
	ui.RuneWithColor(char, col, row, color.White)
}

func (ui *UI) RuneWithColor(char rune, col, row int, color color.Color) {
	if ui.outOfBounds(col, row) {
		return
	}

	ui.runes[row][col] = char
	ui.isImage[row][col] = false
	ui.colors[row][col] = color
}

func (ui *UI) Image(image image.Image, col, row int) {
	if ui.outOfBounds(col, row) {
		return
	}

	ui.images[row][col] = image
	ui.isImage[row][col] = true
}

func (ui *UI) FullImage(img image.Image, startCol, startRow int) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	row := startRow

	var wg sync.WaitGroup

	for y := 0; y < height+ui.rowHeight; y += ui.rowHeight {
		col := startCol

		for x := 0; x < width+ui.colWidth; x += ui.colWidth {
			if ui.outOfBounds(col, row) {
				break
			}

			wg.Add(1)
			x := x
			y := y
			c := col
			row := row

			go func() {
				subImg := image.NewRGBA(image.Rect(0, 0, ui.colWidth, ui.rowHeight))

				start := image.Point{X: x, Y: y}
				draw.Draw(subImg, subImg.Bounds(), img, start, draw.Src)
				ui.Image(subImg, c, row)
				wg.Done()
			}()
			col++
		}

		row++
	}

	wg.Wait()
}

func (ui *UI) render(width, height int) *image.RGBA {
	outImg := image.NewRGBA(image.Rect(0, 0, width, height))
	textImg := image.NewRGBA(outImg.Bounds())

	textDrawer := font.Drawer{
		Dst:  textImg,
		Src:  image.White,
		Face: ui.font,
	}

	var wg sync.WaitGroup

	for row := 0; row < ui.rows; row++ {
		for col := 0; col < ui.cols; col++ {
			if ui.isImage[row][col] && ui.images[row][col] != nil {
				wg.Add(1)
				col := col
				row := row

				go func() {
					img := ui.images[row][col]
					draw.Draw(outImg, ui.rect(col, row), img, image.Point{}, draw.Over)
					wg.Done()
				}()
			} else if !ui.isImage[row][col] && ui.runes[row][col] != 0 {
				char := ui.runes[row][col]
				textDrawer.Dot = fixed.P(col*ui.colWidth, (row+1)*ui.rowHeight-3)
				textDrawer.DrawString(string(char))

				wg.Add(1)
				col := col
				row := row
				go func() {
					textColor := ui.colors[row][col]
					draw.DrawMask(outImg, ui.rect(col, row),
						image.NewUniform(textColor), image.Point{},
						textImg, image.Point{X: col * ui.colWidth, Y: row * ui.rowHeight}, draw.Src)
					wg.Done()
				}()
			}
		}
	}
	wg.Wait()
	return outImg
}

func (ui *UI) rect(col, row int) image.Rectangle {
	return image.Rect(col*ui.colWidth, row*ui.rowHeight, (col+1)*ui.colWidth, (row+1)*ui.rowHeight)
}

//func (ui *UI) Cursor(cursor glfw.Cursor, col, row int) {
//	//	actually should create own types for cursor tho...
//	// user imports other glfw = unlucky
//}

func (ui *UI) outOfBounds(col, row int) bool {
	return col >= ui.cols || row >= ui.rows
}

func (ui *UI) Size(cols, rows int) {
	ui.win.SetSize(cols*ui.colWidth, rows*ui.rowHeight)
}

func (ui *UI) SetTitle(title string) {
	ui.win.SetTitle(title)
}

func (ui *UI) SizeInPixels(width, height int) {
	ui.win.SetSize(width, height)
}
