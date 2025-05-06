package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/image/font"
)

//go:embed help.txt
var helpmsg string

func main() {
	flag.Usage = func() {
		prog := path.Base(os.Args[0])
		out := flag.CommandLine.Output()

		pre, post, _ := strings.Cut(strings.ReplaceAll(helpmsg, "{}", prog), "%%")
		fmt.Fprint(out, pre)
		flag.PrintDefaults()
		fmt.Fprint(out, post)
	}

	dpi := flag.Int("dpi", 72, "dpi")
	hintstr := flag.String("hinting", "none", "hinting: none, full, vertical") // was normal, light, mono, none, light_subpixel
	flag.Parse()

	var hint font.Hinting
	switch *hintstr {
	case "none":
		hint = font.HintingNone
	case "full":
		hint = font.HintingFull
	case "vertical":
		hint = font.HintingVertical
	default:
		panic("invalid hinting")
	}

	if flag.NArg() < 3 {
		flag.Usage()
		os.Exit(1)
	}
	args := flag.Args()

	var sizes []int
	for len(args) > 0 {
		sz, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			break
		}
		sizes = append(sizes, int(sz))
		args = args[1:]
	}
	if len(args) == 0 {
		panic("no prefix")
	}
	prefix := args[0]
	args = args[1:]
	if len(args) == 0 {
		panic("no input")
	}

	os.MkdirAll(path.Dir(prefix), 0755)

	for _, sz := range sizes {
		writeFont(prefix, sz, *dpi, hint, args)
	}
}
