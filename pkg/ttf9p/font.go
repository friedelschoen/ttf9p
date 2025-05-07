package ttf9p

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

func WriteFont(prefix string, ptsz, dpi int, hint font.Hinting, inputs []string) error {
	s := fmt.Sprintf("%s.%d.font", prefix, ptsz)

	fdfont, err := os.Create(s)
	if err != nil {
		return err
	}
	defer fdfont.Close()

	for i, ofile := range inputs {
		content, err := os.ReadFile(ofile)
		if err != nil {
			return err
		}
		fontfile, err := opentype.Parse(content)
		if err != nil {
			return err
		}

		f, err := opentype.NewFace(fontfile, &opentype.FaceOptions{
			Hinting: hint,
			Size:    float64(ptsz),
			DPI:     float64(dpi),
		})
		if err != nil {
			return err
		}

		/* is header */
		if i == 0 {
			fmt.Fprintf(fdfont, "%-4d %d\n", f.Metrics().Height.Ceil(), f.Metrics().Ascent.Ceil())
		}

		ranges, err := GetCharset(fontfile)
		if err != nil {
			return err
		}
		for _, rn := range ranges {
			var w fixed.Int26_6
			start := rn.Min

			for r := rn.Min; r <= rn.Max; r++ {
				advance, _ := f.GlyphAdvance(r)

				if w+advance > Maxsubfwidth {
					if start < r {
						err := writeSubfont(fdfont, prefix, ptsz, f, Range{start, r - 1}, w.Ceil())
						if err != nil {
							return err
						}
					}
					start = r
					w = 0
				}
				w += advance
			}

			if w > 0 {
				err := writeSubfont(fdfont, prefix, ptsz, f, Range{start, rn.Max}, w.Round())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
