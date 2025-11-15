// Package reaction provides chemical reaction manipulation using Indigo library via CGO
// coding=utf-8
// @Project : go-indigo
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction.go
// @Software: GoLand
package reaction

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd

// Windows: link against import libraries (.lib)
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../3rd/windows-x86_64 -lindigo
#cgo windows,386 LDFLAGS: -L${SRCDIR}/../3rd/windows-i386 -lindigo

// Linux: use $ORIGIN for runtime library search
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -Wl,-rpath='$ORIGIN/../3rd/linux-x86_64'
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../3rd/linux-aarch64 -lindigo -Wl,-rpath='$ORIGIN/../3rd/linux-aarch64'

// macOS: use @loader_path (not @executable_path) for shared libraries
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-x86_64 -lindigo -Wl,-rpath,'@loader_path/../3rd/darwin-x86_64'
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-aarch64 -lindigo -Wl,-rpath,'@loader_path/../3rd/darwin-aarch64'

#include <stdlib.h>
#include "indigo.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// Reaction center constants
const (
	RC_NOT_CENTER     = -1
	RC_UNMARKED       = 0
	RC_CENTER         = 1
	RC_UNCHANGED      = 2
	RC_MADE_OR_BROKEN = 4
	RC_ORDER_CHANGED  = 8
)

// Reaction represents a chemical reaction in Indigo
type Reaction struct {
	Handle int
	Closed bool
}

// Close frees the Indigo reaction object
func (r *Reaction) Close() error {
	if r.Closed || r.Handle < 0 {
		return nil
	}

	ret := int(C.indigoFree(C.int(r.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to free reaction: %s", getLastError())
	}

	r.Closed = true
	r.Handle = -1
	return nil
}

// CountReactants returns the number of reactants in the reaction
func (r *Reaction) CountReactants() (int, error) {
	if r.Closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	count := int(C.indigoCountReactants(C.int(r.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count reactants: %s", getLastError())
	}

	return count, nil
}

// CountProducts returns the number of products in the reaction
func (r *Reaction) CountProducts() (int, error) {
	if r.Closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	count := int(C.indigoCountProducts(C.int(r.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count products: %s", getLastError())
	}

	return count, nil
}

// CountCatalysts returns the number of catalysts in the reaction
func (r *Reaction) CountCatalysts() (int, error) {
	if r.Closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	count := int(C.indigoCountCatalysts(C.int(r.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count catalysts: %s", getLastError())
	}

	return count, nil
}

// CountMolecules returns the total number of molecules (reactants + products + catalysts)
func (r *Reaction) CountMolecules() (int, error) {
	if r.Closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	count := int(C.indigoCountMolecules(C.int(r.Handle)))
	if count < 0 {
		return 0, fmt.Errorf("failed to count molecules: %s", getLastError())
	}

	return count, nil
}

// AddReactant adds a molecule as a reactant to the reaction
// moleculeHandle is the Indigo handle of the molecule to add
func (r *Reaction) AddReactant(moleculeHandle int) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoAddReactant(C.int(r.Handle), C.int(moleculeHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to add reactant: %s", getLastError())
	}

	return nil
}

// AddProduct adds a molecule as a product to the reaction
// moleculeHandle is the Indigo handle of the molecule to add
func (r *Reaction) AddProduct(moleculeHandle int) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoAddProduct(C.int(r.Handle), C.int(moleculeHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to add product: %s", getLastError())
	}

	return nil
}

// AddCatalyst adds a molecule as a catalyst to the reaction
// moleculeHandle is the Indigo handle of the molecule to add
func (r *Reaction) AddCatalyst(moleculeHandle int) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoAddCatalyst(C.int(r.Handle), C.int(moleculeHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to add catalyst: %s", getLastError())
	}

	return nil
}

// GetMolecule returns a molecule from the reaction by index
// Index order: reactants, then products, then catalysts
func (r *Reaction) GetMolecule(index int) (int, error) {
	if r.Closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	handle := int(C.indigoGetMolecule(C.int(r.Handle), C.int(index)))
	if handle < 0 {
		return 0, fmt.Errorf("failed to get molecule at index %d: %s", index, getLastError())
	}

	return handle, nil
}

// Clone creates a deep copy of the reaction
func (r *Reaction) Clone() (*Reaction, error) {
	if r.Closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	newHandle := int(C.indigoClone(C.int(r.Handle)))
	if newHandle < 0 {
		return nil, fmt.Errorf("failed to clone reaction: %s", getLastError())
	}

	return newReaction(newHandle), nil
}

// Optimize optimizes the query reaction for faster substructure search
// options is a string with optimization options (empty string for defaults)
func (r *Reaction) Optimize(options string) error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	cOptions := C.CString(options)
	defer C.free(unsafe.Pointer(cOptions))

	ret := int(C.indigoOptimize(C.int(r.Handle), cOptions))
	if ret < 0 {
		return fmt.Errorf("failed to optimize reaction: %s", getLastError())
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

// newReaction is a helper function to create a Reaction object from a handle
// It sets up the finalizer to ensure proper cleanup
func newReaction(handle int) *Reaction {
	r := &Reaction{
		Handle: handle,
		Closed: false,
	}
	runtime.SetFinalizer(r, (*Reaction).Close)
	return r
}
