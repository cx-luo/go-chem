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

// Linux platforms
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -lindigo-renderer -Wl,-rpath,${SRCDIR}/../3rd/linux-x86_64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../3rd/linux-aarch64 -lindigo -lindigo-renderer -Wl,-rpath,${SRCDIR}/../3rd/linux-aarch64

// macOS platforms
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-x86_64 -lindigo -lindigo-renderer -Wl,-rpath,${SRCDIR}/../3rd/darwin-x86_64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-aarch64 -lindigo -lindigo-renderer -Wl,-rpath,${SRCDIR}/../3rd/darwin-aarch64

#include <stdlib.h>
#include "indigo.h"
#include "indigo-renderer.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// indigoSessionID holds the session ID for Indigo (defined in molecule package)
var indigoSessionID C.qword

// rendererInitialized tracks whether the renderer has been initialized
var rendererInitialized = false

func init() {
	// Initialize Indigo session if not already initialized
	indigoSessionID = C.indigoAllocSessionId()
	C.indigoSetSessionId(indigoSessionID)
}

// InitRenderer initializes the Indigo renderer for the current session
// This should be called before using rendering functions
func InitRenderer() error {
	if rendererInitialized {
		return nil // Already initialized
	}

	ret := int(C.indigoRendererInit(indigoSessionID))
	if ret < 0 {
		return fmt.Errorf("failed to initialize renderer: %s", getLastError())
	}

	rendererInitialized = true
	return nil
}

// DisposeRenderer disposes the Indigo renderer
// This should be called when done using rendering functions
func DisposeRenderer() error {
	if !rendererInitialized {
		return nil // Not initialized
	}

	ret := int(C.indigoRendererDispose(indigoSessionID))
	if ret < 0 {
		return fmt.Errorf("failed to dispose renderer: %s", getLastError())
	}

	rendererInitialized = false
	return nil
}

// ResetRenderer resets all rendering settings to defaults
func ResetRenderer() error {
	ret := int(C.indigoRenderReset())
	if ret < 0 {
		return fmt.Errorf("failed to reset renderer: %s", getLastError())
	}
	return nil
}

