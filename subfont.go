package main

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"
	"path"
	"slices"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func toSubfont(fdfont io.Writer, opath string, f font.Face, rn Range, width int) {
	height := f.Metrics().Height.Ceil()
	length := rn.Max - rn.Min + 1

	fmt.Fprintf(fdfont, "0x%X\t0x%X\t%s.%X-%X\n", rn.Min, rn.Max, path.Base(opath), rn.Min, rn.Max)

	img := image.NewGray(image.Rectangle{Max: image.Point{width, height}})

	fcs := make([]Fontchar, length)
	dot := fixed.Point26_6{Y: f.Metrics().Ascent}
	for i := range length {
		bounds, advance, _ := f.GlyphBounds(rune(i) + rn.Min)
		advance = max(advance, bounds.Max.X-bounds.Min.X)

		dr, mask, maskp, _, _ := f.Glyph(dot, rune(i)+rn.Min)
		if !dr.Empty() {
			draw.DrawMask(img, dr, image.White, image.Point{}, mask, maskp, draw.Src)
		}
		fcs[i] = Fontchar{
			X:      dot.X.Round(),
			Top:    uint8(dr.Min.Y),
			Left:   0,
			Bottom: uint8(dr.Max.Y),
			Width:  uint8(advance.Round()),
		}

		dot.X += advance
	}

	pat := fmt.Sprintf("%s.%X-%X", opath, rn.Min, rn.Max)

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

	subfont, err := os.Create(pat)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// final fontchar
	fmt.Fprintf(subfont, "%11s %11d %11d %11d %11d ", format, 0, 0, width, height)
	subfont.Write(pix)
	fmt.Fprintf(subfont, "%11d %11d %11d ", len(fcs), height, f.Metrics().Ascent.Round())
	for _, c := range fcs {
		subfont.Write(c.Encode())
	}
	subfont.Write(Fontchar{X: width}.Encode())
}
