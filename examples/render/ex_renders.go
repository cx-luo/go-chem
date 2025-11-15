// Package render provides example usage of the render package
// This file is for documentation purposes only
package main

import (
	"fmt"
	"github.com/cx-luo/go-indigo/core"
	"log"
	"os"

	"github.com/cx-luo/go-indigo/molecule"
	"github.com/cx-luo/go-indigo/render"
)

// Example 1: Basic molecule rendering
func ExampleBasicRender(indigoInit *core.Indigo) {
	// Load a molecule
	mol, err := indigoInit.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()

	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}
	// Reset renderer settings to defaults
	indigoRender.ResetRenderer()

	// Set PNG format and size
	indigoRender.SetRenderOption("render-output-format", "png")
	indigoRender.SetRenderOptionInt("render-image-width", 800)
	indigoRender.SetRenderOptionInt("render-image-height", 800)

	// Render to file
	if err := indigoRender.RenderToFile(mol.Handle, "benzene.png"); err != nil {
		log.Fatal(err)
	}
}

// Example 2: Using RenderOptions
func ExampleWithOptions(indigoInit *core.Indigo) {
	mol, err := indigoInit.LoadMoleculeFromString("CCO")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()

	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}

	// Reset renderer settings
	indigoRender.ResetRenderer()

	// Configure render options
	opts := indigoRender.Options
	opts.OutputFormat = "svg"
	opts.ImageWidth = 600
	opts.ImageHeight = 400
	opts.BackgroundColor = "0.95, 0.95, 0.95" // Light gray
	opts.BondLength = 50
	opts.ShowAtomIDs = true

	// Apply options
	if err := indigoRender.Apply(); err != nil {
		log.Fatal(err)
	}

	// Render
	if err := indigoRender.RenderToFile(mol.Handle, "ethanol.svg"); err != nil {
		log.Fatal(err)
	}
}

// Example 3: Grid rendering
func ExampleGridRender(indigoInit *core.Indigo) {
	// Load multiple molecules
	molecules := []string{
		"CCO",         // Ethanol
		"c1ccccc1",    // Benzene
		"CC(=O)O",     // Acetic acid
		"CC(C)O",      // Isopropanol
		"c1ccc(O)cc1", // Phenol
		"CCN",         // Ethylamine
	}

	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}

	// Reset renderer settings
	indigoRender.ResetRenderer()

	// Create array
	array, err := indigoInit.CreateArray()
	if err != nil {
		log.Fatal(err)
	}
	defer indigoInit.FreeObject(array)

	// Load and add molecules to array
	for _, smiles := range molecules {
		mol, err := indigoInit.LoadMoleculeFromString(smiles)
		if err != nil {
			log.Printf("Warning: failed to load %s: %v", smiles, err)
			continue
		}
		if err := indigoInit.ArrayAdd(array, mol.Handle); err != nil {
			log.Printf("Warning: failed to add molecule: %v", err)
		}
		// Note: molecules should be kept alive until rendering is done
		mol.Close()
	}

	// Set render options
	opts := indigoRender.Options
	opts.ImageWidth = 1200
	opts.ImageHeight = 800
	if err := indigoRender.Apply(); err != nil {
		log.Fatal(err)
	}

	// Render grid (3 columns)
	if err := indigoRender.RenderGridToFile(array, nil, 3, "molecules_grid.png"); err != nil {
		log.Fatal(err)
	}
}

// Example 4: Rendering a reaction
func ExampleReactionRender(indigoInit *core.Indigo) {
	// Load a reaction
	rxn, err := indigoInit.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		log.Fatal(err)
	}
	defer rxn.Close()
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}
	// Reset renderer settings
	indigoRender.ResetRenderer()

	// Configure for reaction rendering
	indigoRender.SetRenderOption("render-output-format", "png")
	indigoRender.SetRenderOptionInt("render-image-width", 800)
	indigoRender.SetRenderOptionInt("render-image-height", 400)
	indigoRender.SetRenderOption("render-background-color", "1.0, 1.0, 1.0")

	// Render reaction
	if err := indigoRender.RenderToFile(rxn.Handle, "oxidation_reaction.png"); err != nil {
		log.Fatal(err)
	}
}

// Example 5: Render to memory buffer
func ExampleBufferRender(indigoInit *core.Indigo) {
	mol, err := indigoInit.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}
	// Reset renderer settings
	indigoRender.ResetRenderer()

	// Create write buffer
	buffer, err := indigoInit.CreateWriteBuffer()
	if err != nil {
		log.Fatal(err)
	}
	defer indigoInit.FreeObject(buffer)

	// Set format
	if err := indigoRender.SetRenderOption("render-output-format", "png"); err != nil {
		log.Fatal(err)
	}

	// Render to buffer
	if err := indigoRender.Render(mol.Handle, buffer); err != nil {
		log.Fatal(err)
	}

	// Get image data
	imageData, err := indigoInit.GetBufferData(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// Now you can use imageData ([]byte) - write to file, send over network, etc.
	if err := os.WriteFile("benzene_from_buffer.png", imageData, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Buffer render: Generated", len(imageData), "bytes")
}

// Example 6: Batch rendering with different styles
func ExampleBatchRender(indigoInit *core.Indigo) {
	mol, err := indigoInit.LoadMoleculeFromString("CC(C)(C)C")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}
	// Reset renderer settings
	indigoRender.ResetRenderer()

	// Render with different styles
	styles := map[string]string{
		"default":  "hetero",
		"all":      "all",
		"none":     "none",
		"terminal": "terminal-hetero",
	}

	for name, labelMode := range styles {
		if err := indigoRender.SetRenderOption("render-label-mode", labelMode); err != nil {
			log.Printf("Warning: failed to set label mode %s: %v", labelMode, err)
			continue
		}
		if err := indigoRender.RenderToFile(mol.Handle, fmt.Sprintf("neopentane_%s.png", name)); err != nil {
			log.Printf("Warning: failed to render %s: %v", name, err)
		}
	}
}

