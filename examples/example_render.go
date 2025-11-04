// Package render provides example usage of the render package
// This file is for documentation purposes only
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/reaction"
	"github.com/cx-luo/go-chem/render"
)

// Example 1: Basic molecule rendering
func ExampleBasicRender() {
	// Load a molecule
	mol, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()

	// Reset renderer settings to defaults
	render.ResetRenderer()

	// Set PNG format and size
	render.SetRenderOption("render-output-format", "png")
	render.SetRenderOptionInt("render-image-width", 800)
	render.SetRenderOptionInt("render-image-height", 800)

	// Render to file
	if err := render.RenderToFile(mol.Handle(), "benzene.png"); err != nil {
		log.Fatal(err)
	}
}

// Example 2: Using RenderOptions
func ExampleWithOptions() {
	mol, err := molecule.LoadMoleculeFromString("CCO")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()

	// Reset renderer settings
	render.ResetRenderer()

	// Configure render options
	opts := render.DefaultRenderOptions()
	opts.OutputFormat = "svg"
	opts.ImageWidth = 600
	opts.ImageHeight = 400
	opts.BackgroundColor = "0.95, 0.95, 0.95" // Light gray
	opts.BondLength = 50
	opts.ShowAtomIDs = true

	// Apply options
	if err := opts.Apply(); err != nil {
		log.Fatal(err)
	}

	// Render
	if err := render.RenderToFile(mol.Handle(), "ethanol.svg"); err != nil {
		log.Fatal(err)
	}
}

// Example 3: Grid rendering
func ExampleGridRender() {
	// Load multiple molecules
	molecules := []string{
		"CCO",         // Ethanol
		"c1ccccc1",    // Benzene
		"CC(=O)O",     // Acetic acid
		"CC(C)O",      // Isopropanol
		"c1ccc(O)cc1", // Phenol
		"CCN",         // Ethylamine
	}

	// Reset renderer settings
	render.ResetRenderer()

	// Create array
	array, err := render.CreateArray()
	if err != nil {
		log.Fatal(err)
	}
	defer render.FreeObject(array)

	// Load and add molecules to array
	for _, smiles := range molecules {
		mol, err := molecule.LoadMoleculeFromString(smiles)
		if err != nil {
			log.Printf("Warning: failed to load %s: %v", smiles, err)
			continue
		}
		if err := render.ArrayAdd(array, mol.Handle()); err != nil {
			log.Printf("Warning: failed to add molecule: %v", err)
		}
		// Note: molecules should be kept alive until rendering is done
		defer mol.Close()
	}

	// Set render options
	opts := render.DefaultRenderOptions()
	opts.ImageWidth = 1200
	opts.ImageHeight = 800
	if err := opts.Apply(); err != nil {
		log.Fatal(err)
	}

	// Render grid (3 columns)
	if err := render.RenderGridToFile(array, nil, 3, "molecules_grid.png"); err != nil {
		log.Fatal(err)
	}
}

// Example 4: Rendering a reaction
func ExampleReactionRender() {
	// Load a reaction
	rxn, err := reaction.LoadReactionFromString("CCO>>CC=O")
	if err != nil {
		log.Fatal(err)
	}
	defer rxn.Close()

	// Reset renderer settings
	render.ResetRenderer()

	// Configure for reaction rendering
	render.SetRenderOption("render-output-format", "png")
	render.SetRenderOptionInt("render-image-width", 800)
	render.SetRenderOptionInt("render-image-height", 400)
	render.SetRenderOption("render-background-color", "1.0, 1.0, 1.0")

	// Render reaction
	if err := render.RenderToFile(rxn.Handle(), "oxidation_reaction.png"); err != nil {
		log.Fatal(err)
	}
}

