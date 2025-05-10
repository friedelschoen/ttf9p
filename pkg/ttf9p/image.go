package ttf9p

import (
	"fmt"
	"image"
	"io"
	"slices"
)

func WriteImage(subfont io.Writer, img *image.Gray) {
	width, height := img.Rect.Dx(), img.Rect.Dy()
	pix := img.Pix

	format := "k8"
	hasgray := slices.ContainsFunc(pix, func(p byte) bool {
		return p != 0x00 && p != 0xff
	})
	if !hasgray {
		format = "k1"
		newpix := make([]byte, 0, len(pix))
		for y := 0; y < height; y++ {
			b := y * width
			for x := 0; x < width; x += 8 {
				t := 0
				for i := 0; i < 8; i++ {
					t <<= 1
					if x+i < width {
						t |= int(pix[b] & 1)
						b++
					}
				}
				newpix = append(newpix, byte(t))
			}
		}
		pix = newpix
	}

	// final fontchar
	fmt.Fprintf(subfont, "%11s %11d %11d %11d %11d ", format, 0, 0, width, height)
	subfont.Write(pix)
}
