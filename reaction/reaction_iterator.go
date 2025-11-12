// Package reaction provides reaction iteration functionality using Indigo library via CGO
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : reaction_iterator.go
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
	"runtime"
)

// ReactionIterator represents an iterator for molecules in a reaction
type ReactionIterator struct {
	handle int
	closed bool
}

// IterateReactants returns an iterator for all reactants in the reaction
func (r *Reaction) IterateReactants() (*ReactionIterator, error) {
	if r.Closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	handle := int(C.indigoIterateReactants(C.int(r.Handle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to create reactants iterator: %s", getLastError())
	}

	iter := &ReactionIterator{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(iter, (*ReactionIterator).Close)
	return iter, nil
}

// IterateProducts returns an iterator for all products in the reaction
func (r *Reaction) IterateProducts() (*ReactionIterator, error) {
	if r.Closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	handle := int(C.indigoIterateProducts(C.int(r.Handle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to create products iterator: %s", getLastError())
	}

	iter := &ReactionIterator{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(iter, (*ReactionIterator).Close)
	return iter, nil
}

// IterateCatalysts returns an iterator for all catalysts in the reaction
func (r *Reaction) IterateCatalysts() (*ReactionIterator, error) {
	if r.Closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	handle := int(C.indigoIterateCatalysts(C.int(r.Handle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to create catalysts iterator: %s", getLastError())
	}

	iter := &ReactionIterator{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(iter, (*ReactionIterator).Close)
	return iter, nil
}

// IterateMolecules returns an iterator for all molecules (reactants, products, and catalysts) in the reaction
func (r *Reaction) IterateMolecules() (*ReactionIterator, error) {
	if r.Closed {
		return nil, fmt.Errorf("reaction is closed")
	}

	handle := int(C.indigoIterateMolecules(C.int(r.Handle)))
	if handle < 0 {
		return nil, fmt.Errorf("failed to create molecules iterator: %s", getLastError())
	}

	iter := &ReactionIterator{
		handle: handle,
		closed: false,
	}

	runtime.SetFinalizer(iter, (*ReactionIterator).Close)
	return iter, nil
}

// HasNext returns true if there are more items in the iterator
func (iter *ReactionIterator) HasNext() bool {
	if iter.closed {
		return false
	}

	ret := int(C.indigoHasNext(C.int(iter.handle)))
	return ret > 0
}

// Next advances the iterator and returns the handle to the current item
func (iter *ReactionIterator) Next() (int, error) {
	if iter.closed {
		return 0, fmt.Errorf("iterator is closed")
	}

	handle := int(C.indigoNext(C.int(iter.handle)))
	if handle < 0 {
		return 0, fmt.Errorf("failed to get next item: %s", getLastError())
	}

	return handle, nil
}

// Close frees the iterator
func (iter *ReactionIterator) Close() error {
	if iter.closed || iter.handle < 0 {
		return nil
	}

	ret := int(C.indigoFree(C.int(iter.handle)))
	if ret < 0 {
		return fmt.Errorf("failed to free iterator: %s", getLastError())
	}

	iter.closed = true
	iter.handle = -1
	return nil
}