// Example 5: Render to memory buffer
func ExampleBufferRender() {
	mol, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()

	// Reset renderer settings
	render.ResetRenderer()

	// Create write buffer
	buffer, err := render.CreateWriteBuffer()
	if err != nil {
		log.Fatal(err)
	}
	defer render.FreeObject(buffer)

	// Set format
	if err := render.SetRenderOption("render-output-format", "png"); err != nil {
		log.Fatal(err)
	}

	// Render to buffer
	if err := render.Render(mol.Handle(), buffer); err != nil {
		log.Fatal(err)
	}

	// Get image data
	imageData, err := render.GetBufferData(buffer)
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
func ExampleBatchRender() {
	mol, err := molecule.LoadMoleculeFromString("CC(C)(C)C")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()

	// Reset renderer settings
	render.ResetRenderer()

	// Render with different styles
	styles := map[string]string{
		"default":  "hetero",
		"all":      "all",
		"none":     "none",
		"terminal": "terminal-hetero",
	}

	for name, labelMode := range styles {
		if err := render.SetRenderOption("render-label-mode", labelMode); err != nil {
			log.Printf("Warning: failed to set label mode %s: %v", labelMode, err)
			continue
		}
		if err := render.RenderToFile(mol.Handle(), fmt.Sprintf("neopentane_%s.png", name)); err != nil {
			log.Printf("Warning: failed to render %s: %v", name, err)
		}
	}
}

// Example 7: High-quality rendering
func ExampleHighQualityRender() {
	mol, err := molecule.LoadMoleculeFromString("c1ccc(cc1)C(=O)O") // Benzoic acid
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()

	// Reset renderer settings
	render.ResetRenderer()

	// High-quality settings
	opts := &render.RenderOptions{
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
	if err := opts.Apply(); err != nil {
		log.Fatal(err)
	}

	if err := render.RenderToFile(mol.Handle(), "benzoic_acid_hq.svg"); err != nil {
		log.Fatal(err)
	}
}

// Example 8: Rendering with stereochemistry
func ExampleStereoRender() {
	// L-Alanine
	mol, err := molecule.LoadMoleculeFromString("C[C@H](N)C(=O)O")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()

	// Reset renderer settings
	render.ResetRenderer()

	// Different stereo styles
	stereoStyles := []string{"old", "ext", "none"}

	for _, style := range stereoStyles {
		if err := render.SetRenderOption("render-stereo-style", style); err != nil {
			log.Printf("Warning: failed to set stereo style %s: %v", style, err)
			continue
		}
		if err := render.RenderToFile(mol.Handle(), fmt.Sprintf("alanine_%s.png", style)); err != nil {
			log.Printf("Warning: failed to render %s: %v", style, err)
		}
	}
}

// Example 9: Grid with reference atoms (alignment)
func ExampleAlignedGrid() {
	// Series of alcohols - align on oxygen
	alcohols := []string{"CO", "CCO", "CCCO", "CC(C)O"}

	// Reset renderer settings
	render.ResetRenderer()

	array, err := render.CreateArray()
	if err != nil {
		log.Fatal(err)
	}
	defer render.FreeObject(array)

	var mols []*molecule.Molecule
	for _, smiles := range alcohols {
		mol, err := molecule.LoadMoleculeFromString(smiles)
		if err != nil {
			log.Printf("Warning: failed to load %s: %v", smiles, err)
			continue
		}
		mols = append(mols, mol)
		if err := render.ArrayAdd(array, mol.Handle()); err != nil {
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

	if err := render.RenderGridToFile(array, refAtoms, 2, "alcohols_aligned.png"); err != nil {
		log.Fatal(err)
	}
}

// Example 10: Custom colors and styling
func ExampleCustomColors() {
	mol, err := molecule.LoadMoleculeFromString("c1ccccc1")
	if err != nil {
		log.Fatal(err)
	}
	defer mol.Close()

	// Reset renderer settings
	render.ResetRenderer()

	// Custom background color (light blue)
	render.SetRenderOption("render-background-color", "0.9, 0.95, 1.0")

	// Thicker lines
	render.SetRenderOptionFloat("render-relative-thickness", 2.0)

	// Larger bonds
	render.SetRenderOptionInt("render-bond-length", 60)

	if err := render.RenderToFile(mol.Handle(), "benzene_custom.png"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("=== Render Package Example ===")

	// Initialize renderer once at the start
	if err := render.InitRenderer(); err != nil {
		log.Fatalf("Failed to initialize renderer: %v", err)
	}
	defer render.DisposeRenderer()

	// Run all examples
	ExampleBasicRender()
	fmt.Println("✓ Basic render completed")

	ExampleWithOptions()
	fmt.Println("✓ With options completed")

	ExampleGridRender()
	fmt.Println("✓ Grid render completed")

	ExampleReactionRender()
	fmt.Println("✓ Reaction render completed")

	ExampleBufferRender()
	fmt.Println("✓ Buffer render completed")

	ExampleBatchRender()
	fmt.Println("✓ Batch render completed")

	ExampleHighQualityRender()
	fmt.Println("✓ High-quality render completed")

	ExampleStereoRender()
	fmt.Println("✓ Stereo render completed")

	ExampleAlignedGrid()
	fmt.Println("✓ Aligned grid completed")

	ExampleCustomColors()
	fmt.Println("✓ Custom colors completed")

	fmt.Println("\n=== All examples completed successfully! ===")
}
