# Render Examples

Examples demonstrating molecule and reaction rendering using Indigo Renderer.

## Files

- **render_examples.go** - Comprehensive rendering examples including:
  - Basic molecule rendering to PNG/SVG
  - Custom render options (size, colors, styles)
  - Reaction rendering
  - Grid rendering (multiple molecules)
  - Substructure highlighting

## Running Examples

```bash
cd examples/render
go run render_examples.go
```

## Supported Output Formats

- **PNG** - Raster image (default)
- **SVG** - Vector graphics
- **PDF** - Portable Document Format
- **EMF** - Enhanced Metafile (Windows only)

## Quick Start

```go
import (
    "github.com/cx-luo/go-chem/molecule"
    "github.com/cx-luo/go-chem/render"
)

// Initialize renderer
render.InitRenderer()
defer render.DisposeRenderer()

// Load molecule
mol, _ := molecule.LoadMoleculeFromString("CCO")
defer mol.Close()

// Render to file
render.SetRenderOption("render-output-format", "png")
render.RenderToFile(mol.Handle(), "ethanol.png")
```

## Common Render Options

```go
// Image size
render.SetRenderOptionInt("render-image-width", 800)
render.SetRenderOptionInt("render-image-height", 600)

// Background color (RGB, 0-1 range)
render.SetRenderOption("render-background-color", "1.0, 1.0, 1.0")

// Bond length in pixels
render.SetRenderOptionInt("render-bond-length", 40)

// Show atom/bond IDs
render.SetRenderOptionBool("render-atom-ids-visible", true)
render.SetRenderOptionBool("render-bond-ids-visible", true)

// Stereo style
render.SetRenderOption("render-stereo-style", "ext") // none, old, ext, bondmark

// Label mode
render.SetRenderOption("render-label-mode", "hetero") // hetero, terminal-hetero, all, none
```

## Using Render Options Struct

```go
opts := render.DefaultRenderOptions()
opts.OutputFormat = "svg"
opts.ImageWidth = 1024
opts.ImageHeight = 768
opts.BackgroundColor = "0.95, 0.95, 1.0" // Light blue
opts.ShowAtomIDs = true
opts.Apply()
```

## Grid Rendering

Render multiple molecules in a grid:

```go
// Create array
arrayHandle, _ := render.CreateArray()
defer render.FreeObject(arrayHandle)

// Add molecules
render.ArrayAdd(arrayHandle, mol1.Handle())
render.ArrayAdd(arrayHandle, mol2.Handle())
render.ArrayAdd(arrayHandle, mol3.Handle())

// Render grid
render.RenderGridToFile(arrayHandle, nil, 3, "grid.png")
```

## Notes

- Call `InitRenderer()` before any rendering operations
- Call `DisposeRenderer()` when done
- SVG format produces smaller files for web use
- PNG format is best for presentations
- PDF format is ideal for printing

