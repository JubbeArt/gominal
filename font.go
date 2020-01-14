package gominal

import (
	"io/ioutil"
	"os"

	"github.com/flopp/go-findfont"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var potentialFonts = []string{
	"DejaVuSansMono.ttf",
	"NotoMono-Regular.ttf",
	"UbuntuMono-R.ttf",
	"LiberationMono-Regular.ttf",
	"FreeMono.ttf",
}

func loadFont(size int) (font.Face, error) {
	var fontPath = ""
	var err error

	for _, fontName := range potentialFonts {
		fontPath, err = findfont.Find(fontName)

		if err != nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	file, err := os.Open(fontPath)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	ttFont, err := truetype.Parse(bytes)

	if err != nil {
		return nil, err
	}

	fontFace := truetype.NewFace(ttFont, &truetype.Options{
		Size:              float64(size),
		GlyphCacheEntries: 2048,
	})
	return fontFace, nil
}
