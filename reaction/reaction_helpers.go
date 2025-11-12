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
)

// GetReactant returns a reactant molecule by index
func (r *Reaction) GetReactant(index int) (int, error) {
	if r.Closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	// Get iterator
	iterHandle := int(C.indigoIterateReactants(C.int(r.Handle)))
	if iterHandle < 0 {
		return 0, fmt.Errorf("failed to iterate reactants: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	// Find the reactant at the specified index
	currentIndex := 0
	for C.indigoHasNext(C.int(iterHandle)) != 0 {
		molHandle := int(C.indigoNext(C.int(iterHandle)))
		if currentIndex == index {
			return molHandle, nil
		}
		currentIndex++
	}

	return 0, fmt.Errorf("reactant index %d out of range", index)
}

// GetProduct returns a product molecule by index
func (r *Reaction) GetProduct(index int) (int, error) {
	if r.Closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	// Get iterator
	iterHandle := int(C.indigoIterateProducts(C.int(r.Handle)))
	if iterHandle < 0 {
		return 0, fmt.Errorf("failed to iterate products: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	// Find the product at the specified index
	currentIndex := 0
	for C.indigoHasNext(C.int(iterHandle)) != 0 {
		molHandle := int(C.indigoNext(C.int(iterHandle)))
		if currentIndex == index {
			return molHandle, nil
		}
		currentIndex++
	}

	return 0, fmt.Errorf("product index %d out of range", index)
}

// GetCatalyst returns a catalyst molecule by index
func (r *Reaction) GetCatalyst(index int) (int, error) {
	if r.Closed {
		return 0, fmt.Errorf("reaction is closed")
	}

	// Get iterator
	iterHandle := int(C.indigoIterateCatalysts(C.int(r.Handle)))
	if iterHandle < 0 {
		return 0, fmt.Errorf("failed to iterate catalysts: %s", getLastError())
	}
	defer C.indigoFree(C.int(iterHandle))

	// Find the catalyst at the specified index
	currentIndex := 0
	for C.indigoHasNext(C.int(iterHandle)) != 0 {
		molHandle := int(C.indigoNext(C.int(iterHandle)))
		if currentIndex == index {
			return molHandle, nil
		}
		currentIndex++
	}

	return 0, fmt.Errorf("catalyst index %d out of range", index)
}

// Layout performs 2D layout of the reaction
func (r *Reaction) Layout() error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoLayout(C.int(r.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to layout reaction: %s", getLastError())
	}

	return nil
}

// Clean2D performs 2D cleaning of the reaction
func (r *Reaction) Clean2D() error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoClean2d(C.int(r.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to clean2d reaction: %s", getLastError())
	}

	return nil
}

// Aromatize performs aromatization of all molecules in the reaction
func (r *Reaction) Aromatize() error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoAromatize(C.int(r.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to aromatize reaction: %s", getLastError())
	}

	return nil
}

// Dearomatize removes aromaticity from all molecules in the reaction
func (r *Reaction) Dearomatize() error {
	if r.Closed {
		return fmt.Errorf("reaction is closed")
	}

	ret := int(C.indigoDearomatize(C.int(r.Handle)))
	if ret < 0 {
		return fmt.Errorf("failed to dearomatize reaction: %s", getLastError())
	}

	return nil
}

// GetReactantMolecule returns a reactant molecule handle as a Molecule object by index
func (r *Reaction) GetReactantMolecule(index int) (int, error) {
	handle, err := r.GetReactant(index)
	if err != nil {
		return 0, err
	}

	// Clone the molecule to avoid issues with handle ownership
	clonedHandle := int(C.indigoClone(C.int(handle)))
	if clonedHandle < 0 {
		return 0, fmt.Errorf("failed to clone reactant: %s", getLastError())
	}

	return clonedHandle, nil
}

// GetProductMolecule returns a product molecule as a Molecule object by index
func (r *Reaction) GetProductMolecule(index int) (int, error) {
	handle, err := r.GetProduct(index)
	if err != nil {
		return 0, err
	}

	// Clone the molecule to avoid issues with handle ownership
	clonedHandle := int(C.indigoClone(C.int(handle)))
	if clonedHandle < 0 {
		return 0, fmt.Errorf("failed to clone product: %s", getLastError())
	}

	return clonedHandle, nil
}

// GetCatalystMolecule returns a catalyst molecule as a Molecule object by index
func (r *Reaction) GetCatalystMolecule(index int) (int, error) {
	handle, err := r.GetCatalyst(index)
	if err != nil {
		return 0, err
	}

	// Clone the molecule to avoid issues with handle ownership
	clonedHandle := int(C.indigoClone(C.int(handle)))
	if clonedHandle < 0 {
		return 0, fmt.Errorf("failed to clone catalyst: %s", getLastError())
	}

	return clonedHandle, nil
}

// GetAllReactants returns all reactant molecules as Molecule objects
func (r *Reaction) GetAllReactants() ([]int, error) {
	if r.Closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	count, err := r.CountReactants()
	if err != nil {
		return nil, err
	}

	reactants := make([]int, 0, count)

	iterHandle := int(C.indigoIterateReactants(C.int(r.Handle)))
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

		reactants = append(reactants, clonedHandle)
	}

	return reactants, nil
}

// GetAllProducts returns all product molecules as Molecule objects
func (r *Reaction) GetAllProducts() ([]int, error) {
	if r.Closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	count, err := r.CountProducts()
	if err != nil {
		return nil, err
	}

	products := make([]int, 0, count)

	iterHandle := int(C.indigoIterateProducts(C.int(r.Handle)))
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

		products = append(products, clonedHandle)
	}

	return products, nil
}

// GetAllCatalysts returns all catalyst molecules as Molecule objects
func (r *Reaction) GetAllCatalysts() ([]int, error) {
	if r.Closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	count, err := r.CountCatalysts()
	if err != nil {
		return nil, err
	}

	catalysts := make([]int, 0, count)

	iterHandle := int(C.indigoIterateCatalysts(C.int(r.Handle)))
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

		catalysts = append(catalysts, clonedHandle)
	}

	return catalysts, nil
}
