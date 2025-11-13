package render_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cx-luo/go-chem/core"

	"github.com/cx-luo/go-chem/render"
)

var indigoInit *core.Indigo

func init() {
	handle, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}
	indigoInit = handle
}

// TestInitRenderer tests initializing the renderer
func TestInitRenderer(t *testing.T) {
	_, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}

	// Initialize again should not error
	_, err = indigoInit.InitRenderer()
	if err != nil {
		t.Errorf("second initialization should not error: %v", err)
	}
}

// TestDisposeRenderer tests disposing the renderer
func TestDisposeRenderer(t *testing.T) {
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}

	err = indigoRender.DisposeRenderer()
	if err != nil {
		t.Errorf("failed to dispose renderer: %v", err)
	}

	// Dispose again should not error
	err = indigoRender.DisposeRenderer()
	if err != nil {
		t.Errorf("second dispose should not error: %v", err)
	}
}

// TestRenderMoleculeToFile tests rendering a molecule to a file
func TestRenderMoleculeToFile(t *testing.T) {
	// Load a molecule
	mol, err := indigoInit.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer mol.Close()

	// Create temp file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "benzene.png")
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}

	// Set render options
	opts := indigoRender.Options
	opts.ImageWidth = 300
	opts.ImageHeight = 300
	if err := indigoRender.Apply(); err != nil {
		t.Fatalf("failed to apply render options: %v", err)
	}

	// Render to file
	err = indigoRender.RenderToFile(mol.Handle, outputFile)
	if err != nil {
		t.Fatalf("failed to render molecule to file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("output file was not created: %s", outputFile)
	}

	// Verify file has content
	info, _ := os.Stat(outputFile)
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

// TestRenderMoleculeSVG tests rendering a molecule to SVG format
func TestRenderMoleculeSVG(t *testing.T) {
	// Load a molecule
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer mol.Close()

	// Create temp file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "ethanol.svg")
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}
	// Set SVG format
	if err := indigoRender.SetRenderOption("render-output-format", "svg"); err != nil {
		t.Fatalf("failed to set SVG format: %v", err)
	}

	// Render to file
	err = indigoRender.RenderToFile(mol.Handle, outputFile)
	if err != nil {
		t.Fatalf("failed to render molecule to SVG: %v", err)
	}

	// Verify file exists and has content
	info, err := os.Stat(outputFile)
	if err != nil {
		t.Errorf("output file was not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output SVG file is empty")
	}
}

// TestRenderReactionToFile tests rendering a reaction to a file
func TestRenderReactionToFile(t *testing.T) {
	// Load a reaction
	rxn, err := indigoInit.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		t.Fatalf("failed to load reaction: %v", err)
	}
	defer rxn.Close()

	// Create temp file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "reaction.png")
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}
	// Set render options
	opts := indigoRender.Options
	opts.ImageWidth = 600
	opts.ImageHeight = 300
	if err := indigoRender.Apply(); err != nil {
		t.Fatalf("failed to apply render options: %v", err)
	}

	// Render to file
	err = indigoRender.RenderToFile(rxn.Handle, outputFile)
	if err != nil {
		t.Fatalf("failed to render reaction to file: %v", err)
	}

	// Verify file exists and has content
	info, err := os.Stat(outputFile)
	if err != nil {
		t.Errorf("output file was not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

// TestRenderOptions tests various render options
func TestRenderOptions(t *testing.T) {
	tests := []struct {
		name   string
		option string
		value  string
	}{
		{"OutputFormat", "render-output-format", "png"},
		{"BackgroundColor", "render-background-color", "1.0, 1.0, 1.0"},
		{"StereoStyle", "render-stereo-style", "ext"},
		{"LabelMode", "render-label-mode", "hetero"},
	}

	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := indigoRender.SetRenderOption(tt.option, tt.value)
			if err != nil {
				t.Errorf("failed to set option %s: %v", tt.option, err)
			}
		})
	}
}

