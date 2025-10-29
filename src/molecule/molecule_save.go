// Package src coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 16:25
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_save.go
// @Software: GoLand
package src

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

// SavePNG renders the molecule to a PNG file with a simple circular layout.
func (m *Molecule) SavePNG(filename string, size int) error {
	img := renderMoleculeRaster(m, size, size)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// SaveJPEG renders the molecule to a JPEG file with given quality (1..100).
func (m *Molecule) SaveJPEG(filename string, size int, quality int) error {
	if quality <= 0 {
		quality = 85
	}
	img := renderMoleculeRaster(m, size, size)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	opts := &jpeg.Options{Quality: quality}
	return jpeg.Encode(f, img, opts)
}

// SaveSVG writes a simple SVG depiction (lines for bonds, circles for atoms).
func (m *Molecule) SaveSVG(filename string, width int, height int) error {
	coords := computeCircularLayout(m, width, height, 20)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// header
	if _, err = f.WriteString(`<?xml version="1.0" encoding="UTF-8"?>\n`); err != nil {
		return err
	}
	if _, err = f.WriteString(svgf("<svg xmlns='http://www.w3.org/2000/svg' width='%d' height='%d' viewBox='0 0 %d %d'>\n", width, height, width, height)); err != nil {
		return err
	}
	// bonds
	for i, e := range m.Bonds {
		if e.Order < 0 {
			continue
		}
		c1 := coords[e.Beg]
		c2 := coords[e.End]
		stroke := "#444"
		widthPx := 2
		if m.BondOrders[i] == BOND_DOUBLE {
			widthPx = 4
		} else if m.BondOrders[i] == BOND_TRIPLE {
			widthPx = 6
		} else if m.BondOrders[i] == BOND_AROMATIC {
			stroke = "#AA7733"
		}
		if _, err = f.WriteString(svgf("<line x1='%.1f' y1='%.1f' x2='%.1f' y2='%.1f' stroke='%s' stroke-width='%d' stroke-linecap='round'/>\n", c1.X, c1.Y, c2.X, c2.Y, stroke, widthPx)); err != nil {
			return err
		}
	}
	// atoms
	for i := range m.Atoms {
		c := coords[i]
		r := 6.0
		fill := elementColor(m.Atoms[i].Number)
		if _, err = f.WriteString(svgf("<circle cx='%.1f' cy='%.1f' r='%.1f' fill='%s' stroke='#222' stroke-width='1'/>\n", c.X, c.Y, r, fill)); err != nil {
			return err
		}
	}
	if _, err = f.WriteString("</svg>\n"); err != nil {
		return err
	}
	return nil
}

// --- internal rendering helpers ---

type point struct{ X, Y float64 }

func renderMoleculeRaster(m *Molecule, width, height int) image.Image {
	coords := computeCircularLayout(m, width, height, 20)
	bg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(bg, bg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// draw bonds
	for i, e := range m.Bonds {
		if e.Order < 0 {
			continue
		}
		c1 := coords[e.Beg]
		c2 := coords[e.End]
		stroke := color.RGBA{R: 68, G: 68, B: 68, A: 255}
		w := 2
		if m.BondOrders[i] == BOND_DOUBLE {
			w = 4
		} else if m.BondOrders[i] == BOND_TRIPLE {
			w = 6
		} else if m.BondOrders[i] == BOND_AROMATIC {
			stroke = color.RGBA{R: 170, G: 119, B: 51, A: 255}
		}
		drawLine(bg, int(c1.X), int(c1.Y), int(c2.X), int(c2.Y), stroke, w)
	}

	// draw atoms
	for i := range m.Atoms {
		c := coords[i]
		r := 6
		fill := parseHexColor(elementColor(m.Atoms[i].Number))
		drawFilledCircle(bg, int(c.X), int(c.Y), r, fill)
		drawCircle(bg, int(c.X), int(c.Y), r, color.RGBA{A: 255})
	}
	return bg
}

func computeCircularLayout(m *Molecule, width, height int, margin int) []point {
	n := len(m.Atoms)
	coords := make([]point, n)
	if n == 0 {
		return coords
	}
	cx := float64(width / 2)
	cy := float64(height / 2)
	r := math.Min(float64(width), float64(height))/2 - float64(margin)
	if r < 10 {
		r = 10
	}
	for i := 0; i < n; i++ {
		angle := 2 * math.Pi * float64(i) / float64(n)
		x := cx + r*math.Cos(angle)
		y := cy + r*math.Sin(angle)
		coords[i] = point{X: x, Y: y}
	}
	return coords
}

// primitive raster drawing
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.Color, width int) {
	// simple Bresenham with thickness by drawing multiple offsets
	setPix := func(x, y int) {
		if x >= 0 && y >= 0 && x < img.Bounds().Dx() && y < img.Bounds().Dy() {
			img.Set(x, y, col)
		}
	}
	dx := int(math.Abs(float64(x1 - x0)))
	dy := -int(math.Abs(float64(y1 - y0)))
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		for ox := -width / 2; ox <= width/2; ox++ {
			for oy := -width / 2; oy <= width/2; oy++ {
				setPix(x0+ox, y0+oy)
			}
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func drawCircle(img *image.RGBA, cx, cy, r int, col color.Color) {
	x := r
	y := 0
	set := func(x, y int) {
		if x >= 0 && y >= 0 && x < img.Bounds().Dx() && y < img.Bounds().Dy() {
			img.Set(x, y, col)
		}
	}
	for x >= y {
		set(cx+x, cy+y)
		set(cx+y, cy+x)
		set(cx-y, cy+x)
		set(cx-x, cy+y)
		set(cx-x, cy-y)
		set(cx-y, cy-x)
		set(cx+y, cy-x)
		set(cx+x, cy-y)
		y++
		if delta := 2*y + 1 - 2*x; delta > 0 {
			x--
		}
	}
}

func drawFilledCircle(img *image.RGBA, cx, cy, r int, col color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				xp := cx + x
				yp := cy + y
				if xp >= 0 && yp >= 0 && xp < img.Bounds().Dx() && yp < img.Bounds().Dy() {
					img.Set(xp, yp, col)
				}
			}
		}
	}
}

