// Package render provides molecule and reaction rendering using Indigo Renderer via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/04
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : render.go
// @Software: GoLand
package render

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd

// Windows platforms
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../3rd/windows-x86_64 -lindigo -lindigo-renderer
#cgo windows,386 LDFLAGS: -L${SRCDIR}/../3rd/windows-i386 -lindigo -lindigo-renderer

// Linux: use $ORIGIN for runtime library search
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -lindigo-renderer -Wl,-rpath,${SRCDIR}/../3rd/linux-x86_64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../3rd/linux-aarch64 -lindigo -lindigo-renderer -Wl,-rpath,${SRCDIR}/../3rd/linux-aarch64

// macOS: use @loader_path (not @executable_path) for shared libraries
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-x86_64 -lindigo -lindigo-renderer -Wl,-rpath,${SRCDIR}/../3rd/darwin-x86_64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-aarch64 -lindigo -lindigo-renderer -Wl,-rpath,${SRCDIR}/../3rd/darwin-aarch64

#include <stdlib.h>
#include "indigo.h"
#include "indigo-renderer.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Renderer struct {
	Sid                 uint64
	Options             *RenderOptions
	RendererInitialized bool
}

// DisposeRenderer disposes the Indigo renderer
// This should be called when done using rendering functions
func (r *Renderer) DisposeRenderer() error {
	if !r.RendererInitialized {
		return nil // Not initialized
	}

	ret := int(C.indigoRendererDispose(C.ulonglong(r.Sid)))
	if ret < 0 {
		return fmt.Errorf("failed to dispose renderer: %s", getLastError())
	}

	r.RendererInitialized = false
	return nil
}

// ResetRenderer resets all rendering settings to defaults
func (r *Renderer) ResetRenderer() error {
	ret := int(C.indigoRenderReset())
	r.RendererInitialized = false
	if ret < 0 {
		return fmt.Errorf("failed to reset renderer: %s", getLastError())
	}
	return nil
}