// RenderToFile renders an object (molecule or reaction) to a file
// objectHandle: the Indigo handle of the object to render
// filename: the output file path (e.g., "molecule.png", "reaction.svg")
func RenderToFile(objectHandle int, filename string) error {
	if objectHandle < 0 {
		return fmt.Errorf("invalid object handle")
	}

	// Ensure renderer is initialized
	if !rendererInitialized {
		if err := InitRenderer(); err != nil {
			return fmt.Errorf("failed to initialize renderer: %w", err)
		}
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
func Render(objectHandle int, outputHandle int) error {
	if objectHandle < 0 {
		return fmt.Errorf("invalid object handle")
	}
	if outputHandle < 0 {
		return fmt.Errorf("invalid output handle")
	}

	// Ensure renderer is initialized
	if !rendererInitialized {
		if err := InitRenderer(); err != nil {
			return fmt.Errorf("failed to initialize renderer: %w", err)
		}
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
func RenderGridToFile(arrayHandle int, refAtoms []int, nColumns int, filename string) error {
	if arrayHandle < 0 {
		return fmt.Errorf("invalid array handle")
	}
	if nColumns <= 0 {
		return fmt.Errorf("invalid number of columns: %d", nColumns)
	}

	// Ensure renderer is initialized
	if !rendererInitialized {
		if err := InitRenderer(); err != nil {
			return fmt.Errorf("failed to initialize renderer: %w", err)
		}
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
func RenderGrid(arrayHandle int, refAtoms []int, nColumns int, outputHandle int) error {
	if arrayHandle < 0 {
		return fmt.Errorf("invalid array handle")
	}
	if nColumns <= 0 {
		return fmt.Errorf("invalid number of columns: %d", nColumns)
	}
	if outputHandle < 0 {
		return fmt.Errorf("invalid output handle")
	}

	// Ensure renderer is initialized
	if !rendererInitialized {
		if err := InitRenderer(); err != nil {
			return fmt.Errorf("failed to initialize renderer: %w", err)
		}
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
func SetRenderOption(option string, value string) error {
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
func SetRenderOptionInt(option string, value int) error {
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
func SetRenderOptionFloat(option string, value float64) error {
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
func SetRenderOptionBool(option string, value bool) error {
	strValue := "false"
	if value {
		strValue = "true"
	}
	return SetRenderOption(option, strValue)
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

// DefaultRenderOptions returns default rendering options
func DefaultRenderOptions() *RenderOptions {
	return &RenderOptions{
		OutputFormat:      "png",
		ImageWidth:        1600,
		ImageHeight:       1600,
		BackgroundColor:   "1.0, 1.0, 1.0",
		BondLength:        40,
		RelativeThickness: 1.0,
		ShowAtomIDs:       false,
		ShowBondIDs:       false,
		Margins:           "10, 10",
		StereoStyle:       "ext",
		LabelMode:         "hetero",
	}
}

// Apply applies the render options to the renderer
func (opts *RenderOptions) Apply() error {
	if opts.OutputFormat != "" {
		if err := SetRenderOption("render-output-format", opts.OutputFormat); err != nil {
			return err
		}
	}

	if opts.ImageWidth > 0 {
		if err := SetRenderOptionInt("render-image-width", opts.ImageWidth); err != nil {
			return err
		}
	}

	if opts.ImageHeight > 0 {
		if err := SetRenderOptionInt("render-image-height", opts.ImageHeight); err != nil {
			return err
		}
	}

	if opts.BackgroundColor != "" {
		if err := SetRenderOption("render-background-color", opts.BackgroundColor); err != nil {
			return err
		}
	}

	if opts.BondLength > 0 {
		if err := SetRenderOptionInt("render-bond-length", opts.BondLength); err != nil {
			return err
		}
	}

	if opts.RelativeThickness > 0 {
		if err := SetRenderOptionFloat("render-relative-thickness", opts.RelativeThickness); err != nil {
			return err
		}
	}

	if err := SetRenderOptionBool("render-atom-ids-visible", opts.ShowAtomIDs); err != nil {
		return err
	}

	if err := SetRenderOptionBool("render-bond-ids-visible", opts.ShowBondIDs); err != nil {
		return err
	}

	if opts.Margins != "" {
		if err := SetRenderOption("render-margins", opts.Margins); err != nil {
			return err
		}
	}

	if opts.StereoStyle != "" {
		if err := SetRenderOption("render-stereo-style", opts.StereoStyle); err != nil {
			return err
		}
	}

	if opts.LabelMode != "" {
		if err := SetRenderOption("render-label-mode", opts.LabelMode); err != nil {
			return err
		}
	}

	return nil
}

// CreateArray creates an Indigo array for rendering multiple objects
func CreateArray() (int, error) {
	handle := int(C.indigoCreateArray())
	if handle < 0 {
		return 0, fmt.Errorf("failed to create array: %s", getLastError())
	}
	return handle, nil
}

// ArrayAdd adds an object to an array
func ArrayAdd(arrayHandle int, objectHandle int) error {
	if arrayHandle < 0 {
		return fmt.Errorf("invalid array handle")
	}
	if objectHandle < 0 {
		return fmt.Errorf("invalid object handle")
	}

	ret := int(C.indigoArrayAdd(C.int(arrayHandle), C.int(objectHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to add object to array: %s", getLastError())
	}

	return nil
}

// FreeObject frees an Indigo object (array, buffer, etc.)
func FreeObject(handle int) error {
	if handle < 0 {
		return nil // Already invalid
	}

	ret := int(C.indigoFree(C.int(handle)))
	if ret < 0 {
		return fmt.Errorf("failed to free object: %s", getLastError())
	}

	return nil
}

// CreateWriteBuffer creates an output buffer for rendering
func CreateWriteBuffer() (int, error) {
	handle := int(C.indigoWriteBuffer())
	if handle < 0 {
		return 0, fmt.Errorf("failed to create write buffer: %s", getLastError())
	}

	runtime.SetFinalizer(&handle, func(h *int) {
		if *h >= 0 {
			C.indigoFree(C.int(*h))
		}
	})

	return handle, nil
}

// GetBufferData retrieves data from a write buffer
func GetBufferData(bufferHandle int) ([]byte, error) {
	if bufferHandle < 0 {
		return nil, fmt.Errorf("invalid buffer handle")
	}

	var size C.int
	var dataPtr *C.char
	ret := C.indigoToBuffer(C.int(bufferHandle), &dataPtr, &size)
	if ret < 0 || dataPtr == nil {
		return nil, fmt.Errorf("failed to get buffer data: %s", getLastError())
	}

	// Copy C data to Go slice
	// Note: dataPtr is managed by Indigo internally, don't free it
	data := C.GoBytes(unsafe.Pointer(dataPtr), size)
	return data, nil
}

// getLastError retrieves the last error message from Indigo
func getLastError() string {
	errMsg := C.indigoGetLastError()
	if errMsg == nil {
		return "unknown error"
	}
	return C.GoString(errMsg)
}