func elementColor(z int) string {
	switch z {
	case ELEM_H:
		return "#BBBBBB"
	case ELEM_C:
		return "#222222"
	case ELEM_N:
		return "#3366CC"
	case ELEM_O:
		return "#CC3333"
	case ELEM_S:
		return "#CCCC33"
	case ELEM_F, ELEM_Cl, ELEM_Br, ELEM_I:
		return "#33AA66"
	default:
		return "#888888"
	}
}

// tiny helpers to avoid extra imports
func parseHexColor(s string) color.RGBA {
	// expecting #RRGGBB
	r := uint8(0x88)
	g := uint8(0x88)
	b := uint8(0x88)
	if len(s) == 7 && s[0] == '#' {
		r = hexByte(s[1], s[2])
		g = hexByte(s[3], s[4])
		b = hexByte(s[5], s[6])
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func hexNibble(c byte) uint8 {
	if c >= '0' && c <= '9' {
		return uint8(c - '0')
	}
	if c >= 'a' && c <= 'f' {
		return uint8(c-'a') + 10
	}
	if c >= 'A' && c <= 'F' {
		return uint8(c-'A') + 10
	}
	return 0
}

func hexByte(h, l byte) uint8 { return (hexNibble(h) << 4) | hexNibble(l) }

func svgf(format string, a ...any) string {
	return sprintf(format, a...)
}

// minimal fmt.Sprintf replacement to reduce imports in this file
func sprintf(format string, a ...any) string {
	return fmtSprintf(format, a...)
}

// use standard fmt for actual formatting
// kept separate to highlight limited use
func fmtSprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
