package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"golang.org/x/image/math/fixed"

	"golang.org/x/image/font"
)

type drawRequest interface {
	draw(out *image.RGBA)
}

type charDrawRequest struct {
	char      rune
	col       int
	row       int
	textColor color.RGBA
	bg        color.Color
	style     string
}

func (req charDrawRequest) draw(out *image.RGBA) {
	img := image.NewRGBA(rect(req.col, req.row))
	draw.Draw(img, img.Bounds(), image.NewUniform(req.bg), image.Point{}, draw.Src)

	fontFace := fontNormal

	if req.style == styleBold {
		fontFace = fontBold
	}

	drawer := font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(req.textColor),
		Face: fontFace,
	}

	drawer.Dot = fixed.P(req.col*colWidth+1, (req.row+1)*rowHeight-3)
	drawer.DrawString(string(req.char))

	draw.Draw(out, img.Bounds(), img, img.Bounds().Min, draw.Src)
}

type imageDrawRequest struct {
	img image.Image
	col int
	row int
}

func (req imageDrawRequest) draw(out *image.RGBA) {
	start := image.Point{X: req.col * colWidth, Y: req.row * rowHeight}
	bounds := req.img.Bounds()
	cols := int(math.Floor(float64(bounds.Dx()) / float64(colWidth)))
	rows := int(math.Floor(float64(bounds.Dy()) / float64(rowHeight)))

	rect := image.Rect(
		start.X, start.Y,
		start.X+cols*colWidth, start.Y+rows*rowHeight)

	// should be black
	draw.Draw(out, rect, image.White, image.Point{}, draw.Src)
	draw.Draw(out, bounds.Add(start), req.img, image.Point{}, draw.Src)
}

type clearDrawRequest struct{}

func (clearDrawRequest) draw(out *image.RGBA) {
	draw.Draw(out, out.Bounds(), image.Black, image.Point{}, draw.Src)
}

func rect(col, row int) image.Rectangle {
	return image.Rect(col*colWidth, row*rowHeight, (col+1)*colWidth, (row+1)*rowHeight)
}
