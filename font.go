package main

import (
	"fmt"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	Maxsubfwidth = 3000 /* rough */
)

func writeFont(prefix string, ptsz, dpi int, hint font.Hinting, inputs []string) {
	s := fmt.Sprintf("%s.%d.font", prefix, ptsz)

	fdfont, err := os.Create(s)
	if err != nil {
		panic(err)
	}
	defer fdfont.Close()

	for i, ofile := range inputs {
		content, err := os.ReadFile(ofile)
		if err != nil {
			panic(err)
		}
		fontfile, err := opentype.Parse(content)
		if err != nil {
			panic(err)
		}

		f, err := opentype.NewFace(fontfile, &opentype.FaceOptions{
			Hinting: hint,
			Size:    float64(ptsz),
			DPI:     float64(dpi),
		})
		if err != nil {
			panic(err)
		}

		if i == 0 {
			fmt.Fprintf(fdfont, "%-4d %d\n", f.Metrics().Height.Ceil(), f.Metrics().Ascent.Round())
		}

		ranges, err := GetCharset(fontfile)
		if err != nil {
			panic(err)
		}
		for _, rn := range ranges {
			var w fixed.Int26_6
			start := rn.Min

			for r := rn.Min; r <= rn.Max; r++ {
				advance, _ := f.GlyphAdvance(r)

				if w+advance > Maxsubfwidth {
					if start < r {
						writeSubfont(fdfont, prefix, ptsz, f, Range{start, r - 1}, w.Ceil())
					}
					start = r
					w = 0
				}
				w += advance
			}

			if w > 0 {
				writeSubfont(fdfont, prefix, ptsz, f, Range{start, rn.Max}, w.Round())
			}
		}
	}
}
