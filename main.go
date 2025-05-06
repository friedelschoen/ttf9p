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
		fmt.Fprintf(os.Stderr, "invalid hinting")
		flag.Usage()
		os.Exit(1)
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
	if len(sizes) == 0 {
		fmt.Fprintf(os.Stderr, "missing font-sizes")
		flag.Usage()
		os.Exit(1)
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "missing output-prefix")
		flag.Usage()
		os.Exit(1)
	}
	prefix := args[0]
	args = args[1:]
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "missing input-font")
		flag.Usage()
		os.Exit(1)
	}

	err := os.MkdirAll(path.Dir(prefix), 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create directory: %v", err)
		os.Exit(1)
	}

	for _, sz := range sizes {
		err := writeFont(prefix, sz, *dpi, hint, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to write font: %v", err)
			os.Exit(1)
		}
	}
}
