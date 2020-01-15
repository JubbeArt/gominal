package main

//type FontStyle int
//
//const (
//	FontNormal FontStyle = iota
//	FontBold
//)
//
//type UI struct {
//	runes           [][]rune
//	images          [][]image.Image
//	isImage         [][]bool
//	textColors      [][]color.Color
//	textBackgrounds [][]color.Color
//	textStyles      [][]FontStyle
//
//	colWidth  int
//	rowHeight int
//
//	cols int
//	rows int
//
//	fontBold   font.Face
//	fontNormal font.Face
//}
//
//func (ui *UI) setup(width, height int) {
//	ui.cols = width / ui.colWidth
//	ui.rows = height / ui.rowHeight
//
//	ui.runes = make([][]rune, ui.rows)
//	ui.textColors = make([][]color.Color, ui.rows)
//	ui.textBackgrounds = make([][]color.Color, ui.rows)
//	ui.textStyles = make([][]FontStyle, ui.rows)
//	ui.isImage = make([][]bool, ui.rows)
//	ui.images = make([][]image.Image, ui.rows)
//
//	for row := range ui.runes {
//		ui.runes[row] = make([]rune, ui.cols)
//		ui.textColors[row] = make([]color.Color, ui.cols)
//		ui.textBackgrounds[row] = make([]color.Color, ui.cols)
//		ui.textStyles[row] = make([]FontStyle, ui.cols)
//		ui.isImage[row] = make([]bool, ui.cols)
//		ui.images[row] = make([]image.Image, ui.cols)
//	}
//}
//
//func (ui *UI) Rune(char rune, col, row int) {
//	ui.RuneFull(char, col, row, color.White, color.Black, FontNormal)
//}
//
//func (ui *UI) RuneFull(char rune, col, row int, textColor color.Color, bgColor color.Color, fontStyle FontStyle) {
//	if ui.outOfBounds(col, row) {
//		return
//	}
//
//	ui.runes[row][col] = char
//	ui.isImage[row][col] = false
//	ui.textColors[row][col] = textColor
//	ui.textBackgrounds[row][col] = bgColor
//	ui.textStyles[row][col] = fontStyle
//}
//
//func (ui *UI) Image(image image.Image, col, row int) {
//	if ui.outOfBounds(col, row) {
//		return
//	}
//
//	ui.images[row][col] = image
//	ui.isImage[row][col] = true
//}
//
//func (ui *UI) Rows() int { return ui.rows }
//func (ui *UI) Cols() int { return ui.cols }
//
//func (ui *UI) FullImage(img image.Image, startCol, startRow int) {
//	bounds := img.Bounds()
//	width, height := bounds.Dx(), bounds.Dy()
//
//	row := startRow
//
//	var wg sync.WaitGroup
//
//	for y := 0; y < height+ui.rowHeight; y += ui.rowHeight {
//		col := startCol
//
//		for x := 0; x < width+ui.colWidth; x += ui.colWidth {
//			if ui.outOfBounds(col, row) {
//				break
//			}
//
//			wg.Add(1)
//			x := x
//			y := y
//			c := col
//			row := row
//
//			go func() {
//				subImg := image.NewRGBA(image.Rect(0, 0, ui.colWidth, ui.rowHeight))
//
//				start := image.Point{X: x, Y: y}
//				draw.Draw(subImg, subImg.Bounds(), img, start, draw.Src)
//				ui.Image(subImg, c, row)
//				wg.Done()
//			}()
//			col++
//		}
//
//		row++
//	}
//
//	wg.Wait()
//}
//
//func (ui *UI) render(width, height int) *image.RGBA {
//	outImg := image.NewRGBA(image.Rect(0, 0, width, height))
//	textImg := image.NewRGBA(outImg.Bounds())
//
//	normalDrawer := font.Drawer{
//		Dst:  textImg,
//		Src:  image.White,
//		Face: ui.fontNormal,
//	}
//
//	boldDrawer := font.Drawer{
//		Dst:  textImg,
//		Src:  image.White,
//		Face: ui.fontBold,
//	}
//
//	var wg sync.WaitGroup
//
//	for row := 0; row < ui.rows; row++ {
//		for col := 0; col < ui.cols; col++ {
//			if ui.isImage[row][col] && ui.images[row][col] != nil {
//				wg.Add(1)
//				col := col
//				row := row
//
//				go func() {
//					img := ui.images[row][col]
//					draw.Draw(outImg, ui.rect(col, row), img, image.Point{}, draw.Over)
//					wg.Done()
//				}()
//			} else if !ui.isImage[row][col] && ui.runes[row][col] != 0 {
//				char := ui.runes[row][col]
//				textDrawer := normalDrawer
//
//				if ui.textStyles[row][col] == FontBold {
//					textDrawer = boldDrawer
//				}
//
//				textDrawer.Dot = fixed.P(col*ui.colWidth, (row+1)*ui.rowHeight-3)
//				textDrawer.DrawString(string(char))
//
//				wg.Add(1)
//				col := col
//				row := row
//				go func() {
//					bgColor := ui.textBackgrounds[row][col]
//
//					if bgColor != color.Black {
//						draw.Draw(outImg, ui.rect(col, row), image.NewUniform(bgColor), image.Point{}, draw.Src)
//					}
//
//					textColor := ui.textColors[row][col]
//					textStart := image.Point{X: col * ui.colWidth, Y: row * ui.rowHeight}
//
//					if textColor == color.White {
//						draw.Draw(
//							outImg, ui.rect(col, row),
//							textImg, textStart,
//							draw.Over)
//					} else {
//						draw.DrawMask(
//							outImg, ui.rect(col, row),
//							image.NewUniform(textColor), image.Point{},
//							textImg, textStart,
//							draw.Over)
//					}
//
//					wg.Done()
//				}()
//			}
//		}
//	}
//	wg.Wait()
//	return outImg
//}
//
//func (ui *UI) rect(col, row int) image.Rectangle {
//	return image.Rect(col*ui.colWidth, row*ui.rowHeight, (col+1)*ui.colWidth, (row+1)*ui.rowHeight)
//}
//
//func (ui *UI) outOfBounds(col, row int) bool {
//	return col >= ui.cols || row >= ui.rows
//}
