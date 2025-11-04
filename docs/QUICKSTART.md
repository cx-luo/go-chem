# 5分钟快速开始

本指南帮助你在 5 分钟内开始使用 go-chem。

## 前置要求

- Go 1.20+
- C 编译器（gcc/clang/MSVC）

## 安装

### 1. 克隆仓库

```bash
git clone https://github.com/cx-luo/go-chem.git
cd go-chem
```

### 2. 配置环境

#### Windows (PowerShell)

```powershell
$env:CGO_ENABLED="1"
$env:CGO_CFLAGS="-I$PWD/3rd"
$env:CGO_LDFLAGS="-L$PWD/3rd/windows-x86_64"
$env:PATH="$env:PATH;$PWD/3rd/windows-x86_64"
```

#### Linux / macOS

```bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/3rd"

# Linux
export CGO_LDFLAGS="-L$(pwd)/3rd/linux-x86_64"
export LD_LIBRARY_PATH="$(pwd)/3rd/linux-x86_64:$LD_LIBRARY_PATH"

# macOS (M1/M2)
export CGO_LDFLAGS="-L$(pwd)/3rd/darwin-aarch64"
export DYLD_LIBRARY_PATH="$(pwd)/3rd/darwin-aarch64:$DYLD_LIBRARY_PATH"
```

### 3. 验证安装

```bash
go test ./test/molecule/... -v
```

如果测试通过，安装成功！

## 第一个程序

创建 `main.go`:

```go
package main

import (
    "fmt"
    "github.com/cx-luo/go-chem/molecule"
)

func main() {
    // 从 SMILES 加载分子
    mol, err := molecule.LoadMoleculeFromString("CCO")
    if err != nil {
        panic(err)
    }
    defer mol.Close()

    // 计算分子量
    mw, _ := mol.MolecularWeight()
    fmt.Printf("分子量: %.2f\n", mw)

    // 获取分子式
    formula, _ := mol.GrossFormula()
    fmt.Printf("分子式: %s\n", formula)

    // 转换为 canonical SMILES
    smiles, _ := mol.ToCanonicalSmiles()
    fmt.Printf("SMILES: %s\n", smiles)
}
```

运行：

```bash
go run main.go
```

输出：

```
分子量: 46.07
分子式: C2H6O
SMILES: CCO
```

## 常用操作

### 加载分子

```go
// 从 SMILES
mol, _ := molecule.LoadMoleculeFromString("c1ccccc1")

// 从文件
mol, _ := molecule.LoadMoleculeFromFile("molecule.mol")

// 从 InChI
molecule.InitInChI()
mol, _ := molecule.LoadFromInChI("InChI=1S/C6H6/c1-2-4-6-5-3-1/h1-6H")
```

### 保存分子

```go
// 保存为 MOL 文件
mol.SaveToFile("output.mol")

// 转换为 SMILES
smiles, _ := mol.ToSmiles()

// 转换为 JSON
json, _ := mol.ToJSON()
```

### 构建分子

```go
// 创建空分子
mol, _ := molecule.CreateMolecule()
defer mol.Close()

// 添加原子
c1, _ := mol.AddAtom("C")
c2, _ := mol.AddAtom("C")
o, _ := mol.AddAtom("O")

// 添加化学键
mol.AddBond(c1, c2, molecule.BOND_SINGLE)
mol.AddBond(c2, o, molecule.BOND_SINGLE)
```

### 渲染分子

```go
import "github.com/cx-luo/go-chem/render"

// 初始化渲染器
render.InitRenderer()
defer render.DisposeRenderer()

// 设置选项
render.SetRenderOptionInt("render-image-width", 800)
render.SetRenderOptionInt("render-image-height", 600)

// 渲染
render.RenderToFile(mol.Handle(), "molecule.png")
```

### 处理反应

```go
import "github.com/cx-luo/go-chem/reaction"

// 加载反应
rxn, _ := reaction.LoadReactionFromString("CCO>>CC=O")
defer rxn.Close()

// 自动原子映射
rxn.Automap("discard")

// 保存
rxn.SaveRxnfileToFile("reaction.rxn")
```

### 生成 InChI

```go
// 初始化
molecule.InitInChI()
defer molecule.DisposeInChI()

// 生成 InChI
inchi, _ := mol.ToInChI()
fmt.Println("InChI:", inchi)

// 生成 InChIKey
key, _ := mol.ToInChIKey()
fmt.Println("InChIKey:", key)
```

## 完整示例

### 批量处理分子

```go
package main

import (
    "fmt"
    "github.com/cx-luo/go-chem/molecule"
)

func main() {
    smilesList := []string{
        "CCO",           // 乙醇
        "c1ccccc1",      // 苯
        "CC(=O)O",       // 乙酸
        "CCN",           // 乙胺
    }

    for i, smiles := range smilesList {
        mol, err := molecule.LoadMoleculeFromString(smiles)
        if err != nil {
            fmt.Printf("Error loading %s: %v\n", smiles, err)
            continue
        }

        mw, _ := mol.MolecularWeight()
        formula, _ := mol.GrossFormula()
        atoms, _ := mol.CountAtoms()

        fmt.Printf("%d. SMILES: %-15s Formula: %-8s MW: %6.2f Atoms: %d\n",
            i+1, smiles, formula, mw, atoms)

        mol.Close()
    }
}
```