// TestRenderOptionsInt tests integer render options
func TestRenderOptionsInt(t *testing.T) {
	tests := []struct {
		name   string
		option string
		value  int
	}{
		{"ImageWidth", "render-image-width", 800},
		{"ImageHeight", "render-image-height", 600},
		{"BondLength", "render-bond-length", 50},
	}
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := indigoRender.SetRenderOptionInt(tt.option, tt.value)
			if err != nil {
				t.Errorf("failed to set option %s: %v", tt.option, err)
			}
		})
	}
}

// TestRenderOptionsFloat tests float render options
func TestRenderOptionsFloat(t *testing.T) {
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}
	err = indigoRender.SetRenderOptionFloat("render-relative-thickness", 1.5)
	if err != nil {
		t.Errorf("failed to set relative thickness: %v", err)
	}
}

// TestRenderOptionsBool tests boolean render options
func TestRenderOptionsBool(t *testing.T) {
	tests := []struct {
		name   string
		option string
		value  bool
	}{
		{"ShowAtomIDs", "render-atom-ids-visible", true},
		{"ShowBondIDs", "render-bond-ids-visible", false},
	}
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := indigoRender.SetRenderOptionBool(tt.option, tt.value)
			if err != nil {
				t.Errorf("failed to set option %s: %v", tt.option, err)
			}
		})
	}
}

// TestDefaultRenderOptions tests default render options
func TestDefaultRenderOptions(t *testing.T) {
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}

	opts := indigoRender.Options

	if opts.OutputFormat != "png" {
		t.Errorf("expected default format 'png', got '%s'", opts.OutputFormat)
	}

	if opts.ImageWidth != 1600 {
		t.Errorf("expected default width 1600, got %d", opts.ImageWidth)
	}

	// Apply options
	err = indigoRender.Apply()
	if err != nil {
		t.Errorf("failed to apply default options: %v", err)
	}
}

// TestRenderOptionsApply tests applying custom render options
func TestRenderOptionsApply(t *testing.T) {
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}
	indigoRender.Options = &render.RenderOptions{
		OutputFormat:      "svg",
		ImageWidth:        500,
		ImageHeight:       500,
		BackgroundColor:   "0.9, 0.9, 0.9",
		BondLength:        30,
		RelativeThickness: 1.2,
		ShowAtomIDs:       true,
		ShowBondIDs:       false,
		Margins:           "20, 20",
		StereoStyle:       "ext",
		LabelMode:         "all",
	}

	err = indigoRender.Apply()
	if err != nil {
		t.Errorf("failed to apply custom options: %v", err)
	}
}

// TestResetRenderer tests resetting the renderer
func TestResetRenderer(t *testing.T) {
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		t.Fatalf("failed to initialize renderer: %v", err)
	}
	// Set some options
	indigoRender.SetRenderOptionInt("render-image-width", 800)
	indigoRender.SetRenderOption("render-output-format", "svg")

	// Reset
	err = indigoRender.ResetRenderer()
	if err != nil {
		t.Errorf("failed to reset renderer: %v", err)
	}
}

// TestRenderGrid tests rendering multiple molecules in a grid
func TestRenderGrid(t *testing.T) {
	// Create molecules
	mol1, _ := indigoInit.LoadMoleculeFromString("CCO")
	defer mol1.Close()

	mol2, _ := indigoInit.LoadMoleculeFromString("c1ccccc1")
	defer mol2.Close()

	mol3, _ := indigoInit.LoadMoleculeFromString("CC(=O)O")
	defer mol3.Close()

	// Create array
	arrayHandle, err := indigoInit.CreateArray()
	if err != nil {
		t.Fatalf("failed to create array: %v", err)
	}
	defer indigoInit.FreeObject(arrayHandle)

	// Add molecules to array
	if err := indigoInit.ArrayAdd(arrayHandle, mol1.Handle); err != nil {
		t.Fatalf("failed to add mol1 to array: %v", err)
	}
	if err := indigoInit.ArrayAdd(arrayHandle, mol2.Handle); err != nil {
		t.Fatalf("failed to add mol2 to array: %v", err)
	}
	if err := indigoInit.ArrayAdd(arrayHandle, mol3.Handle); err != nil {
		t.Fatalf("failed to add mol3 to array: %v", err)
	}

	// Create temp file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "grid.png")

	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		panic(err)
	}

	// Render grid
	err = indigoRender.RenderGridToFile(arrayHandle, nil, 2, outputFile)
	if err != nil {
		t.Fatalf("failed to render grid: %v", err)
	}

	// Verify file exists and has content
	info, err := os.Stat(outputFile)
	if err != nil {
		t.Errorf("output file was not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output grid file is empty")
	}
}