// Example 7: High-quality rendering
func ExampleHighQualityRender(indigoInit *core.Indigo) {
	mol, err := indigoInit.LoadMoleculeFromString("c1ccc(cc1)C(=O)O") // Benzoic acid
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}
	// Reset renderer settings
	indigoRender.ResetRenderer()

	// High-quality settings
	indigoRender.Options = &render.RenderOptions{
		OutputFormat:      "svg", // Vector format for scalability
		ImageWidth:        2400,  // Large size
		ImageHeight:       2400,
		BackgroundColor:   "1.0, 1.0, 1.0",
		BondLength:        60,
		RelativeThickness: 1.5,
		ShowAtomIDs:       false,
		ShowBondIDs:       false,
		Margins:           "50, 50",
		StereoStyle:       "ext",
		LabelMode:         "hetero",
	}
	if err := indigoRender.Apply(); err != nil {
		log.Fatal(err)
	}

	if err := indigoRender.RenderToFile(mol.Handle, "benzoic_acid_hq.svg"); err != nil {
		log.Fatal(err)
	}
}

// Example 8: Rendering with stereochemistry
func ExampleStereoRender(indigoInit *core.Indigo) {
	// L-Alanine
	mol, err := indigoInit.LoadMoleculeFromString("C[C@H](N)C(=O)O")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}
	// Reset renderer settings
	indigoRender.ResetRenderer()

	// Different stereo styles
	stereoStyles := []string{"old", "ext", "none"}

	for _, style := range stereoStyles {
		if err := indigoRender.SetRenderOption("render-stereo-style", style); err != nil {
			log.Printf("Warning: failed to set stereo style %s: %v", style, err)
			continue
		}
		if err := indigoRender.RenderToFile(mol.Handle, fmt.Sprintf("alanine_%s.png", style)); err != nil {
			log.Printf("Warning: failed to render %s: %v", style, err)
		}
	}
}

// Example 9: Grid with reference atoms (alignment)
func ExampleAlignedGrid(indigoInit *core.Indigo) {
	// Series of alcohols - align on oxygen
	alcohols := []string{"CO", "CCO", "CCCO", "CC(C)O"}
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}
	// Reset renderer settings
	indigoRender.ResetRenderer()

	array, err := indigoInit.CreateArray()
	if err != nil {
		log.Fatal(err)
	}
	defer indigoInit.FreeObject(array)

	var mols []*molecule.Molecule
	for _, smiles := range alcohols {
		mol, err := indigoInit.LoadMoleculeFromString(smiles)
		if err != nil {
			log.Printf("Warning: failed to load %s: %v", smiles, err)
			continue
		}
		mols = append(mols, mol)
		if err := indigoInit.ArrayAdd(array, mol.Handle); err != nil {
			log.Printf("Warning: failed to add molecule: %v", err)
		}
	}
	defer func() {
		for _, mol := range mols {
			mol.Close()
		}
	}()

	// Reference atoms (oxygen atom index in each molecule)
	// This aligns all molecules on their oxygen atoms
	refAtoms := []int{0, 2, 3, 2} // Oxygen positions

	if err := indigoRender.RenderGridToFile(array, refAtoms, 2, "alcohols_aligned.png"); err != nil {
		log.Fatal(err)
	}
}

// Example 10: Custom colors and styling
func ExampleCustomColors(indigoInit *core.Indigo) {
	mol, err := indigoInit.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}
	// Reset renderer settings
	indigoRender.ResetRenderer()

	// Custom background color (light blue)
	indigoRender.SetRenderOption("render-background-color", "0.9, 0.95, 1.0")

	// Thicker lines
	indigoRender.SetRenderOptionFloat("render-relative-thickness", 2.0)

	// Larger bonds
	indigoRender.SetRenderOptionInt("render-bond-length", 60)

	if err := indigoRender.RenderToFile(mol.Handle, "benzene_custom.png"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	indigoInit, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Render Package Example ===")
	indigoRender, err := indigoInit.InitRenderer()
	if err != nil {
		fmt.Printf("failed to initialize renderer: %v", err)
	}

	defer indigoRender.DisposeRenderer()

	// Run all examples
	ExampleBasicRender(indigoInit)
	fmt.Println("✓ Basic render completed")

	ExampleWithOptions(indigoInit)
	fmt.Println("✓ With options completed")

	ExampleGridRender(indigoInit)
	fmt.Println("✓ Grid render completed")

	ExampleReactionRender(indigoInit)
	fmt.Println("✓ Reaction render completed")

	ExampleBufferRender(indigoInit)
	fmt.Println("✓ Buffer render completed")

	ExampleBatchRender(indigoInit)
	fmt.Println("✓ Batch render completed")

	ExampleHighQualityRender(indigoInit)
	fmt.Println("✓ High-quality render completed")

	ExampleStereoRender(indigoInit)
	fmt.Println("✓ Stereo render completed")

	ExampleAlignedGrid(indigoInit)
	fmt.Println("✓ Aligned grid completed")

	ExampleCustomColors(indigoInit)
	fmt.Println("✓ Custom colors completed")

	fmt.Println("\n=== All examples completed successfully! ===")
}
