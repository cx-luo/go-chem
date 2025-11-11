// Package molecule provides molecule saving functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_saver.go
// @Software: GoLand
package molecule

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd

// Windows platforms
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../3rd/windows-x86_64 -lindigo
#cgo windows,386 LDFLAGS: -L${SRCDIR}/../3rd/windows-i386 -lindigo

// Linux platforms
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux-x86_64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../3rd/linux-aarch64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux-aarch64

// macOS platforms
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-x86_64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/darwin-x86_64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-aarch64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/darwin-aarch64

#include <stdlib.h>
#include "indigo.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ToSmiles converts the molecule to SMILES format
func (m *Molecule) ToSmiles() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoSmiles(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCanonicalSmiles converts the molecule to canonical SMILES format
func (m *Molecule) ToCanonicalSmiles() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoCanonicalSmiles(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to canonical SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToSmarts converts the molecule to SMARTS format
func (m *Molecule) ToSmarts() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoSmarts(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to SMARTS: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCanonicalSmarts converts the molecule to canonical SMARTS format
func (m *Molecule) ToCanonicalSmarts() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoCanonicalSmarts(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to canonical SMARTS: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToMolfile returns the molecule as a MOL file string
func (m *Molecule) ToMolfile() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoMolfile(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to MOL: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToFile saves the molecule to a file in MOL format
func (m *Molecule) SaveToFile(filename string) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveMolfileToFile(C.int(m.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save to file %s: %s", filename, getLastError())
	}

	return nil
}

// SaveMolfile saves the molecule to an output object
// outputHandle is the Indigo handle of an output object
func (m *Molecule) SaveMolfile(outputHandle int) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	ret := int(C.indigoSaveMolfile(C.int(m.Handle), C.int(outputHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to save to output: %s", getLastError())
	}

	return nil
}

// ToJSON converts the molecule to JSON format
func (m *Molecule) ToJSON() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	// Create a string output buffer
	bufferHandle := int(C.indigoWriteBuffer())
	if bufferHandle < 0 {
		return "", fmt.Errorf("failed to create buffer: %s", getLastError())
	}
	defer C.indigoFree(C.int(bufferHandle))

	// Save as JSON
	ret := int(C.indigoSaveJson(C.int(m.Handle), C.int(bufferHandle)))
	if ret < 0 {
		return "", fmt.Errorf("failed to save as JSON: %s", getLastError())
	}

	// Get the string
	cStr := C.indigoToString(C.int(bufferHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get JSON string: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToJSONFile saves the molecule to a file in JSON format
func (m *Molecule) SaveToJSONFile(filename string) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveJsonToFile(C.int(m.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save to JSON file %s: %s", filename, getLastError())
	}

	return nil
}

// ToBase64String converts the molecule to base64 string
func (m *Molecule) ToBase64String() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoToBase64String(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to base64: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCXSmiles converts the molecule to ChemAxon Extended SMILES format
func (m *Molecule) ToCXSmiles() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	// Set the SMILES saving format to chemaxon
	cOption := C.CString("smiles-saving-format")
	cValue := C.CString("chemaxon")
	cDefaultValue := C.CString("daylight")
	defer C.free(unsafe.Pointer(cOption))
	defer C.free(unsafe.Pointer(cValue))
	defer C.free(unsafe.Pointer(cDefaultValue))

	ret := int(C.indigoSetOption(cOption, cValue))
	if ret < 0 {
		return "", fmt.Errorf("failed to set chemaxon format: %s", getLastError())
	}

	// Reset to default after getting the SMILES
	defer C.indigoSetOption(cOption, cDefaultValue)

	cStr := C.indigoSmiles(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to CXSmiles: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCanonicalCXSmiles converts the molecule to canonical ChemAxon Extended SMILES format
// Note: According to Indigo API, canonical SMILES already includes ChemAxon extensions
func (m *Molecule) ToCanonicalCXSmiles() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	// For canonical CXSMILES, we use canonicalSmiles directly
	// The canonical SMILES in Indigo automatically includes ChemAxon extensions
	cStr := C.indigoCanonicalSmiles(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to canonical CXSmiles: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCML converts the molecule to CML (Chemical Markup Language) format
func (m *Molecule) ToCML() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoCml(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to CML: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCDXML converts the molecule to CDXML format
func (m *Molecule) ToCDXML() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	cStr := C.indigoCdxml(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to CDXML: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCDXBase64 converts the molecule to base64-encoded CDX format
func (m *Molecule) ToCDXBase64() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	// Create a buffer for CDX
	bufferHandle := int(C.indigoWriteBuffer())
	if bufferHandle < 0 {
		return "", fmt.Errorf("failed to create buffer: %s", getLastError())
	}
	defer C.indigoFree(C.int(bufferHandle))

	// Save as CDX
	ret := int(C.indigoSaveCdx(C.int(m.Handle), C.int(bufferHandle)))
	if ret < 0 {
		return "", fmt.Errorf("failed to save as CDX: %s", getLastError())
	}

	// Convert to base64
	cStr := C.indigoToBase64String(C.int(bufferHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert CDX to base64: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToSDF saves the molecule components to an SDF file
func (m *Molecule) SaveToSDF(filename string) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	// Create a file output
	outputHandle := int(C.indigoWriteFile(cFilename))
	if outputHandle < 0 {
		return fmt.Errorf("failed to create output file: %s", getLastError())
	}
	defer C.indigoFree(C.int(outputHandle))

	// Create SDF saver
	cFormat := C.CString("sdf")
	defer C.free(unsafe.Pointer(cFormat))

	saverHandle := int(C.indigoCreateSaver(C.int(outputHandle), cFormat))
	if saverHandle < 0 {
		return fmt.Errorf("failed to create SDF saver: %s", getLastError())
	}
	defer C.indigoFree(C.int(saverHandle))

	// Iterate through components and save each
	iterHandle := int(C.indigoIterateComponents(C.int(m.Handle)))
	if iterHandle < 0 {
		return fmt.Errorf("failed to iterate components: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	for C.indigoHasNext(C.int(iterHandle)) != 0 {
		compHandle := int(C.indigoNext(C.int(iterHandle)))
		if compHandle < 0 {
			return fmt.Errorf("failed to get next component: %s", getLastError())
		}

		cloneHandle := int(C.indigoClone(C.int(compHandle)))
		if cloneHandle < 0 {
			return fmt.Errorf("failed to clone component: %s", getLastError())
		}

		// Use indigoAppend instead of indigoSdfAppend (following Python API pattern)
		ret := int(C.indigoAppend(C.int(saverHandle), C.int(cloneHandle)))
		C.indigoFree(C.int(cloneHandle))
		if ret < 0 {
			return fmt.Errorf("failed to append to SDF: %s", getLastError())
		}
	}

	// Close the saver
	ret := int(C.indigoClose(C.int(saverHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to close SDF saver: %s", getLastError())
	}

	return nil
}

// ToSDF converts the molecule to SDF format string
func (m *Molecule) ToSDF() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	// Create a buffer for SDF
	bufferHandle := int(C.indigoWriteBuffer())
	if bufferHandle < 0 {
		return "", fmt.Errorf("failed to create buffer: %s", getLastError())
	}
	defer C.indigoFree(C.int(bufferHandle))

	// Create SDF saver
	cFormat := C.CString("sdf")
	defer C.free(unsafe.Pointer(cFormat))

	saverHandle := int(C.indigoCreateSaver(C.int(bufferHandle), cFormat))
	if saverHandle < 0 {
		return "", fmt.Errorf("failed to create SDF saver: %s", getLastError())
	}
	defer C.indigoFree(C.int(saverHandle))

	// Iterate through components and save each
	iterHandle := int(C.indigoIterateComponents(C.int(m.Handle)))
	if iterHandle < 0 {
		return "", fmt.Errorf("failed to iterate components: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	for C.indigoHasNext(C.int(iterHandle)) != 0 {
		compHandle := int(C.indigoNext(C.int(iterHandle)))
		if compHandle < 0 {
			return "", fmt.Errorf("failed to get next component: %s", getLastError())
		}

		cloneHandle := int(C.indigoClone(C.int(compHandle)))
		if cloneHandle < 0 {
			return "", fmt.Errorf("failed to clone component: %s", getLastError())
		}

		// Use indigoAppend instead of indigoSdfAppend (following Python API pattern)
		ret := int(C.indigoAppend(C.int(saverHandle), C.int(cloneHandle)))
		C.indigoFree(C.int(cloneHandle))
		if ret < 0 {
			return "", fmt.Errorf("failed to append to SDF: %s", getLastError())
		}
	}

	// Close the saver
	ret := int(C.indigoClose(C.int(saverHandle)))
	if ret < 0 {
		return "", fmt.Errorf("failed to close SDF saver: %s", getLastError())
	}

	// Get the string
	cStr := C.indigoToString(C.int(bufferHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get SDF string: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToCMLFile saves the molecule to a file in CML format
func (m *Molecule) SaveToCMLFile(filename string) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveCmlToFile(C.int(m.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save to CML file %s: %s", filename, getLastError())
	}

	return nil
}

// SaveToCDXMLFile saves the molecule to a file in CDXML format
func (m *Molecule) SaveToCDXMLFile(filename string) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveCdxmlToFile(C.int(m.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save to CDXML file %s: %s", filename, getLastError())
	}

	return nil
}

// SaveToCDXFile saves the molecule to a file in CDX format
func (m *Molecule) SaveToCDXFile(filename string) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveCdxToFile(C.int(m.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save to CDX file %s: %s", filename, getLastError())
	}

	return nil
}

// ToDaylightSmiles converts the molecule to Daylight SMILES format
// This explicitly uses the Daylight format (as opposed to ChemAxon)
func (m *Molecule) ToDaylightSmiles() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	// Explicitly set Daylight format
	cOption := C.CString("smiles-saving-format")
	cValue := C.CString("daylight")
	defer C.free(unsafe.Pointer(cOption))
	defer C.free(unsafe.Pointer(cValue))

	ret := int(C.indigoSetOption(cOption, cValue))
	if ret < 0 {
		return "", fmt.Errorf("failed to set daylight format: %s", getLastError())
	}

	cStr := C.indigoSmiles(C.int(m.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert to Daylight SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToRDF converts the molecule to RDF (Reaction Data Format) format string
// Note: This is typically used for reactions, but can work with molecules
func (m *Molecule) ToRDF() (string, error) {
	if m.Closed {
		return "", fmt.Errorf("molecule is closed")
	}

	// Create a buffer for RDF
	bufferHandle := int(C.indigoWriteBuffer())
	if bufferHandle < 0 {
		return "", fmt.Errorf("failed to create buffer: %s", getLastError())
	}
	defer C.indigoFree(C.int(bufferHandle))

	// Create RDF saver
	cFormat := C.CString("rdf")
	defer C.free(unsafe.Pointer(cFormat))

	saverHandle := int(C.indigoCreateSaver(C.int(bufferHandle), cFormat))
	if saverHandle < 0 {
		return "", fmt.Errorf("failed to create RDF saver: %s", getLastError())
	}
	defer C.indigoFree(C.int(saverHandle))

	// Append the molecule
	ret := int(C.indigoAppend(C.int(saverHandle), C.int(m.Handle)))
	if ret < 0 {
		return "", fmt.Errorf("failed to append to RDF: %s", getLastError())
	}

	// Close the saver
	ret = int(C.indigoClose(C.int(saverHandle)))
	if ret < 0 {
		return "", fmt.Errorf("failed to close RDF saver: %s", getLastError())
	}

	// Get the string
	cStr := C.indigoToString(C.int(bufferHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get RDF string: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToRDFFile saves the molecule to a file in RDF format
func (m *Molecule) SaveToRDFFile(filename string) error {
	if m.Closed {
		return fmt.Errorf("molecule is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	// Create a file output
	outputHandle := int(C.indigoWriteFile(cFilename))
	if outputHandle < 0 {
		return fmt.Errorf("failed to create output file: %s", getLastError())
	}
	defer C.indigoFree(C.int(outputHandle))

	// Create RDF saver
	cFormat := C.CString("rdf")
	defer C.free(unsafe.Pointer(cFormat))

	saverHandle := int(C.indigoCreateSaver(C.int(outputHandle), cFormat))
	if saverHandle < 0 {
		return fmt.Errorf("failed to create RDF saver: %s", getLastError())
	}
	defer C.indigoFree(C.int(saverHandle))

	// Append the molecule
	ret := int(C.indigoAppend(C.int(saverHandle), C.int(m.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to append to RDF: %s", getLastError())
	}

	// Close the saver
	ret = int(C.indigoClose(C.int(saverHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to close RDF saver: %s", getLastError())
	}

	return nil
}

// ToBuffer converts the molecule to a binary buffer
func (m *Molecule) ToBuffer() ([]byte, error) {
	if m.Closed {
		return nil, fmt.Errorf("molecule is closed")
	}

	// Create a buffer
	bufferHandle := int(C.indigoWriteBuffer())
	if bufferHandle < 0 {
		return nil, fmt.Errorf("failed to create buffer: %s", getLastError())
	}
	defer C.indigoFree(C.int(bufferHandle))

	// Save molecule to buffer
	ret := int(C.indigoSaveMolfile(C.int(m.Handle), C.int(bufferHandle)))
	if ret < 0 {
		return nil, fmt.Errorf("failed to save to buffer: %s", getLastError())
	}

	// Get the buffer content
	cStr := C.indigoToString(C.int(bufferHandle))
	if cStr == nil {
		return nil, fmt.Errorf("failed to get buffer content: %s", getLastError())
	}

	return []byte(C.GoString(cStr)), nil
}

// ToKet converts the molecule to KET (Ketcher JSON) format
// KET format is the JSON format used by Ketcher editor
func (m *Molecule) ToKet() (string, error) {
	return m.ToJSON()
}

// SaveToKetFile saves the molecule to a file in KET (Ketcher JSON) format
func (m *Molecule) SaveToKetFile(filename string) error {
	return m.SaveToJSONFile(filename)
}
