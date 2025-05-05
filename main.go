package main

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/pflag"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	Substitute   rune = 0xfffd
	Maxsubfwidth      = 3000 /* rough */
)

func main() {
	ptsz := 16
	dpi := 72
	hint := font.HintingNone

	hintstr := ""

	pflag.IntVarP(&ptsz, "point", "s", 16, "point size")
	pflag.IntVarP(&ptsz, "dpi", "d", 72, "dpi")
	pflag.StringVarP(&hintstr, "hinting", "H", "none", "hinting: none, full, vertical") // was normal, light, mono, none, light_subpixel

	pflag.Parse()

	switch hintstr {
	case "none":
		hint = font.HintingNone
	case "full":
		hint = font.HintingFull
	case "vertical":
		hint = font.HintingVertical
	default:
		panic("invalid hinting")
	}

	if pflag.NArg() < 2 {
		pflag.Usage()
		os.Exit(1)
	}
	args := pflag.Args()
	opath := fmt.Sprintf("%s.%d", args[len(args)-1], ptsz)

	s := fmt.Sprintf("%s.font", opath)

	os.MkdirAll(path.Dir(opath), 0755)

	fdfont, err := os.Create(s)
	if err != nil {
		panic(err)
	}
	defer fdfont.Close()

	for i, ofile := range args[:len(args)-1] {
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
						toSubfont(fdfont, opath, f, Range{start, r - 1}, w.Ceil())
					}
					start = r
					w = 0
				}
				w += advance
			}

			if w > 0 {
				toSubfont(fdfont, opath, f, Range{start, rn.Max}, w.Round())
			}
		}
	}
}
