package ttf9p

import (
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type Range struct {
	Min rune
	Max rune
}

func GetCharset(font *opentype.Font) ([]Range, error) {
	var buf sfnt.Buffer
	var ranges [][2]rune
	const maxRune = 0x10FFFF
	const maxrange = 16 // of Maxrange zoals in C-code

	var (
		start, end rune
		inRange    bool
	)

	for r := rune(1); r <= maxRune; r++ {
		index, err := font.GlyphIndex(&buf, r)
		if err != nil || index == 0 {
			if inRange {
				ranges = append(ranges, [2]rune{start, end})
				inRange = false
			}
			continue
		}

		if !inRange {
			start = r
			end = r
			inRange = true
		} else if r == end+1 {
			end = r
		} else {
			ranges = append(ranges, [2]rune{start, end})
			start = r
			end = r
		}
	}

	if inRange {
		ranges = append(ranges, [2]rune{start, end})
	}

	// split ranges into chunks of Maxrange
	var final []Range
	for _, r := range ranges {
		for i := r[0]; i <= r[1]; i += maxrange + 1 {
			j := i + maxrange
			if j > r[1] {
				j = r[1]
			}
			final = append(final, Range{i, j})
		}
	}

	return final, nil
}
