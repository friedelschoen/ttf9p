package ttf9p

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"
	"path"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Fontchar struct {
	X      int
	Top    uint8
	Bottom uint8
	Left   uint8
	Width  uint8
}

func (c Fontchar) Encode() []byte {
	return []byte{
		byte(c.X >> 0),
		byte(c.X >> 8),
		c.Top,
		c.Bottom,
		c.Left,
		c.Width,
	}
}

func writeSubfont(fdfont io.Writer, prefix string, ptsz int, f font.Face, rn Range, width int) error {
	height := f.Metrics().Height.Ceil()
	length := rn.Max - rn.Min + 1

	fmt.Fprintf(fdfont, "0x%X\t0x%X\t%s.%d.%X-%X\n", rn.Min, rn.Max, path.Base(prefix), ptsz, rn.Min, rn.Max)

	img := image.NewGray(image.Rectangle{Max: image.Point{width, height}})

	fcs := make([]Fontchar, length)
	dot := fixed.Point26_6{Y: f.Metrics().Ascent}
	for i := range length {
		dr, mask, maskp, advance, _ := f.Glyph(dot, rune(i)+rn.Min)
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

	pat := fmt.Sprintf("%s.%d.%X-%X", prefix, ptsz, rn.Min, rn.Max)

	subfont, err := os.Create(pat)
	if err != nil {
		return err
	}
	defer f.Close()
	WriteImage(subfont, img)

	fmt.Fprintf(subfont, "%11d %11d %11d ", len(fcs), height, f.Metrics().Ascent.Round())
	for _, c := range fcs {
		subfont.Write(c.Encode())
	}
	subfont.Write(Fontchar{X: width}.Encode())
	return nil
}
