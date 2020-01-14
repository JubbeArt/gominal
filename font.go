package gominal

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func loadFonts(size int) (font.Face, font.Face) {
	fontNormal, _ := truetype.Parse(fontBytes)
	fontBold, _ := truetype.Parse(fontBoldBytes)

	ffNormal := truetype.NewFace(fontNormal, &truetype.Options{
		Size:              float64(size),
		GlyphCacheEntries: 2048,
	})

	ffBold := truetype.NewFace(fontBold, &truetype.Options{
		Size:              float64(size),
		GlyphCacheEntries: 2048,
	})

	return ffNormal, ffBold
}
