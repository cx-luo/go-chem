// Package reaction provides reaction saving functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction_saver.go
// @Software: GoLand
package reaction

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

// ToRxnfile returns the reaction as an RXN file string
func (r *Reaction) ToRxnfile() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoRxnfile(C.int(r.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to RXN: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToFile saves the reaction to a file in RXN format
func (r *Reaction) SaveToFile(filename string) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveRxnfileToFile(C.int(r.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save reaction to file %s: %s", filename, getLastError())
	}

	return nil
}

// SaveRxnfile saves the reaction to an output object
// outputHandle is the Indigo handle of an output object
func (r *Reaction) SaveRxnfile(outputHandle int) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoSaveRxnfile(C.int(r.Handle), C.int(outputHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to save reaction to output: %s", getLastError())
	}

	return nil
}

// ToSmiles converts the reaction to SMILES format
func (r *Reaction) ToSmiles() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoSmiles(C.int(r.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCanonicalSmiles converts the reaction to canonical SMILES format
func (r *Reaction) ToCanonicalSmiles() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoCanonicalSmiles(C.int(r.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to canonical SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// CreateStringOutput creates a string output buffer
// Returns the Indigo handle for the output buffer
func CreateStringOutput() (int, error) {
	handle := int(C.indigoWriteBuffer())
	if handle < 0 {
		return 0, fmt.Errorf("failed to create string output: %s", getLastError())
	}
	return handle, nil
}

// GetStringOutput retrieves the string from an output buffer
func GetStringOutput(outputHandle int) (string, error) {
	cStr := C.indigoToString(C.int(outputHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get string from output: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// FreeOutput frees an output buffer
func FreeOutput(outputHandle int) error {
	ret := int(C.indigoFree(C.int(outputHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to free output: %s", getLastError())
	}
	return nil
}

// ToSmarts converts the reaction to SMARTS format
func (r *Reaction) ToSmarts() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoSmarts(C.int(r.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to SMARTS: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCanonicalSmarts converts the reaction to canonical SMARTS format
func (r *Reaction) ToCanonicalSmarts() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoCanonicalSmarts(C.int(r.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to canonical SMARTS: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCXSmiles converts the reaction to ChemAxon Extended SMILES format
func (r *Reaction) ToCXSmiles() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
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

	cStr := C.indigoSmiles(C.int(r.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to CXSmiles: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToDaylightSmiles converts the reaction to Daylight SMILES format
// This explicitly uses the Daylight format (as opposed to ChemAxon)
func (r *Reaction) ToDaylightSmiles() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
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

	cStr := C.indigoSmiles(C.int(r.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to Daylight SMILES: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToCML converts the reaction to CML (Chemical Markup Language) format
func (r *Reaction) ToCML() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoCml(C.int(r.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to CML: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToCMLFile saves the reaction to a file in CML format
func (r *Reaction) SaveToCMLFile(filename string) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveCmlToFile(C.int(r.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save reaction to CML file %s: %s", filename, getLastError())
	}

	return nil
}

// ToCDXML converts the reaction to CDXML format
func (r *Reaction) ToCDXML() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
	}

	cStr := C.indigoCdxml(C.int(r.Handle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction to CDXML: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToCDXMLFile saves the reaction to a file in CDXML format
func (r *Reaction) SaveToCDXMLFile(filename string) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveCdxmlToFile(C.int(r.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save reaction to CDXML file %s: %s", filename, getLastError())
	}

	return nil
}

// SaveToCDXFile saves the reaction to a file in CDX format
func (r *Reaction) SaveToCDXFile(filename string) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveCdxToFile(C.int(r.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save reaction to CDX file %s: %s", filename, getLastError())
	}

	return nil
}

// ToCDXBase64 converts the reaction to base64-encoded CDX format
func (r *Reaction) ToCDXBase64() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
	}

	// Create a buffer for CDX
	bufferHandle := int(C.indigoWriteBuffer())
	if bufferHandle < 0 {
		return "", fmt.Errorf("failed to create buffer: %s", getLastError())
	}
	defer C.indigoFree(C.int(bufferHandle))

	// Save as CDX
	ret := int(C.indigoSaveCdx(C.int(r.Handle), C.int(bufferHandle)))
	if ret < 0 {
		return "", fmt.Errorf("failed to save reaction as CDX: %s", getLastError())
	}

	// Convert to base64
	cStr := C.indigoToBase64String(C.int(bufferHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to convert reaction CDX to base64: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// ToJSON converts the reaction to JSON format
func (r *Reaction) ToJSON() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
	}

	// Create a string output buffer
	bufferHandle := int(C.indigoWriteBuffer())
	if bufferHandle < 0 {
		return "", fmt.Errorf("failed to create buffer: %s", getLastError())
	}
	defer C.indigoFree(C.int(bufferHandle))

	// Save as JSON
	ret := int(C.indigoSaveJson(C.int(r.Handle), C.int(bufferHandle)))
	if ret < 0 {
		return "", fmt.Errorf("failed to save reaction as JSON: %s", getLastError())
	}

	// Get the string
	cStr := C.indigoToString(C.int(bufferHandle))
	if cStr == nil {
		return "", fmt.Errorf("failed to get reaction JSON string: %s", getLastError())
	}

	return C.GoString(cStr), nil
}

// SaveToJSONFile saves the reaction to a file in JSON format
func (r *Reaction) SaveToJSONFile(filename string) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := int(C.indigoSaveJsonToFile(C.int(r.Handle), cFilename))
	if ret < 0 {
		return fmt.Errorf("failed to save reaction to JSON file %s: %s", filename, getLastError())
	}

	return nil
}

// ToRDF converts the reaction to RDF (Reaction Data Format) format string
func (r *Reaction) ToRDF() (string, error) {
	if r.Closed {
		return "", fmt.Errorf("reaction is closed")
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

	// Append the reaction
	ret := int(C.indigoAppend(C.int(saverHandle), C.int(r.Handle)))
	if ret < 0 {
		return "", fmt.Errorf("failed to append reaction to RDF: %s", getLastError())
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

// SaveToRDFFile saves the reaction to a file in RDF format
func (r *Reaction) SaveToRDFFile(filename string) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
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

	// Append the reaction
	ret := int(C.indigoAppend(C.int(saverHandle), C.int(r.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to append reaction to RDF: %s", getLastError())
	}

	// Close the saver
	ret = int(C.indigoClose(C.int(saverHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to close RDF saver: %s", getLastError())
	}

	return nil
}

// ToBuffer converts the reaction to a binary buffer
func (r *Reaction) ToBuffer() ([]byte, error) {
	if r.Closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	// Create a buffer
	bufferHandle := int(C.indigoWriteBuffer())
	if bufferHandle < 0 {
		return nil, fmt.Errorf("failed to create buffer: %s", getLastError())
	}
	defer C.indigoFree(C.int(bufferHandle))

	// Save reaction to buffer
	ret := int(C.indigoSaveRxnfile(C.int(r.Handle), C.int(bufferHandle)))
	if ret < 0 {
		return nil, fmt.Errorf("failed to save reaction to buffer: %s", getLastError())
	}

	// Get the buffer content
	cStr := C.indigoToString(C.int(bufferHandle))
	if cStr == nil {
		return nil, fmt.Errorf("failed to get buffer content: %s", getLastError())
	}

	return []byte(C.GoString(cStr)), nil
}

// ToKet converts the reaction to KET (Ketcher JSON) format
// KET format is the JSON format used by Ketcher editor
func (r *Reaction) ToKet() (string, error) {
	return r.ToJSON()
}

// SaveToKetFile saves the reaction to a file in KET (Ketcher JSON) format
func (r *Reaction) SaveToKetFile(filename string) error {
	return r.SaveToJSONFile(filename)
}
