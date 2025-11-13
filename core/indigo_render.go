// Package core coding=utf-8
// @Project : go-chem
// @Time    : 2025/11/13 10:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : indigo_render.go
// @Software: GoLand
package core

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

	"github.com/cx-luo/go-chem/render"
)

// DefaultRenderOptions returns default rendering options
func defaultRenderOptions() *render.RenderOptions {
	return &render.RenderOptions{
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

// InitRenderer initializes the Indigo renderer with default options for the current session
// This should be called before using rendering functions
func (in *Indigo) InitRenderer() (*render.Renderer, error) {
	ret := int(C.indigoRendererInit(C.ulonglong(in.sid)))
	if ret < 0 {
		return nil, fmt.Errorf("failed to initialize renderer: %s", lastErrorString())
	}

	return &render.Renderer{Sid: in.sid, Options: defaultRenderOptions(), RendererInitialized: true}, nil
}
