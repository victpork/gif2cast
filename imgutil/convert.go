package imgutil

import (
	"image"
	"image/color"
	"image/gif"

	"github.com/bamiaux/rez"
)

// ConvertPalettedImgToRGBA converts a paletted image format (e.g. GIF)
// into RGBA
func ConvertPalettedImgToRGBA(src *image.Paletted) (dst *image.RGBA) {
	dst = &image.RGBA{
		Pix:    make([]uint8, 0, len(src.Pix)*4),
		Rect:   src.Rect,
		Stride: src.Stride * 4,
	}

	for i := range src.Pix {
		c := src.Palette[src.Pix[i]].(color.RGBA)
		dst.Pix = append(dst.Pix, c.R, c.G, c.B, c.A)
	}
	return
}

// Resize takes a GIF image and resize it with given width and height
// Output is RGBA format in a stack.
func Resize(gifImg *gif.GIF, w, h int) (dst []*image.RGBA, err error) {
	dst = make([]*image.RGBA, 0, len(gifImg.Image))

	converter, err := rez.NewConverter(&rez.ConverterConfig{
		Input: rez.Descriptor{
			Width:      gifImg.Image[0].Rect.Dx(),
			Height:     gifImg.Image[0].Rect.Dy(),
			Interlaced: false,
			Ratio:      rez.Ratio444,
			Pack:       4,
			Planes:     1,
		},
		Output: rez.Descriptor{
			Width:      w,
			Height:     h,
			Interlaced: false,
			Ratio:      rez.Ratio444,
			Pack:       4,
			Planes:     1,
		},
	}, rez.NewBilinearFilter())
	if err != nil {
		return nil, err
	}
	for i := range gifImg.Image {
		dstFrame := image.NewRGBA(image.Rect(0, 0, w, h))
		err = converter.Convert(dstFrame, ConvertPalettedImgToRGBA(gifImg.Image[i]))
		if err != nil {
			return nil, err
		}
		dst = append(dst, dstFrame)
	}
	return dst, nil
}
