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

	// Initialize renderer
	if err := render.InitRenderer(); err != nil {
		log.Fatal(err)
	}
	defer render.DisposeRenderer()

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
	mol, _ := molecule.LoadMoleculeFromString("CCO")
	defer mol.Close()

	render.InitRenderer()
	defer render.DisposeRenderer()

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
	render.RenderToFile(mol.Handle(), "ethanol.svg")
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

	render.InitRenderer()
	defer render.DisposeRenderer()

	// Create array
	array, _ := render.CreateArray()
	defer render.FreeObject(array)

	// Load and add molecules to array
	for _, smiles := range molecules {
		mol, _ := molecule.LoadMoleculeFromString(smiles)
		render.ArrayAdd(array, mol.Handle())
		// Note: molecules should be kept alive until rendering is done
		defer mol.Close()
	}

	// Set render options
	opts := render.DefaultRenderOptions()
	opts.ImageWidth = 1200
	opts.ImageHeight = 800
	opts.Apply()

	// Render grid (3 columns)
	render.RenderGridToFile(array, nil, 3, "molecules_grid.png")
}

// Example 4: Rendering a reaction
func ExampleReactionRender() {
	// Load a reaction
	rxn, _ := reaction.LoadReactionFromString("CCO>>CC=O")
	defer rxn.Close()

	render.InitRenderer()
	defer render.DisposeRenderer()

	// Configure for reaction rendering
	render.SetRenderOption("render-output-format", "png")
	render.SetRenderOptionInt("render-image-width", 800)
	render.SetRenderOptionInt("render-image-height", 400)
	render.SetRenderOption("render-background-color", "1.0, 1.0, 1.0")

	// Render reaction
	render.RenderToFile(rxn.Handle(), "oxidation_reaction.png")
}

// Example 5: Render to memory buffer
func ExampleBufferRender() {
	mol, _ := molecule.LoadMoleculeFromString("c1ccccc1")
	defer mol.Close()

	render.InitRenderer()
	defer render.DisposeRenderer()

	// Create write buffer
	buffer, _ := render.CreateWriteBuffer()
	defer render.FreeObject(buffer)

	// Set format
	render.SetRenderOption("render-output-format", "png")

	// Render to buffer
	if err := render.Render(mol.Handle(), buffer); err != nil {
		log.Fatal(err)
	}

	// Get image data
	imageData, _ := render.GetBufferData(buffer)

	// Now you can use imageData ([]byte) - write to file, send over network, etc.
	os.WriteFile("benzene_from_buffer.png", imageData, 0644)
}

// Example 6: Batch rendering with different styles
func ExampleBatchRender() {
	mol, _ := molecule.LoadMoleculeFromString("CC(C)(C)C")
	defer mol.Close()

	render.InitRenderer()
	defer render.DisposeRenderer()

	// Render with different styles
	styles := map[string]string{
		"default":  "hetero",
		"all":      "all",
		"none":     "none",
		"terminal": "terminal-hetero",
	}

	for name, labelMode := range styles {
		render.SetRenderOption("render-label-mode", labelMode)
		render.RenderToFile(mol.Handle(), fmt.Sprintf("neopentane_%s.png", name))
	}
}

// Example 7: High-quality rendering
func ExampleHighQualityRender() {
	mol, _ := molecule.LoadMoleculeFromString("c1ccc(cc1)C(=O)O") // Benzoic acid
	defer mol.Close()

	render.InitRenderer()
	defer render.DisposeRenderer()

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
	opts.Apply()

	render.RenderToFile(mol.Handle(), "benzoic_acid_hq.svg")
}

// Example 8: Rendering with stereochemistry
func ExampleStereoRender() {
	// L-Alanine
	mol, _ := molecule.LoadMoleculeFromString("C[C@H](N)C(=O)O")
	defer mol.Close()

	render.InitRenderer()
	defer render.DisposeRenderer()

	// Different stereo styles
	stereoStyles := []string{"old", "ext", "bondmark"}

	for _, style := range stereoStyles {
		render.SetRenderOption("render-stereo-style", style)
		render.RenderToFile(mol.Handle(), fmt.Sprintf("alanine_%s.png", style))
	}
}

// Example 9: Grid with reference atoms (alignment)
func ExampleAlignedGrid() {
	// Series of alcohols - align on oxygen
	alcohols := []string{"CO", "CCO", "CCCO", "CC(C)O"}

	render.InitRenderer()
	defer render.DisposeRenderer()

	array, _ := render.CreateArray()
	defer render.FreeObject(array)

	var mols []*molecule.Molecule
	for _, smiles := range alcohols {
		mol, _ := molecule.LoadMoleculeFromString(smiles)
		mols = append(mols, mol)
		render.ArrayAdd(array, mol.Handle())
	}
	defer func() {
		for _, mol := range mols {
			mol.Close()
		}
	}()

	// Reference atoms (oxygen atom index in each molecule)
	// This aligns all molecules on their oxygen atoms
	refAtoms := []int{0, 2, 3, 2} // Oxygen positions

	render.RenderGridToFile(array, refAtoms, 2, "alcohols_aligned.png")
}

// Example 10: Custom colors and styling
func ExampleCustomColors() {
	mol, _ := molecule.LoadMoleculeFromString("c1ccccc1")
	defer mol.Close()

	render.InitRenderer()
	defer render.DisposeRenderer()

	// Custom background color (light blue)
	render.SetRenderOption("render-background-color", "0.9, 0.95, 1.0")

	// Thicker lines
	render.SetRenderOptionFloat("render-relative-thickness", 2.0)

	// Larger bonds
	render.SetRenderOptionInt("render-bond-length", 60)

	render.RenderToFile(mol.Handle(), "benzene_custom.png")
}