// RenderToFile renders an object (molecule or reaction) to a file
// objectHandle: the Indigo handle of the object to render
// filename: the output file path (e.g., "molecule.png", "reaction.svg")
func (r *Renderer) RenderToFile(objectHandle int, filename string) error {
	if objectHandle < 0 {
		return fmt.Errorf("invalid object handle")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoRenderToFile(C.int(objectHandle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to render to file %s: %s", filename, getLastError())
	}

	return nil
}

// Render renders an object to an output buffer
// objectHandle: the Indigo handle of the object to render
// outputHandle: the Indigo handle of the output buffer (from indigoWriteBuffer or indigoWriteFile)
func (r *Renderer) Render(objectHandle int, outputHandle int) error {
	if objectHandle < 0 {
		return fmt.Errorf("invalid object handle")
	}
	if outputHandle < 0 {
		return fmt.Errorf("invalid output handle")
	}

	ret := int(C.indigoRender(C.int(objectHandle), C.int(outputHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to render: %s", getLastError())
	}

	return nil
}

// RenderGridToFile renders a grid of molecules to a file
// arrayHandle: the Indigo handle of an array of molecules (created with indigoCreateArray)
// refAtoms: optional array of reference atom indices (nil for automatic)
// nColumns: number of columns in the grid
// filename: the output file path
func (r *Renderer) RenderGridToFile(arrayHandle int, refAtoms []int, nColumns int, filename string) error {
	if arrayHandle < 0 {
		return fmt.Errorf("invalid array handle")
	}
	if nColumns <= 0 {
		return fmt.Errorf("invalid number of columns: %d", nColumns)
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	var refAtomsPtr *C.int
	if len(refAtoms) > 0 {
		// Convert Go slice to C array without C.malloc (to avoid cgo malloc issues)
		cRefAtoms := make([]C.int, len(refAtoms))
		for i, v := range refAtoms {
			cRefAtoms[i] = C.int(v)
		}
		refAtomsPtr = &cRefAtoms[0]
	}

	ret := int(C.indigoRenderGridToFile(C.int(arrayHandle), refAtomsPtr, C.int(nColumns), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to render grid to file %s: %s", filename, getLastError())
	}

	return nil
}

// RenderGrid renders a grid of molecules to an output buffer
// arrayHandle: the Indigo handle of an array of molecules
// refAtoms: optional array of reference atom indices (nil for automatic)
// nColumns: number of columns in the grid
// outputHandle: the Indigo handle of the output buffer
func (r *Renderer) RenderGrid(arrayHandle int, refAtoms []int, nColumns int, outputHandle int) error {
	if arrayHandle < 0 {
		return fmt.Errorf("invalid array handle")
	}
	if nColumns <= 0 {
		return fmt.Errorf("invalid number of columns: %d", nColumns)
	}
	if outputHandle < 0 {
		return fmt.Errorf("invalid output handle")
	}

	var refAtomsPtr *C.int
	if len(refAtoms) > 0 {
		// Convert Go slice to C array
		cRefAtoms := make([]C.int, len(refAtoms))
		for i, v := range refAtoms {
			cRefAtoms[i] = C.int(v)
		}
		refAtomsPtr = &cRefAtoms[0]
	}

	ret := int(C.indigoRenderGrid(C.int(arrayHandle), refAtomsPtr, C.int(nColumns), C.int(outputHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to render grid: %s", getLastError())
	}

	return nil
}

// SetRenderOption sets a rendering option
// Common options:
//   - "render-output-format": "png", "svg", "pdf", "emf" (Windows)
//   - "render-image-width": width in pixels (default: 1600)
//   - "render-image-height": height in pixels (default: 1600)
//   - "render-background-color": "1.0, 1.0, 1.0" (RGB, white background)
//   - "render-bond-length": bond length in pixels (default: 40)
//   - "render-relative-thickness": line thickness (default: 1.0)
//   - "render-atom-ids-visible": "true" or "false"
//   - "render-bond-ids-visible": "true" or "false"
//   - "render-margins": "10, 10" (x, y margins)
//   - "render-stereo-style": "none", "old", "ext", "bondmark"
//   - "render-label-mode": "hetero", "terminal-hetero", "all", "none"
func (r *Renderer) SetRenderOption(option string, value string) error {
	cOption := C.CString(option)
	defer C.free(unsafe.Pointer(cOption))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	ret := int(C.indigoSetOption(cOption, cValue))
	if ret < 0 {
		return fmt.Errorf("failed to set render option %s: %s", option, getLastError())
	}

	return nil
}

// SetRenderOptionInt sets a rendering option with an integer value
func (r *Renderer) SetRenderOptionInt(option string, value int) error {
	cOption := C.CString(option)
	defer C.free(unsafe.Pointer(cOption))

	cValue := C.CString(fmt.Sprintf("%d", value))
	defer C.free(unsafe.Pointer(cValue))

	ret := int(C.indigoSetOption(cOption, cValue))
	if ret < 0 {
		return fmt.Errorf("failed to set render option %s: %s", option, getLastError())
	}

	return nil
}

// SetRenderOptionFloat sets a rendering option with a float value
func (r *Renderer) SetRenderOptionFloat(option string, value float64) error {
	cOption := C.CString(option)
	defer C.free(unsafe.Pointer(cOption))

	cValue := C.CString(fmt.Sprintf("%f", value))
	defer C.free(unsafe.Pointer(cValue))

	ret := int(C.indigoSetOption(cOption, cValue))
	if ret < 0 {
		return fmt.Errorf("failed to set render option %s: %s", option, getLastError())
	}

	return nil
}

// SetRenderOptionBool sets a rendering option with a boolean value
func (r *Renderer) SetRenderOptionBool(option string, value bool) error {
	strValue := "false"
	if value {
		strValue = "true"
	}
	return r.SetRenderOption(option, strValue)
}

// RenderOptions provides a convenient way to configure rendering settings
type RenderOptions struct {
	OutputFormat      string  // "png", "svg", "pdf", "emf"
	ImageWidth        int     // Width in pixels
	ImageHeight       int     // Height in pixels
	BackgroundColor   string  // RGB color (e.g., "1.0, 1.0, 1.0")
	BondLength        int     // Bond length in pixels
	RelativeThickness float64 // Line thickness
	ShowAtomIDs       bool    // Show atom IDs
	ShowBondIDs       bool    // Show bond IDs
	Margins           string  // Margins (e.g., "10, 10")
	StereoStyle       string  // "none", "old", "ext", "bondmark"
	LabelMode         string  // "hetero", "terminal-hetero", "all", "none"
}

// Apply applies the render options to the renderer
func (r *Renderer) Apply() error {
	opts := r.Options
	if opts.OutputFormat != "" {
		if err := r.SetRenderOption("render-output-format", opts.OutputFormat); err != nil {
			return err
		}
	}

	if opts.ImageWidth > 0 {
		if err := r.SetRenderOptionInt("render-image-width", opts.ImageWidth); err != nil {
			return err
		}
	}

	if opts.ImageHeight > 0 {
		if err := r.SetRenderOptionInt("render-image-height", opts.ImageHeight); err != nil {
			return err
		}
	}

	if opts.BackgroundColor != "" {
		if err := r.SetRenderOption("render-background-color", opts.BackgroundColor); err != nil {
			return err
		}
	}

	if opts.BondLength > 0 {
		if err := r.SetRenderOptionInt("render-bond-length", opts.BondLength); err != nil {
			return err
		}
	}

	if opts.RelativeThickness > 0 {
		if err := r.SetRenderOptionFloat("render-relative-thickness", opts.RelativeThickness); err != nil {
			return err
		}
	}

	if err := r.SetRenderOptionBool("render-atom-ids-visible", opts.ShowAtomIDs); err != nil {
		return err
	}

	if err := r.SetRenderOptionBool("render-bond-ids-visible", opts.ShowBondIDs); err != nil {
		return err
	}

	if opts.Margins != "" {
		if err := r.SetRenderOption("render-margins", opts.Margins); err != nil {
			return err
		}
	}

	if opts.StereoStyle != "" {
		if err := r.SetRenderOption("render-stereo-style", opts.StereoStyle); err != nil {
			return err
		}
	}

	if opts.LabelMode != "" {
		if err := r.SetRenderOption("render-label-mode", opts.LabelMode); err != nil {
			return err
		}
	}

	return nil
}

// getLastError retrieves the last error message from Indigo
func getLastError() string {
	errMsg := C.indigoGetLastError()
	if errMsg == nil {
		return "unknown error"
	}
	return C.GoString(errMsg)
}
