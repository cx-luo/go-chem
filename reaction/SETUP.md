# Setup Guide for Reaction Package

This guide explains how to set up your environment to use the reaction package.

## Prerequisites

- Go 1.16 or later
- GCC or MinGW (for CGO compilation)
- Indigo library files in the `3rd` directory

## Windows Setup

### 1. Ensure DLL files are accessible

The reaction package requires several DLL files from the Indigo library:

- `indigo.dll`
- `msvcp140.dll`
- `vcruntime140.dll`
- `vcruntime140_1.dll`

These files are located in the `3rd` directory. You have two options:

**Option A: Add to System PATH (Recommended)**

```cmd
set PATH=%PATH%;D:\for_github\go-chem\3rd
```

Or permanently add the `3rd` directory to your system PATH through System Properties.

**Option B: Copy DLLs to executable directory**

```cmd
copy 3rd\*.dll <your_project_directory>\
```

### 2. Set CGO environment variables

```cmd
set CGO_ENABLED=1
set CGO_CFLAGS=-ID:/for_github/go-chem/3rd
set CGO_LDFLAGS=-LD:/for_github/go-chem/3rd
```

### 3. Build and run

```cmd
cd reaction
go build
```

## Linux Setup

### 1. Ensure shared libraries are accessible

```bash
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/path/to/go-chem/3rd
```

Or add to `~/.bashrc`:

```bash
echo 'export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/path/to/go-chem/3rd' >> ~/.bashrc
source ~/.bashrc
```

### 2. Set CGO environment variables

```bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I/path/to/go-chem/3rd"
export CGO_LDFLAGS="-L/path/to/go-chem/3rd -lindigo"
```

### 3. Build and run

```bash
cd reaction
go build
```

## Running Tests

### From project root

```bash
# Set PATH first (Windows)
set PATH=%PATH%;D:\for_github\go-chem\3rd

# Run tests
cd test/reaction
go test -v
```

### Run specific test

```bash
go test -v -run TestLoadReactionFromString
```

### Run with coverage

```bash
go test -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Running Examples

```bash
# Set PATH first (Windows)
set PATH=%PATH%;D:\for_github\go-chem\3rd

# Run example
cd reaction
go run example_reaction.go
```

## Troubleshooting

### Error: `exit status 0xc0000135` (Windows)

This error means a required DLL cannot be found. Solutions:

1. Add the `3rd` directory to your PATH
2. Copy all DLL files to the same directory as your executable
3. Use a tool like Dependency Walker to identify missing DLLs

### Error: `cannot find -lindigo` (Linux)

This means the linker cannot find the Indigo library. Solutions:

1. Verify `libindigo.so` exists in the `3rd` directory
2. Set `LD_LIBRARY_PATH` correctly
3. Check CGO_LDFLAGS includes the correct path

### Error: `undefined reference to 'indigoCreateReaction'`

This indicates a linking problem. Solutions:

1. Ensure CGO_LDFLAGS includes `-lindigo`
2. Verify the Indigo library is the correct architecture (32-bit vs 64-bit)
3. Check that the library file is not corrupted

### Error: `runtime error: cgo argument has Go pointer to Go pointer`

This is a Go CGO safety check. The package handles this correctly, but if you see this error:

1. Update to the latest Go version
2. Check that you're not passing Go pointers to C incorrectly
3. Review the CGO documentation

## Development Setup

### VS Code

Add to `.vscode/settings.json`:

```json
{
    "go.toolsEnvVars": {
        "CGO_ENABLED": "1",
        "CGO_CFLAGS": "-ID:/for_github/go-chem/3rd",
        "CGO_LDFLAGS": "-LD:/for_github/go-chem/3rd"
    }
}
```

### GoLand

1. Go to Settings → Go → Build Tags & Vendoring
2. Add custom environment variables:
   - `CGO_ENABLED=1`
   - `CGO_CFLAGS=-ID:/for_github/go-chem/3rd`
   - `CGO_LDFLAGS=-LD:/for_github/go-chem/3rd`

## Testing Indigo Installation

Create a simple test file to verify Indigo is working:

```go
package main

import (
    "fmt"
    "github.com/cx-luo/go-chem/reaction"
)

func main() {
    r, err := reaction.CreateReaction()
    if err != nil {
        panic(err)
    }
    defer r.Close()
    fmt.Println("Indigo is working correctly!")
}
```

Run with:

```bash
go run test_indigo.go
```

If this works, your environment is set up correctly.

## Additional Resources

- [Indigo Documentation](https://lifescience.opensource.epam.com/indigo/)
- [CGO Documentation](https://golang.org/cmd/cgo/)
- [Go Testing](https://golang.org/pkg/testing/)
