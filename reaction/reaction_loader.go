// Package reaction provides reaction loading functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-indigo
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction_loader.go
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

//// LoadReactionFromString loads a reaction from a string
//// The string should contain reaction data in a supported format (RXN, SMILES, etc.)
//func LoadReactionFromString(data string) (*Reaction, error) {
//	cData := C.CString(data)
//	defer C.free(unsafe.Pointer(cData))
//
//	handle := int(C.indigoLoadReactionFromString(cData))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load reaction from string: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadReactionFromFile loads a reaction from a file
//func LoadReactionFromFile(filename string) (*Reaction, error) {
//	cFilename := C.CString(filename)
//	defer C.free(unsafe.Pointer(cFilename))
//
//	handle := int(C.indigoLoadReactionFromFile(cFilename))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load reaction from file %s: %s", filename, getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadReactionFromBuffer loads a reaction from a byte buffer
//func LoadReactionFromBuffer(buffer []byte) (*Reaction, error) {
//	if len(buffer) == 0 {
//		return nil, fmt.Errorf("empty buffer")
//	}
//
//	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
//	handle := int(C.indigoLoadReactionFromBuffer(cBuffer, C.int(len(buffer))))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load reaction from buffer: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadQueryReactionFromString loads a query reaction from a string
//func LoadQueryReactionFromString(data string) (*Reaction, error) {
//	cData := C.CString(data)
//	defer C.free(unsafe.Pointer(cData))
//
//	handle := int(C.indigoLoadQueryReactionFromString(cData))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load query reaction from string: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadQueryReactionFromFile loads a query reaction from a file
//func LoadQueryReactionFromFile(filename string) (*Reaction, error) {
//	cFilename := C.CString(filename)
//	defer C.free(unsafe.Pointer(cFilename))
//
//	handle := int(C.indigoLoadQueryReactionFromFile(cFilename))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load query reaction from file %s: %s", filename, getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadQueryReactionFromBuffer loads a query reaction from a byte buffer
//func LoadQueryReactionFromBuffer(buffer []byte) (*Reaction, error) {
//	if len(buffer) == 0 {
//		return nil, fmt.Errorf("empty buffer")
//	}
//
//	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
//	handle := int(C.indigoLoadQueryReactionFromBuffer(cBuffer, C.int(len(buffer))))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load query reaction from buffer: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadReactionSmartsFromString loads a reaction SMARTS from a string
//func LoadReactionSmartsFromString(smarts string) (*Reaction, error) {
//	cSmarts := C.CString(smarts)
//	defer C.free(unsafe.Pointer(cSmarts))
//
//	handle := int(C.indigoLoadReactionSmartsFromString(cSmarts))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load reaction SMARTS from string: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadReactionSmartsFromFile loads a reaction SMARTS from a file
//func LoadReactionSmartsFromFile(filename string) (*Reaction, error) {
//	cFilename := C.CString(filename)
//	defer C.free(unsafe.Pointer(cFilename))
//
//	handle := int(C.indigoLoadReactionSmartsFromFile(cFilename))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load reaction SMARTS from file %s: %s", filename, getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadReactionSmartsFromBuffer loads a reaction SMARTS from a byte buffer
//func LoadReactionSmartsFromBuffer(buffer []byte) (*Reaction, error) {
//	if len(buffer) == 0 {
//		return nil, fmt.Errorf("empty buffer")
//	}
//
//	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
//	handle := int(C.indigoLoadReactionSmartsFromBuffer(cBuffer, C.int(len(buffer))))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load reaction SMARTS from buffer: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadReactionWithLibFromString loads a reaction from a string with a monomer library
//func LoadReactionWithLibFromString(data string, monomerLibraryHandle int) (*Reaction, error) {
//	cData := C.CString(data)
//	defer C.free(unsafe.Pointer(cData))
//
//	handle := int(C.indigoLoadReactionWithLibFromString(cData, C.int(monomerLibraryHandle)))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load reaction with library from string: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadReactionWithLibFromFile loads a reaction from a file with a monomer library
//func LoadReactionWithLibFromFile(filename string, monomerLibraryHandle int) (*Reaction, error) {
//	cFilename := C.CString(filename)
//	defer C.free(unsafe.Pointer(cFilename))
//
//	handle := int(C.indigoLoadReactionWithLibFromFile(cFilename, C.int(monomerLibraryHandle)))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load reaction with library from file %s: %s", filename, getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadReactionWithLibFromBuffer loads a reaction from a byte buffer with a monomer library
//func LoadReactionWithLibFromBuffer(buffer []byte, monomerLibraryHandle int) (*Reaction, error) {
//	if len(buffer) == 0 {
//		return nil, fmt.Errorf("empty buffer")
//	}
//
//	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
//	handle := int(C.indigoLoadReactionWithLibFromBuffer(cBuffer, C.int(len(buffer)), C.int(monomerLibraryHandle)))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load reaction with library from buffer: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadQueryReactionWithLibFromString loads a query reaction from a string with a monomer library
//func LoadQueryReactionWithLibFromString(data string, monomerLibraryHandle int) (*Reaction, error) {
//	cData := C.CString(data)
//	defer C.free(unsafe.Pointer(cData))
//
//	handle := int(C.indigoLoadQueryReactionWithLibFromString(cData, C.int(monomerLibraryHandle)))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load query reaction with library from string: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadQueryReactionWithLibFromFile loads a query reaction from a file with a monomer library
//func LoadQueryReactionWithLibFromFile(filename string, monomerLibraryHandle int) (*Reaction, error) {
//	cFilename := C.CString(filename)
//	defer C.free(unsafe.Pointer(cFilename))
//
//	handle := int(C.indigoLoadQueryReactionWithLibFromFile(cFilename, C.int(monomerLibraryHandle)))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load query reaction with library from file %s: %s", filename, getLastError())
//	}
//
//	return newReaction(handle), nil
//}
//
//// LoadQueryReactionWithLibFromBuffer loads a query reaction from a byte buffer with a monomer library
//func LoadQueryReactionWithLibFromBuffer(buffer []byte, monomerLibraryHandle int) (*Reaction, error) {
//	if len(buffer) == 0 {
//		return nil, fmt.Errorf("empty buffer")
//	}
//
//	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
//	handle := int(C.indigoLoadQueryReactionWithLibFromBuffer(cBuffer, C.int(len(buffer)), C.int(monomerLibraryHandle)))
//	if handle < 0 {
//		return nil, fmt.Errorf("failed to load query reaction with library from buffer: %s", getLastError())
//	}
//
//	return newReaction(handle), nil
//}
