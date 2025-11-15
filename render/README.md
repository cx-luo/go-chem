# Render Package

Go package for rendering chemical structures (molecules and reactions) using the Indigo Renderer library via CGO.

## Features

- Render molecules and reactions to image files (PNG, SVG, PDF, EMF)
- Support for multiple output formats
- Flexible rendering options (size, colors, styles, etc.)
- Grid rendering for multiple molecules
- Memory buffer rendering
- Reference atom alignment

## Installation

This package requires the Indigo and Indigo-Renderer shared libraries to be available on your system.

### Prerequisites

- Indigo library (`libindigo.so`/`libindigo.dll`)
- Indigo-Renderer library (`libindigo-renderer.so`/`libindigo-renderer.dll`)
- CGO-enabled Go compiler

### Library Setup

Place the required libraries in one of the following locations:

- Windows: `3rd/windows-x86_64/` or `3rd/windows-i386/`
- Linux: `3rd/linux-x86_64/` or `3rd/linux-aarch64/`
- macOS: `3rd/darwin-x86_64/` or `3rd/darwin-aarch64/`

## Usage

### Basic Rendering

```go
package main

import (
 "github.com/cx-luo/go-indigo/molecule"
 "github.com/cx-luo/go-indigo/render"
)

func main() {
 // Load a molecule
 mol, _ := molecule.LoadMoleculeFromString("c1ccccc1")
 defer mol.Close()

 // Initialize renderer
 render.InitRenderer()
 defer render.DisposeRenderer()

 // Set output format
 render.SetRenderOption("render-output-format", "png")

 // Render to file
 render.RenderToFile(mol.Handle(), "benzene.png")
}
```

### Using Render Options

```go
// Create custom options
opts := &render.RenderOptions{
 OutputFormat:      "svg",
 ImageWidth:        800,
 ImageHeight:       600,
 BackgroundColor:   "1.0, 1.0, 1.0",
 BondLength:        40,
 RelativeThickness: 1.2,
 ShowAtomIDs:       false,
 ShowBondIDs:       false,
 StereoStyle:       "ext",
 LabelMode:         "hetero",
}

// Apply options
opts.Apply()

// Or use default options
defaultOpts := render.DefaultRenderOptions()
defaultOpts.Apply()
```

### Grid Rendering

```go
// Create molecules
mol1, _ := molecule.LoadMoleculeFromString("CCO")
mol2, _ := molecule.LoadMoleculeFromString("c1ccccc1")
mol3, _ := molecule.LoadMoleculeFromString("CC(=O)O")

// Create array
array, _ := render.CreateArray()
defer render.FreeObject(array)

// Add molecules to array
render.ArrayAdd(array, mol1.Handle())
render.ArrayAdd(array, mol2.Handle())
render.ArrayAdd(array, mol3.Handle())

// Render grid (2 columns)
render.RenderGridToFile(array, nil, 2, "molecules_grid.png")
```

### Render to Memory Buffer

```go
// Create write buffer
buffer, _ := render.CreateWriteBuffer()
defer render.FreeObject(buffer)

// Render to buffer
render.Render(mol.Handle(), buffer)

// Get data
data, _ := render.GetBufferData(buffer)
// Now you can use data ([]byte) as needed
```

### Render Reactions

```go
import "github.com/cx-luo/go-indigo/reaction"

// Load a reaction
rxn, _ := reaction.LoadReactionFromString("CCO>>CC=O")
defer rxn.Close()

// Render reaction
render.RenderToFile(rxn.Handle(), "reaction.png")
```

## Render Options

### Common Options

| Option | Values | Description |
|--------|--------|-------------|
| `render-output-format` | png, svg, pdf, emf | Output file format |
| `render-image-width` | integer | Image width in pixels |
| `render-image-height` | integer | Image height in pixels |
| `render-background-color` | "R, G, B" | Background color (0.0-1.0) |
| `render-bond-length` | integer | Bond length in pixels |
| `render-relative-thickness` | float | Line thickness multiplier |
| `render-atom-ids-visible` | true/false | Show atom indices |
| `render-bond-ids-visible` | true/false | Show bond indices |
| `render-margins` | "x, y" | Margins in pixels |
| `render-stereo-style` | none, old, ext, bondmark | Stereochemistry display |
| `render-label-mode` | hetero, terminal-hetero, all, none | Atom label display |

### Setting Options

```go
// String options
render.SetRenderOption("render-output-format", "svg")

// Integer options
render.SetRenderOptionInt("render-image-width", 1024)

// Float options
render.SetRenderOptionFloat("render-relative-thickness", 1.5)

// Boolean options
render.SetRenderOptionBool("render-atom-ids-visible", true)
```

## Supported Formats

- **PNG**: Raster image format (default)
- **SVG**: Vector image format
- **PDF**: Portable Document Format
- **EMF**: Enhanced Metafile (Windows only)

## API Reference

### Initialization

- `InitRenderer()` - Initialize the renderer
- `DisposeRenderer()` - Dispose the renderer
- `ResetRenderer()` - Reset all render settings

### Rendering Functions

- `RenderToFile(objectHandle, filename)` - Render to file
- `Render(objectHandle, outputHandle)` - Render to output buffer
- `RenderGridToFile(arrayHandle, refAtoms, nColumns, filename)` - Render grid to file
- `RenderGrid(arrayHandle, refAtoms, nColumns, outputHandle)` - Render grid to buffer

### Configuration

- `SetRenderOption(option, value)` - Set string option
- `SetRenderOptionInt(option, value)` - Set integer option
- `SetRenderOptionFloat(option, value)` - Set float option
- `SetRenderOptionBool(option, value)` - Set boolean option
- `DefaultRenderOptions()` - Get default options
- `RenderOptions.Apply()` - Apply option set

### Array Management

- `CreateArray()` - Create object array
- `ArrayAdd(arrayHandle, objectHandle)` - Add object to array
- `FreeObject(handle)` - Free Indigo object

### Buffer Operations

- `CreateWriteBuffer()` - Create write buffer
- `GetBufferData(bufferHandle)` - Get buffer contents

## Examples

See the test file `test/render/render_test.go` for more examples.

## Error Handling

All rendering functions return an error. Always check for errors:

```go
err := render.RenderToFile(mol.Handle(), "output.png")
if err != nil {
 log.Fatalf("Render failed: %v", err)
}
```

## Performance Tips

1. Initialize the renderer once and reuse it
2. Use appropriate image sizes (smaller = faster)
3. For batch rendering, use grid rendering when possible
4. Dispose of objects when done to free memory

## License

This package uses the Indigo toolkit, which is licensed under the Apache License 2.0.
