package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/spf13/pflag"
	"golang.org/x/image/font"
)

func main() {
	pflag.Usage = func() {
		prog := path.Base(os.Args[0])
		out := pflag.CommandLine.Output()
		help := `Usage:
  %[1]s [flags] <size>... <prefix> <fontfile>...

Arguments:
  <size>        One or more point sizes to render (e.g. 12 16 24)
  <prefix>      Output prefix for generated files
  <fontfile>    One or more TTF/OTF font files to include

Flags:
`
		out.Write([]byte(fmt.Sprintf(help, prog)))
		pflag.PrintDefaults()
		example := `
Example:
  %[1]s -d 96 -H full 12 16 output fonts/DejaVuSans.ttf
  => creates output.12.font and output.16.font + subfonts from the given TTF

Hinting options:
  none      disables hinting
  full      enables full hinting
  vertical  enables vertical hinting only
`
		out.Write([]byte(fmt.Sprintf(example, prog)))
	}

	dpi := pflag.IntP("dpi", "d", 72, "dpi")
	hintstr := pflag.StringP("hinting", "H", "none", "hinting: none, full, vertical") // was normal, light, mono, none, light_subpixel
	pflag.Parse()

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

	if pflag.NArg() < 3 {
		pflag.Usage()
		os.Exit(1)
	}
	args := pflag.Args()

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
