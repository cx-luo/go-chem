# Molecule 包

分子处理包提供了完整的化学分子操作功能，基于 Indigo 库通过 CGO 封装。

## 📋 目录

- [功能特性](#功能特性)
- [快速开始](#快速开始)
- [API 文档](#api-文档)
- [使用示例](#使用示例)
- [元素支持](#元素支持)

## 功能特性

### 核心功能

- ✅ **分子创建和管理**
  - 创建空分子
  - 从头构建分子
  - 分子克隆
  - 分子关闭和资源管理

- ✅ **分子加载**
  - SMILES 字符串
  - MOL 文件
  - SDF 文件
  - Query Molecule
  - SMARTS 模式
  - InChI 字符串

- ✅ **分子保存**
  - SMILES / Canonical SMILES
  - MOL 文件
  - JSON 格式
  - Base64 编码
  - SMARTS

- ✅ **分子构建**
  - 添加原子（元素符号）
  - 添加化学键（单键、双键、三键、芳香键）
  - 设置原子电荷
  - 设置同位素
  - 设置自由基
  - 添加 R-site
  - 合并分子

- ✅ **分子属性计算**
  - 分子量
  - 单同位素质量
  - 最丰富质量
  - 质量组成
  - 总分子式
  - 分子式
  - TPSA（拓扑极性表面积）
  - 可旋转键数量
  - 重原子计数

- ✅ **分子操作**
  - 芳香化 / 去芳香化
  - 氢原子折叠 / 展开
  - 2D 布局
  - 2D 清理
  - 标准化
  - 归一化
  - 离子化（指定 pH）

- ✅ **结构分析**
  - 原子计数
  - 化学键计数
  - 重原子计数
  - 连通组件计数
  - SSSR（最小环集）计数

- ✅ **InChI 支持**
  - InChI 生成
  - InChIKey 生成
  - 从 InChI 加载
  - 警告和日志信息
  - 辅助信息

## 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/cx-luo/go-indigo/molecule"
)

func main() {
    // 从 SMILES 加载分子
    mol, err := molecule.LoadMoleculeFromString("CCO")
    if err != nil {
        panic(err)
    }
    defer mol.Close()

    // 获取基本信息
    atomCount, _ := mol.CountAtoms()
    bondCount, _ := mol.CountBonds()
    fmt.Printf("原子数: %d, 键数: %d\n", atomCount, bondCount)

    // 计算分子量
    mw, _ := mol.MolecularWeight()
    fmt.Printf("分子量: %.2f\n", mw)

    // 转换为 canonical SMILES
    smiles, _ := mol.ToCanonicalSmiles()
    fmt.Printf("Canonical SMILES: %s\n", smiles)
}
```

## API 文档

### 分子创建

```go
// 创建空分子
mol, err := molecule.CreateMolecule()

// 创建查询分子
queryMol, err := molecule.CreateQueryMolecule()
```

### 分子加载

```go
// 从 SMILES 加载
mol, err := molecule.LoadMoleculeFromString("c1ccccc1")

// 从文件加载
mol, err := molecule.LoadMoleculeFromFile("molecule.mol")

// 从缓冲区加载
mol, err := molecule.LoadMoleculeFromBuffer(data)

// 从 InChI 加载
mol, err := molecule.LoadFromInChI("InChI=1S/C6H6/c1-2-4-6-5-3-1/h1-6H")

// 加载 Query Molecule
queryMol, err := molecule.LoadQueryMoleculeFromString("[#6]CO")

// 加载 SMARTS
smarts, err := molecule.LoadSmartsFromString("[OH]")

// 通用结构加载（带参数）
mol, err := molecule.LoadStructureFromString(data, "")
```

### 分子保存

```go
// 转换为 SMILES
smiles, err := mol.ToSmiles()

// 转换为 Canonical SMILES
canonical, err := mol.ToCanonicalSmiles()

// 转换为 SMARTS
smarts, err := mol.ToSmarts()

// 转换为 MOL 格式字符串
molfile, err := mol.ToMolfile()

// 保存到文件
err := mol.SaveToFile("output.mol")

// 转换为 JSON
json, err := mol.ToJSON()

// 保存为 JSON 文件
err := mol.SaveToJSONFile("output.json")

// 转换为 Base64
base64, err := mol.ToBase64String()
```

### 分子构建

```go
// 创建空分子
mol, _ := molecule.CreateMolecule()
defer mol.Close()

// 添加原子
c1, _ := mol.AddAtom("C")
c2, _ := mol.AddAtom("C")
o, _ := mol.AddAtom("O")

// 添加化学键
bond1, _ := mol.AddBond(c1, c2, molecule.BOND_SINGLE)
bond2, _ := mol.AddBond(c2, o, molecule.BOND_SINGLE)

// 设置原子属性
molecule.SetCharge(o, -1)     // 设置电荷
molecule.SetIsotope(c1, 13)   // 设置同位素
molecule.SetRadical(c2, 2)    // 设置自由基

// 重置原子
molecule.ResetAtom(c1, "N")   // 将碳改为氮

// 添加 R-site
rsite, _ := mol.AddRSite("R1")

// 合并两个分子
mol1.Merge(mol2)
```

### 分子属性

```go
// 基本计数
atomCount, _ := mol.CountAtoms()
bondCount, _ := mol.CountBonds()
heavyAtoms, _ := mol.CountHeavyAtoms()

// 质量相关
mw, _ := mol.MolecularWeight()           // 分子量
monoMass, _ := mol.MonoisotopicMass()    // 单同位素质量
abundantMass, _ := mol.MostAbundantMass() // 最丰富质量
massComp, _ := mol.MassComposition()     // 质量组成

// 分子式
grossFormula, _ := mol.GrossFormula()       // 总分子式
molFormula, _ := mol.MolecularFormula()     // 分子式

// 药物相关属性
tpsa, _ := mol.TPSA(false)               // TPSA
rotatableBonds, _ := mol.NumRotatableBonds() // 可旋转键数

// 结构分析
components, _ := mol.CountComponents()    // 连通组件数
rings, _ := mol.CountSSSR()               // 最小环集数
```

### 分子操作

```go
// 芳香化处理
mol.Aromatize()
mol.Dearomatize()

// 氢原子处理
mol.FoldHydrogens()
mol.UnfoldHydrogens()

// 2D 操作
mol.Layout()    // 2D 布局
mol.Clean2D()   // 2D 清理

// 标准化
mol.Normalize("")        // 归一化
mol.Standardize()        // 标准化
mol.Ionize(7.0, 0.5)    // 在 pH 7.0 离子化

// 克隆
newMol, _ := mol.Clone()
defer newMol.Close()
```

### 分子属性管理

```go
// 名称
mol.SetName("Aspirin")
name, _ := mol.Name()

// 自定义属性
mol.SetProperty("CAS", "50-78-2")
has, _ := mol.HasProperty("CAS")
value, _ := mol.GetProperty("CAS")
mol.RemoveProperty("CAS")
```

### InChI 功能

```go
// 初始化 InChI 模块
molecule.InitInChI()
defer molecule.DisposeInChI()

// 生成 InChI
inchi, err := mol.ToInChI()

// 生成 InChIKey
key, err := mol.ToInChIKey()

// 或者从 InChI 字符串生成 Key
key, err := molecule.InChIToKey(inchi)

// 获取详细信息
result, err := mol.ToInChIWithInfo()
fmt.Println("InChI:", result.InChI)
fmt.Println("Key:", result.Key)
fmt.Println("Warning:", result.Warning)
fmt.Println("Log:", result.Log)
fmt.Println("AuxInfo:", result.AuxInfo)

// 从 InChI 加载分子
mol, err := molecule.LoadFromInChI("InChI=1S/C2H6O/c1-2-3/h3H,2H2,1H3")

// 辅助函数
warning := molecule.InChIWarning()
log := molecule.InChILog()
auxInfo := molecule.InChIAuxInfo()

// 重置选项
molecule.ResetInChIOptions()

// 获取版本
version := molecule.InChIVersion()
```

## 使用示例

### 示例 1: 从 SMILES 加载并分析

```go
func Example1() {
    // 加载阿司匹林
    mol, _ := molecule.LoadMoleculeFromString("CC(=O)Oc1ccccc1C(=O)O")
    defer mol.Close()

    // 计算属性
    mw, _ := mol.MolecularWeight()
    formula, _ := mol.GrossFormula()
    tpsa, _ := mol.TPSA(false)

    fmt.Printf("分子量: %.2f\n", mw)
    fmt.Printf("分子式: %s\n", formula)
    fmt.Printf("TPSA: %.2f\n", tpsa)
}
```

### 示例 2: 构建甲烷分子

```go
func Example2() {
    mol, _ := molecule.CreateMolecule()
    defer mol.Close()

    // 添加碳原子
    c, _ := mol.AddAtom("C")

    // 添加 4 个氢原子
    for i := 0; i < 4; i++ {
        h, _ := mol.AddAtom("H")
        mol.AddBond(c, h, molecule.BOND_SINGLE)
    }

    // 转换为 SMILES
    smiles, _ := mol.ToSmiles()
    fmt.Println("SMILES:", smiles)  // 输出: C
}
```

### 示例 3: 分子标准化流程

```go
func Example3() {
    mol, _ := molecule.LoadMoleculeFromString("[Na+].CC(=O)[O-]")
    defer mol.Close()

    // 标准化流程
    mol.Normalize("")
    mol.Standardize()
    mol.Aromatize()
    mol.FoldHydrogens()

    // 输出标准化后的 SMILES
    smiles, _ := mol.ToCanonicalSmiles()
    fmt.Println("标准 SMILES:", smiles)
}
```

### 示例 4: 批量处理分子

```go
func Example4() {
    smilesList := []string{
        "CCO",
        "c1ccccc1",
        "CC(=O)O",
        "CCN",
    }

    for i, smiles := range smilesList {
        mol, _ := molecule.LoadMoleculeFromString(smiles)

        mw, _ := mol.MolecularWeight()
        atoms, _ := mol.CountAtoms()

        fmt.Printf("%d. SMILES: %s, MW: %.2f, Atoms: %d\n",
            i+1, smiles, mw, atoms)

        mol.Close()
    }
}
```

### 示例 5: InChI 转换

```go
func Example5() {
    // 初始化
    molecule.InitInChI()
    defer molecule.DisposeInChI()

    mol, _ := molecule.LoadMoleculeFromString("c1ccccc1")
    defer mol.Close()

    // 生成 InChI 和 InChIKey
    result, _ := mol.ToInChIWithInfo()

    fmt.Println("InChI:", result.InChI)
    fmt.Println("InChIKey:", result.Key)

    if result.Warning != "" {
        fmt.Println("警告:", result.Warning)
    }
}
```

### 示例 6: 文件格式转换

```go
func Example6() {
    // SMILES -> MOL 文件
    mol, _ := molecule.LoadMoleculeFromString("CCO")
    defer mol.Close()

    mol.SaveToFile("ethanol.mol")

    // MOL 文件 -> SMILES
    mol2, _ := molecule.LoadMoleculeFromFile("ethanol.mol")
    defer mol2.Close()

    smiles, _ := mol2.ToCanonicalSmiles()
    fmt.Println("SMILES:", smiles)
}
```

## 元素支持

package 提供完整的元素周期表支持（通过 `elements.go`）：

### 元素常量

```go
const (
    ELEM_H  = 1   // 氢
    ELEM_C  = 6   // 碳
    ELEM_N  = 7   // 氮
    ELEM_O  = 8   // 氧
    ELEM_F  = 9   // 氟
    ELEM_P  = 15  // 磷
    ELEM_S  = 16  // 硫
    ELEM_Cl = 17  // 氯
    ELEM_Br = 35  // 溴
    ELEM_I  = 53  // 碘
    // ... 更多元素
)
```

### 特殊元素

```go
const (
    ELEM_PSEUDO   = -1  // 伪原子
    ELEM_RSITE    = -2  // R-site
    ELEM_TEMPLATE = -3  // 模板原子
)

const (
    RADICAL_SINGLET = 2
    RADICAL_DOUBLET = 3
    RADICAL_TRIPLET = 4
)
```

### 化学键类型

```go
const (
    BOND_SINGLE   = 1  // 单键
    BOND_DOUBLE   = 2  // 双键
    BOND_TRIPLE   = 3  // 三键
    BOND_AROMATIC = 4  // 芳香键
)
```

## 性能考虑

1. **资源管理**: 始终使用 `defer mol.Close()` 确保资源释放
2. **批量操作**: 对于大量分子，考虑使用 goroutine 并行处理
3. **内存使用**: Clone 操作会复制整个分子，注意内存使用
4. **CGO 开销**: 频繁的小操作可能有 CGO 调用开销

## 错误处理

所有函数都返回 error，务必检查错误：

```go
mol, err := molecule.LoadMoleculeFromString("CCO")
if err != nil {
    log.Fatalf("加载分子失败: %v", err)
}
defer mol.Close()
```

## 相关文档

- [Reaction 包文档](../reaction/README.md)
- [Render 包文档](../render/README.md)
- [环境设置](../reaction/SETUP.md)
- [示例代码](../examples/molecule/)

## 许可证

本包基于 Apache License 2.0 许可证。
