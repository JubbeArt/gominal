package gominal

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func loadFonts(size int) (font.Face, font.Face) {
	fontNormal, _ := truetype.Parse(fontBytes)
	fontBold, _ := truetype.Parse(fontBoldBytes)

	options := &truetype.Options{
		Size:              float64(size),
		GlyphCacheEntries: 2048,
	}

	faceNormal := truetype.NewFace(fontNormal, options)
	faceBold := truetype.NewFace(fontBold, options)

	return faceNormal, faceBold
}
