// Package reaction provides helper methods for reactions using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/08
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction_helpers.go
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

	"github.com/cx-luo/go-chem/molecule"
)

// HasNext checks if an iterator has more items
func HasNext(iterHandle int) bool {
	return C.indigoHasNext(C.int(iterHandle)) != 0
}

// Next moves to the next item in an iterator and returns its handle
func Next(iterHandle int) (int, error) {
	handle := int(C.indigoNext(C.int(iterHandle)))
	if handle < 0 {
		return 0, fmt.Errorf("failed to get next item: %s", getLastError())
	}
	return handle, nil
}

// FreeIterator frees an iterator handle
func FreeIterator(iterHandle int) error {
	if iterHandle < 0 {
		return nil
	}

	ret := int(C.indigoFree(C.int(iterHandle)))
	if ret < 0 {
		return fmt.Errorf("failed to free iterator: %s", getLastError())
	}

	return nil
}

// GetReactant returns a reactant molecule handle by index
func (r *Reaction) GetReactant(index int) (int, error) {
	if r.closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	// Get iterator
	iterHandle := int(C.indigoIterateReactants(C.int(r.handle)))
	if iterHandle < 0 {
		return 0, fmt.Errorf("failed to iterate reactants: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	// Find the reactant at the specified index
	currentIndex := 0
	for HasNext(iterHandle) {
		molHandle, err := Next(iterHandle)
		if err != nil {
			return 0, err
		}
		if currentIndex == index {
			return molHandle, nil
		}
		currentIndex++
	}

	return 0, fmt.Errorf("reactant index %d out of range", index)
}

// GetProduct returns a product molecule handle by index
func (r *Reaction) GetProduct(index int) (int, error) {
	if r.closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	// Get iterator
	iterHandle := int(C.indigoIterateProducts(C.int(r.handle)))
	if iterHandle < 0 {
		return 0, fmt.Errorf("failed to iterate products: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	// Find the product at the specified index
	currentIndex := 0
	for HasNext(iterHandle) {
		molHandle, err := Next(iterHandle)
		if err != nil {
			return 0, err
		}
		if currentIndex == index {
			return molHandle, nil
		}
		currentIndex++
	}

	return 0, fmt.Errorf("product index %d out of range", index)
}

// GetCatalyst returns a catalyst molecule handle by index
func (r *Reaction) GetCatalyst(index int) (int, error) {
	if r.closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	// Get iterator
	iterHandle := int(C.indigoIterateCatalysts(C.int(r.handle)))
	if iterHandle < 0 {
		return 0, fmt.Errorf("failed to iterate catalysts: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	// Find the catalyst at the specified index
	currentIndex := 0
	for HasNext(iterHandle) {
		molHandle, err := Next(iterHandle)
		if err != nil {
			return 0, err
		}
		if currentIndex == index {
			return molHandle, nil
		}
		currentIndex++
	}

	return 0, fmt.Errorf("catalyst index %d out of range", index)
}

// Layout performs 2D layout of the reaction
func (r *Reaction) Layout() error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoLayout(C.int(r.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to layout reaction: %s", getLastError())
	}

	return nil
}

// Clean2D performs 2D cleaning of the reaction
func (r *Reaction) Clean2D() error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoClean2d(C.int(r.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to clean2d reaction: %s", getLastError())
	}

	return nil
}

// Aromatize performs aromatization of all molecules in the reaction
func (r *Reaction) Aromatize() error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoAromatize(C.int(r.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to aromatize reaction: %s", getLastError())
	}

	return nil
}

// Dearomatize removes aromaticity from all molecules in the reaction
func (r *Reaction) Dearomatize() error {
	if r.closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoDearomatize(C.int(r.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to dearomatize reaction: %s", getLastError())
	}

	return nil
}

// GetReactantMolecule returns a reactant molecule as a Molecule object by index
func (r *Reaction) GetReactantMolecule(index int) (*molecule.Molecule, error) {
	handle, err := r.GetReactant(index)
	if err != nil {
		return nil, err
	}

	// Clone the molecule to avoid issues with handle ownership
	clonedHandle := int(C.indigoClone(C.int(handle)))
	if clonedHandle < 0 {
		return nil, fmt.Errorf("failed to clone reactant: %s", getLastError())
	}

	return molecule.NewMoleculeFromHandle(clonedHandle)
}

// GetProductMolecule returns a product molecule as a Molecule object by index
func (r *Reaction) GetProductMolecule(index int) (*molecule.Molecule, error) {
	handle, err := r.GetProduct(index)
	if err != nil {
		return nil, err
	}

	// Clone the molecule to avoid issues with handle ownership
	clonedHandle := int(C.indigoClone(C.int(handle)))
	if clonedHandle < 0 {
		return nil, fmt.Errorf("failed to clone product: %s", getLastError())
	}

	return molecule.NewMoleculeFromHandle(clonedHandle)
}

// GetCatalystMolecule returns a catalyst molecule as a Molecule object by index
func (r *Reaction) GetCatalystMolecule(index int) (*molecule.Molecule, error) {
	handle, err := r.GetCatalyst(index)
	if err != nil {
		return nil, err
	}

	// Clone the molecule to avoid issues with handle ownership
	clonedHandle := int(C.indigoClone(C.int(handle)))
	if clonedHandle < 0 {
		return nil, fmt.Errorf("failed to clone catalyst: %s", getLastError())
	}

	return molecule.NewMoleculeFromHandle(clonedHandle)
}

// GetAllReactants returns all reactant molecules as Molecule objects
func (r *Reaction) GetAllReactants() ([]*molecule.Molecule, error) {
	if r.closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	count, err := r.CountReactants()
	if err != nil {
		return nil, err
	}

	reactants := make([]*molecule.Molecule, 0, count)

	iterHandle := int(C.indigoIterateReactants(C.int(r.handle)))
	if iterHandle < 0 {
		return nil, fmt.Errorf("failed to iterate reactants: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	for C.indigoHasNext(C.int(iterHandle)) != 0 {
		molHandle := int(C.indigoNext(C.int(iterHandle)))
		if molHandle < 0 {
			continue
		}

		// Clone the molecule to avoid ownership issues
		clonedHandle := int(C.indigoClone(C.int(molHandle)))
		if clonedHandle < 0 {
			continue
		}

		mol, err := molecule.NewMoleculeFromHandle(clonedHandle)
		if err != nil {
			C.indigoFree(C.int(clonedHandle))
			continue
		}

		reactants = append(reactants, mol)
	}

	return reactants, nil
}

// GetAllProducts returns all product molecules as Molecule objects
func (r *Reaction) GetAllProducts() ([]*molecule.Molecule, error) {
	if r.closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	count, err := r.CountProducts()
	if err != nil {
		return nil, err
	}

	products := make([]*molecule.Molecule, 0, count)

	iterHandle := int(C.indigoIterateProducts(C.int(r.handle)))
	if iterHandle < 0 {
		return nil, fmt.Errorf("failed to iterate products: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	for C.indigoHasNext(C.int(iterHandle)) != 0 {
		molHandle := int(C.indigoNext(C.int(iterHandle)))
		if molHandle < 0 {
			continue
		}

		// Clone the molecule to avoid ownership issues
		clonedHandle := int(C.indigoClone(C.int(molHandle)))
		if clonedHandle < 0 {
			continue
		}

		mol, err := molecule.NewMoleculeFromHandle(clonedHandle)
		if err != nil {
			C.indigoFree(C.int(clonedHandle))
			continue
		}

		products = append(products, mol)
	}

	return products, nil
}

// GetAllCatalysts returns all catalyst molecules as Molecule objects
func (r *Reaction) GetAllCatalysts() ([]*molecule.Molecule, error) {
	if r.closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	count, err := r.CountCatalysts()
	if err != nil {
		return nil, err
	}

	catalysts := make([]*molecule.Molecule, 0, count)

	iterHandle := int(C.indigoIterateCatalysts(C.int(r.handle)))
	if iterHandle < 0 {
		return nil, fmt.Errorf("failed to iterate catalysts: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	for C.indigoHasNext(C.int(iterHandle)) != 0 {
		molHandle := int(C.indigoNext(C.int(iterHandle)))
		if molHandle < 0 {
			continue
		}

		// Clone the molecule to avoid ownership issues
		clonedHandle := int(C.indigoClone(C.int(molHandle)))
		if clonedHandle < 0 {
			continue
		}

		mol, err := molecule.NewMoleculeFromHandle(clonedHandle)
		if err != nil {
			C.indigoFree(C.int(clonedHandle))
			continue
		}

		catalysts = append(catalysts, mol)
	}

	return catalysts, nil
}

func (r *Reaction) GetAllMolecules() ([]*molecule.Molecule, error) {
	if r.closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	count, err := r.CountMolecules()
	if err != nil {
		return nil, err
	}

	molecules := make([]*molecule.Molecule, 0, count)

	iterHandle := int(C.indigoIterateMolecules(C.int(r.handle)))
	if iterHandle < 0 {
		return nil, fmt.Errorf("failed to iterate molecules: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	for C.indigoHasNext(C.int(iterHandle)) != 0 {
		molHandle := int(C.indigoNext(C.int(iterHandle)))
		if molHandle < 0 {
			continue
		}

		// Clone the molecule to avoid ownership issues
		clonedHandle := int(C.indigoClone(C.int(molHandle)))
		if clonedHandle < 0 {
			continue
		}

		mol, err := molecule.NewMoleculeFromHandle(clonedHandle)
		if err != nil {
			C.indigoFree(C.int(clonedHandle))
			continue
		}

		molecules = append(molecules, mol)
	}

	return molecules, nil
}
