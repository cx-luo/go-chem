// Package render2d provides a small Go-native molecule renderer inspired by
// Indigo's core/render2d module.
package render2d

import "image/color"

// Point is a screen-space 2D coordinate.
type Point struct {
	X float64
	Y float64
}

// Options controls molecule rendering.
type Options struct {
	Width  int
	Height int
	Margin float64

	BackgroundColor   string
	BondColor         string
	AromaticBondColor string
	AtomOutlineColor  string

	BondLineWidth float64
	AtomRadius    float64
	FontSize      float64

	ShowCarbonLabels         bool
	ShowTerminalCarbonLabels bool
	UseAtomColors            bool
}

// DefaultOptions returns conservative defaults for a single molecule drawing.
func DefaultOptions() Options {
	return Options{
		Width:                    300,
		Height:                   300,
		Margin:                   24,
		BackgroundColor:          "#FFFFFF",
		BondColor:                "#333333",
		AromaticBondColor:        "#AA7733",
		AtomOutlineColor:         "#222222",
		BondLineWidth:            2,
		AtomRadius:               6,
		FontSize:                 14,
		ShowCarbonLabels:         false,
		ShowTerminalCarbonLabels: false,
		UseAtomColors:            true,
	}
}

func normalizeOptions(options []Options) Options {
	opt := DefaultOptions()
	if len(options) > 0 {
		user := options[0]
		if user.Width != 0 {
			opt.Width = user.Width
		}
		if user.Height != 0 {
			opt.Height = user.Height
		}
		if user.Margin != 0 {
			opt.Margin = user.Margin
		}
		if user.BackgroundColor != "" {
			opt.BackgroundColor = user.BackgroundColor
		}
		if user.BondColor != "" {
			opt.BondColor = user.BondColor
		}
		if user.AromaticBondColor != "" {
			opt.AromaticBondColor = user.AromaticBondColor
		}
		if user.AtomOutlineColor != "" {
			opt.AtomOutlineColor = user.AtomOutlineColor
		}
		if user.BondLineWidth != 0 {
			opt.BondLineWidth = user.BondLineWidth
		}
		if user.AtomRadius != 0 {
			opt.AtomRadius = user.AtomRadius
		}
		if user.FontSize != 0 {
			opt.FontSize = user.FontSize
		}
		opt.ShowCarbonLabels = user.ShowCarbonLabels
		opt.ShowTerminalCarbonLabels = user.ShowTerminalCarbonLabels
		if user.UseAtomColors {
			opt.UseAtomColors = true
		}
	}
	if opt.Width <= 0 {
		opt.Width = 300
	}
	if opt.Height <= 0 {
		opt.Height = 300
	}
	if opt.Margin < 0 {
		opt.Margin = 0
	}
	return opt
}

func parseHexColor(s string) color.RGBA {
	if len(s) != 7 || s[0] != '#' {
		return color.RGBA{R: 136, G: 136, B: 136, A: 255}
	}
	return color.RGBA{
		R: hexByte(s[1], s[2]),
		G: hexByte(s[3], s[4]),
		B: hexByte(s[5], s[6]),
		A: 255,
	}
}

func hexByte(h, l byte) uint8 {
	return (hexNibble(h) << 4) | hexNibble(l)
}

func hexNibble(c byte) uint8 {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	default:
		return 0
	}
}
