# Go-Chem InChI CGO 实现

本项目只保留基于官方 InChI 动态库的 CGO 实现。原 pure Go InChI 生成器已移除，避免继续暴露不完整或不兼容的 InChI 结果。

## 快速开始

```go
import "github.com/cx-luo/go-chem/molecule"

loader := molecule.SmilesLoader{}
mol, err := loader.Parse("CCO")
if err != nil {
    log.Fatal(err)
}

generator := molecule.NewInChIGeneratorCGO()
result, err := generator.GenerateInChI(mol)
if err != nil {
    log.Fatal(err)
}

fmt.Println("InChI:", result.InChI)
fmt.Println("InChIKey:", result.InChIKey)
```

## 运行示例

```bash
# Linux
export LD_LIBRARY_PATH=$(pwd)/3rd:$LD_LIBRARY_PATH
CGO_ENABLED=1 go run examples/inchi_cgo_example.go

# Windows
set PATH=%PATH%;%CD%\3rd
go run examples/inchi_cgo_example.go
```

## API

- `molecule.NewInChIGeneratorCGO()`: 创建 CGO InChI 生成器。
- `(*InChIGeneratorCGO).GenerateInChI(mol)`: 从 `Molecule` 生成 InChI 和 InChIKey。
- `(*InChIGeneratorCGO).SetOptions(options)`: 设置官方 InChI 参数，例如 `FixedH`、`RecMet`。
- `molecule.GenerateInChIKeyCGO(inchi)`: 通过官方库从 InChI 生成 InChIKey。
- `molecule.GetInChIVersion()`: 获取当前链接的 InChI 库版本。

## 依赖

InChI 动态库文件放在 `3rd/` 目录：

- Windows: `3rd/libinchi.dll`
- Linux: `3rd/libinchi.so`
- 头文件: `3rd/inchi_api.h`

更多构建和部署说明见 `docs/INCHI_CGO_GUIDE.md`。