// TestRenderGridWithRefAtoms tests rendering grid with reference atoms
func TestRenderGridWithRefAtoms(t *testing.T) {
	// Create molecules
	mol1, _ := indigoInit.LoadMoleculeFromString("CCO")
	defer mol1.Close()

	mol2, _ := indigoInit.LoadMoleculeFromString("CCCO")
	defer mol2.Close()

	// Create array
	arrayHandle, err := indigoInit.CreateArray()
	if err != nil {
		t.Fatalf("failed to create array: %v", err)
	}
	defer indigoInit.FreeObject(arrayHandle)

	// Add molecules
	indigoInit.ArrayAdd(arrayHandle, mol1.Handle)
	indigoInit.ArrayAdd(arrayHandle, mol2.Handle)

	// Create temp file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "grid_ref.png")

	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		panic(err)
	}
	// Render grid with reference atoms
	refAtoms := []int{0, 0} // First atom of each molecule
	err = indigoRender.RenderGridToFile(arrayHandle, refAtoms, 2, outputFile)
	if err != nil {
		t.Fatalf("failed to render grid with ref atoms: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}

// TestRenderToBuffer tests rendering to a memory buffer
func TestRenderToBuffer(t *testing.T) {
	// Load a molecule
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer mol.Close()

	// Create write buffer
	bufferHandle, err := indigoInit.CreateWriteBuffer()
	if err != nil {
		t.Fatalf("failed to create write buffer: %v", err)
	}
	defer indigoInit.FreeObject(bufferHandle)

	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		panic(err)
	}
	// Set format to PNG
	indigoRender.SetRenderOption("render-output-format", "png")

	// Render to buffer
	err = indigoRender.Render(mol.Handle, bufferHandle)
	if err != nil {
		t.Fatalf("failed to render to buffer: %v", err)
	}

	// Get buffer data
	data, err := indigoInit.GetBufferData(bufferHandle)
	if err != nil {
		t.Fatalf("failed to get buffer data: %v", err)
	}

	if len(data) == 0 {
		t.Error("buffer data is empty")
	}

	// PNG files start with specific magic bytes
	if len(data) < 8 {
		t.Error("buffer data too short to be a valid PNG")
	}
}

// TestRenderInvalidHandle tests rendering with invalid handle
func TestRenderInvalidHandle(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "invalid.png")

	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		panic(err)
	}
	defer indigoRender.DisposeRenderer()
	err = indigoRender.RenderToFile(-1, outputFile)
	if err != nil {
		t.Error("expected error when rendering with invalid handle")
	}
}

// TestRenderMultipleFormats tests rendering to different formats
func TestRenderMultipleFormats(t *testing.T) {
	// Load a molecule
	mol, err := indigoInit.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		t.Fatalf("failed to load molecule: %v", err)
	}
	defer mol.Close()

	tmpDir := t.TempDir()

	formats := []string{"png", "svg"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			outputFile := filepath.Join(tmpDir, "benzene."+format)

			indigoRender, err := indigoInit.InitRenderer()
			if err != nil {
				panic(err)
			}
			// Set format
			if err := indigoRender.SetRenderOption("render-output-format", format); err != nil {
				t.Fatalf("failed to set format %s: %v", format, err)
			}

			// Render
			if err := indigoRender.RenderToFile(mol.Handle, outputFile); err != nil {
				t.Fatalf("failed to render to %s: %v", format, err)
			}

			// Verify
			info, err := os.Stat(outputFile)
			if err != nil {
				t.Errorf("output file not created for format %s: %v", format, err)
			}
			if info.Size() == 0 {
				t.Errorf("output file is empty for format %s", format)
			}
		})
	}
}
