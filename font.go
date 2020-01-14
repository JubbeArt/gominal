package gominal

import (
	"io/ioutil"
	"os"

	"github.com/flopp/go-findfont"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func loadFont(name string, size int) (font.Face, error) {

	fontPath, err := findfont.Find(name)
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
