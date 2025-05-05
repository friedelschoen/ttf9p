# ttfs

**ttfs** is a tool that converts TTF/OTF fonts into Plan 9 subfont files, compatible with software using the Plan 9 font format.

It parses TrueType or OpenType fonts, extracts supported Unicode glyph ranges, renders them to bitmaps, and writes `.font` metadata and subfont files (`*.XXXX-YYYY`) using the Plan 9 format (either `k8` or `k1` encoding depending on grayscale usage).

This project is based on [`ttfs` by ftrvxmtrx](https://git.sr.ht/~ft/ttfs/tree) which was written using C and SDL2.

## Features

- Supports TrueType and OpenType fonts.
- Outputs a `.font` file and a set of subfont image files.
- Implements `k1` fallback when no grayscale is used.
- Lightweight and fast!

## Installation

```sh
git clone https://github.com/friedelschoen/ttfs
cd ttfs
go build
````

## Usage

```sh
./ttfs [options] font1.ttf [font2.ttf ...] myfont/myfont.16
```

Example:

```sh
./ttfs -p 14 -H none DejaVuSansMono.ttf /lib/font/bit/dejavu/dejavu.14
```

This will generate:

* `dejavu/dejavu.14.font` metadata file
* one or more subfont image files like `dejavu/dejavu.14.0020-007F`, `dejavu/dejavu.14.00A0-00FF`, etc.

## Notes

* The maximum width per subfont is capped to ensure compatibility.
* Subfonts are written using grayscale (`k8`) unless the image only uses black/white pixels, in which case `k1` is used.

## License

Zlib-style license. See [LICENSE](LICENSE).