输出：

```
1. SMILES: CCO             Formula: C2H6O    MW:  46.07 Atoms: 9
2. SMILES: c1ccccc1        Formula: C6H6     MW:  78.11 Atoms: 12
3. SMILES: CC(=O)O         Formula: C2H4O2   MW:  60.05 Atoms: 8
4. SMILES: CCN             Formula: C2H7N    MW:  45.08 Atoms: 10
```

### 分子属性计算

```go
package main

import (
    "fmt"
    "github.com/cx-luo/go-chem/molecule"
)

func main() {
    // 加载阿司匹林
    mol, _ := molecule.LoadMoleculeFromString("CC(=O)Oc1ccccc1C(=O)O")
    defer mol.Close()

    fmt.Println("阿司匹林 (Aspirin) 性质：")
    fmt.Println("=========================")

    // 基本信息
    formula, _ := mol.GrossFormula()
    fmt.Printf("分子式: %s\n", formula)

    mw, _ := mol.MolecularWeight()
    fmt.Printf("分子量: %.2f\n", mw)

    atoms, _ := mol.CountAtoms()
    fmt.Printf("原子数: %d\n", atoms)

    bonds, _ := mol.CountBonds()
    fmt.Printf("键数: %d\n", bonds)

    // 药物性质
    tpsa, _ := mol.TPSA(false)
    fmt.Printf("TPSA: %.2f Ų\n", tpsa)

    rotBonds, _ := mol.NumRotatableBonds()
    fmt.Printf("可旋转键: %d\n", rotBonds)

    heavyAtoms, _ := mol.CountHeavyAtoms()
    fmt.Printf("重原子: %d\n", heavyAtoms)
}
```

### 渲染分子网格

```go
package main

import (
    "github.com/cx-luo/go-chem/molecule"
    "github.com/cx-luo/go-chem/render"
)

func main() {
    // 创建分子
    molecules := []string{
        "CCO",
        "c1ccccc1",
        "CC(=O)O",
        "CCN",
    }

    // 初始化渲染器
    render.InitRenderer()
    defer render.DisposeRenderer()

    // 创建数组
    array, _ := render.CreateArray()
    defer render.FreeObject(array)

    // 加载并添加分子
    for _, smiles := range molecules {
        mol, _ := molecule.LoadMoleculeFromString(smiles)
        render.ArrayAdd(array, mol.Handle())
        mol.Close()
    }

    // 设置选项
    render.SetRenderOptionInt("render-image-width", 1200)
    render.SetRenderOptionInt("render-image-height", 800)

    // 渲染为 2x2 网格
    render.RenderGridToFile(array, nil, 2, "molecules_grid.png")

    println("网格渲染完成: molecules_grid.png")
}
```

## 下一步

### 学习资源

1. **查看示例**: [examples/](../examples/) 目录包含更多示例
2. **阅读文档**: [docs/](../docs/) 目录包含详细文档
3. **API 参考**: [docs/API.md](API.md) 完整 API 说明
4. **环境配置**: [docs/SETUP.md](SETUP.md) 详细配置指南

### 推荐学习路径

1. **第1天**: 阅读本文档，运行基本示例
2. **第2天**: [examples/molecule/basic_usage.go](../examples/molecule/basic_usage.go)
3. **第3天**: [examples/molecule/molecule_io.go](../examples/molecule/molecule_io.go)
4. **第4天**: [examples/molecule/molecule_builder.go](../examples/molecule/molecule_builder.go)
5. **第5天**: [examples/example_render.go](../examples/example_render.go)
6. **第6天**: [examples/example_reaction.go](../examples/example_reaction.go)
7. **第7天**: 创建自己的项目！

### 常见问题

#### Q: 找不到 DLL/SO 文件

A: 确保设置了 PATH (Windows) 或 LD_LIBRARY_PATH (Linux)：

```bash
# Windows
set PATH=%PATH%;D:\path\to\go-chem\3rd\windows-x86_64

# Linux
export LD_LIBRARY_PATH=$(pwd)/3rd/linux-x86_64:$LD_LIBRARY_PATH
```

#### Q: CGO 编译错误

A: 确保安装了 C 编译器并设置 CGO_ENABLED=1：

```bash
# 检查 CGO
go env CGO_ENABLED  # 应该输出 "1"

# Windows: 安装 MinGW-w64
# Linux: sudo apt-get install build-essential
# macOS: xcode-select --install
```

#### Q: 分子加载失败

A: 检查 SMILES 格式是否正确，使用错误处理：

```go
mol, err := molecule.LoadMoleculeFromString(smiles)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
defer mol.Close()
```

## 获取帮助

- **文档**: 查看 [docs/](../docs/) 目录
- **示例**: 查看 [examples/](../examples/) 目录
- **Issue**: 在 GitHub 创建 issue
- **邮件**: <chengxiang.luo@foxmail.com>

---

🎉 恭喜！你已经完成了快速开始教程！

现在你可以：

- ✅ 加载和保存分子
- ✅ 计算分子属性
- ✅ 渲染分子结构
- ✅ 处理化学反应
- ✅ 生成 InChI

继续探索 [examples/](../examples/) 目录了解更多功能！
