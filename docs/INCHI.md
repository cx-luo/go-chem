# InChI 功能完整指南

InChI (IUPAC International Chemical Identifier) 是一个用于表示化学物质结构的国际标准标识符。

## 目录

- [什么是 InChI](#什么是-inchi)
- [快速开始](#快速开始)
- [API 参考](#api-参考)
- [使用示例](#使用示例)
- [InChI 格式详解](#inchi-格式详解)
- [高级用法](#高级用法)
- [常见问题](#常见问题)

## 什么是 InChI

### InChI 简介

InChI 是一个文本字符串，用于唯一标识化学物质。它由 IUPAC 开发，具有以下特点：

- **唯一性**: 同一化学结构总是生成相同的 InChI
- **层次化**: 包含多个层，逐步添加结构信息
- **标准化**: 遵循国际标准，不同软件生成的结果一致
- **可读性**: 人类可以理解基本结构信息

### InChIKey

InChIKey 是 InChI 的定长哈希表示（27字符），格式为：

```
XXXXXXXXXXXXXX-YYYYYYYYYY-Z
```

- 第一部分（14字符）: 主哈希（连接层）
- 第二部分（10字符）: 立体化学和同位素信息
- 第三部分（1字符）: InChI 版本和选项

### InChI vs SMILES

| 特性 | InChI | SMILES |
|------|-------|--------|
| 标准化 | 是 | 否（多种变体） |
| 唯一性 | 保证唯一 | 同一分子可有多个表示 |
| 可读性 | 较低 | 较高 |
| 数据库检索 | 优秀（InChIKey） | 需要规范化 |
| 立体化学 | 完整支持 | 支持 |

## 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/cx-luo/go-chem/molecule"
)

func main() {
    // 1. 初始化 InChI 模块
    err := molecule.InitInChI()
    if err != nil {
        panic(err)
    }
    defer molecule.DisposeInChI()

    // 2. 加载分子
    mol, err := molecule.LoadMoleculeFromString("CCO")
    if err != nil {
        panic(err)
    }
    defer mol.Close()

    // 3. 生成 InChI
    inchi, err := mol.ToInChI()
    if err != nil {
        panic(err)
    }
    fmt.Println("InChI:", inchi)
    // 输出: InChI=1S/C2H6O/c1-2-3/h3H,2H2,1H3

    // 4. 生成 InChIKey
    key, err := mol.ToInChIKey()
    if err != nil {
        panic(err)
    }
    fmt.Println("InChIKey:", key)
    // 输出: LFQSCWFLJHTTHZ-UHFFFAOYSA-N
}
```

## API 参考

### 初始化和释放

#### InitInChI

```go
func InitInChI() error
```

初始化 InChI 模块。必须在使用 InChI 功能前调用。

**返回值:**

- `error`: 错误信息

#### DisposeInChI

```go
func DisposeInChI() error
```

释放 InChI 模块资源。应在程序结束时调用。

**返回值:**

- `error`: 错误信息

#### InChIVersion

```go
func InChIVersion() string
```

返回 InChI 库版本。

**返回值:**

- `string`: 版本字符串

### 生成 InChI

#### ToInChI

```go
func (m *Molecule) ToInChI() (string, error)
```

从分子生成标准 InChI 字符串。

**返回值:**

- `string`: InChI 字符串
- `error`: 错误信息

#### ToInChIKey

```go
func (m *Molecule) ToInChIKey() (string, error)
```

从分子生成 InChIKey。

**返回值:**

- `string`: InChIKey（27字符）
- `error`: 错误信息

#### InChIToKey

```go
func InChIToKey(inchi string) (string, error)
```

从 InChI 字符串直接生成 InChIKey。

**参数:**

- `inchi` (string): InChI 字符串

**返回值:**

- `string`: InChIKey
- `error`: 错误信息

### 加载分子

#### LoadFromInChI

```go
func LoadFromInChI(inchi string) (*Molecule, error)
```

从 InChI 字符串加载分子。

**参数:**

- `inchi` (string): InChI 字符串

**返回值:**

- `*Molecule`: 分子对象
- `error`: 错误信息

### 详细信息

#### ToInChIWithInfo

```go
func (m *Molecule) ToInChIWithInfo() (*InChIResult, error)
```

生成 InChI 并返回详细信息。

**返回值:**

- `*InChIResult`: 包含 InChI、InChIKey、警告、日志等
- `error`: 错误信息

#### InChIResult 结构

```go
type InChIResult struct {
    InChI   string  // InChI 字符串
    Key     string  // InChIKey
    Warning string  // 警告信息
    Log     string  // 日志信息
    AuxInfo string  // 辅助信息
}
```

### 辅助函数

#### InChIWarning

```go
func InChIWarning() string
```

获取最后一次 InChI 生成的警告信息。

#### InChILog

```go
func InChILog() string
```

获取最后一次 InChI 生成的日志信息。

#### InChIAuxInfo

```go
func InChIAuxInfo() string
```

获取最后一次 InChI 生成的辅助信息。

#### ResetInChIOptions

```go
func ResetInChIOptions() error
```

重置 InChI 选项到默认值。

## 使用示例

### 示例 1: 基本 InChI 生成

```go
func Example1() {
    molecule.InitInChI()
    defer molecule.DisposeInChI()

    // 苯
    mol, _ := molecule.LoadMoleculeFromString("c1ccccc1")
    defer mol.Close()

    inchi, _ := mol.ToInChI()
    key, _ := mol.ToInChIKey()

    fmt.Println("分子: 苯")
    fmt.Println("InChI:", inchi)
    fmt.Println("InChIKey:", key)
}
```

**输出:**

```
分子: 苯
InChI: InChI=1S/C6H6/c1-2-4-6-5-3-1/h1-6H
InChIKey: UHOVQNZJYSORNB-UHFFFAOYSA-N
```

### 示例 2: 批量生成 InChIKey

```go
func Example2() {
    molecule.InitInChI()
    defer molecule.DisposeInChI()

    molecules := map[string]string{
        "甲醇":   "CO",
        "乙醇":   "CCO",
        "丙醇":   "CCCO",
        "乙酸":   "CC(=O)O",
    }

    for name, smiles := range molecules {
        mol, _ := molecule.LoadMoleculeFromString(smiles)
        key, _ := mol.ToInChIKey()
        fmt.Printf("%-10s %s\n", name, key)
        mol.Close()
    }
}
```

### 示例 3: InChI 详细信息

```go
func Example3() {
    molecule.InitInChI()
    defer molecule.DisposeInChI()

    mol, _ := molecule.LoadMoleculeFromString("CC(=O)O")
    defer mol.Close()

    // 获取详细信息
    result, _ := mol.ToInChIWithInfo()

    fmt.Println("InChI:", result.InChI)
    fmt.Println("InChIKey:", result.Key)

    if result.Warning != "" {
        fmt.Println("警告:", result.Warning)
    }

    if result.Log != "" {
        fmt.Println("日志:", result.Log)
    }

    if result.AuxInfo != "" {
        fmt.Println("辅助信息:", result.AuxInfo)
    }
}
```

### 示例 4: InChI 往返转换

```go
func Example4() {
    molecule.InitInChI()
    defer molecule.DisposeInChI()

    // 原始分子
    mol1, _ := molecule.LoadMoleculeFromString("c1ccccc1")
    defer mol1.Close()

    // 生成 InChI
    inchi, _ := mol1.ToInChI()
    fmt.Println("InChI:", inchi)

    // 从 InChI 重新加载
    mol2, _ := molecule.LoadFromInChI(inchi)
    defer mol2.Close()

    // 验证结构
    atoms1, _ := mol1.CountAtoms()
    atoms2, _ := mol2.CountAtoms()
    fmt.Printf("原子数匹配: %v\n", atoms1 == atoms2)
}
```

### 示例 5: 分子比较（通过 InChIKey）

```go
func Example5() {
    molecule.InitInChI()
    defer molecule.DisposeInChI()

    // 不同 SMILES 表示的同一分子
    mol1, _ := molecule.LoadMoleculeFromString("CCO")
    mol2, _ := molecule.LoadMoleculeFromString("OCC")
    defer mol1.Close()
    defer mol2.Close()

    key1, _ := mol1.ToInChIKey()
    key2, _ := mol2.ToInChIKey()

    if key1 == key2 {
        fmt.Println("这是同一个分子！")
        fmt.Println("InChIKey:", key1)
    }
}
```

### 示例 6: 数据库存储

```go
type MoleculeRecord struct {
    ID        int
    Name      string
    SMILES    string
    InChI     string
    InChIKey  string
}

func Example6() {
    molecule.InitInChI()
    defer molecule.DisposeInChI()

    records := []MoleculeRecord{}

    smiles_list := []string{"CCO", "c1ccccc1", "CC(=O)O"}

    for i, smiles := range smiles_list {
        mol, _ := molecule.LoadMoleculeFromString(smiles)

        inchi, _ := mol.ToInChI()
        key, _ := mol.ToInChIKey()

        record := MoleculeRecord{
            ID:       i + 1,
            SMILES:   smiles,
            InChI:    inchi,
            InChIKey: key,
        }

        records = append(records, record)
        mol.Close()
    }

    // 打印记录
    for _, rec := range records {
        fmt.Printf("ID: %d\n", rec.ID)
        fmt.Printf("  SMILES: %s\n", rec.SMILES)
        fmt.Printf("  InChIKey: %s\n", rec.InChIKey)
    }
}
```

## InChI 格式详解

### InChI 层次结构

标准 InChI 由多个层组成：

```
InChI=1S/C2H6O/c1-2-3/h3H,2H2,1H3
│      │ │    │ │        │
│      │ │    │ │        └─ H 层（氢原子）
│      │ │    │ └────────── C 层（连接性）
│      │ │    └──────────── F 层（化学式）
│      │ └───────────────── 版本
│      └─────────────────── 标准 InChI 标记
└────────────────────────── InChI 前缀
```

### 各层说明

1. **版本层**: `1S` 表示标准 InChI 版本 1
2. **化学式层 (F)**: `C2H6O` - Hill 系统排序
3. **连接层 (C)**: `c1-2-3` - 原子连接关系
4. **氢原子层 (H)**: `h3H,2H2,1H3` - 氢原子分布

### InChIKey 格式

```
LFQSCWFLJHTTHZ-UHFFFAOYSA-N
│             │ │         │└─ 质子化状态
│             │ │         └── 立体化学
│             │ └──────────── 立体和同位素层
│             └────────────── 连接层哈希
└──────────────────────────── 连接层哈希
```

## 高级用法

### 错误处理

```go
func SafeInChIGeneration(smiles string) (string, string, error) {
    if err := molecule.InitInChI(); err != nil {
        return "", "", fmt.Errorf("InChI 初始化失败: %w", err)
    }
    defer molecule.DisposeInChI()

    mol, err := molecule.LoadMoleculeFromString(smiles)
    if err != nil {
        return "", "", fmt.Errorf("SMILES 解析失败: %w", err)
    }
    defer mol.Close()

    inchi, err := mol.ToInChI()
    if err != nil {
        return "", "", fmt.Errorf("InChI 生成失败: %w", err)
    }

    key, err := mol.ToInChIKey()
    if err != nil {
        return inchi, "", fmt.Errorf("InChIKey 生成失败: %w", err)
    }

    return inchi, key, nil
}
```

### 并发处理

```go
func ConcurrentInChI(smilesList []string) map[string]string {
    molecule.InitInChI()
    defer molecule.DisposeInChI()

    results := make(map[string]string)
    mutex := &sync.Mutex{}
    wg := &sync.WaitGroup{}

    for _, smiles := range smilesList {
        wg.Add(1)
        go func(s string) {
            defer wg.Done()

            mol, err := molecule.LoadMoleculeFromString(s)
            if err != nil {
                return
            }
            defer mol.Close()

            key, err := mol.ToInChIKey()
            if err != nil {
                return
            }

            mutex.Lock()
            results[s] = key
            mutex.Unlock()
        }(smiles)
    }

    wg.Wait()
    return results
}
```

## 常见问题

### Q: 什么时候需要调用 InitInChI?

A: 在任何 InChI 功能使用前调用一次。库会自动检查初始化状态。

### Q: InChI 生成失败怎么办?

A: 检查分子结构是否有效，查看警告信息：

```go
result, err := mol.ToInChIWithInfo()
if err != nil {
    fmt.Println("错误:", err)
}
if result.Warning != "" {
    fmt.Println("警告:", result.Warning)
}
```

### Q: 如何处理立体化学?

A: Indigo InChI 自动处理立体化学信息，确保分子加载时包含立体信息。

### Q: InChIKey 可以反推回分子吗?

A: 不可以。InChIKey 是单向哈希，但 InChI 字符串可以转回分子。

### Q: 不同软件生成的 InChI 一致吗?

A: 标准 InChI 应该一致，但某些边缘情况可能有差异。

## 性能考虑

- InChI 生成比 SMILES 慢（需要标准化）
- InChIKey 生成很快（基于 InChI 的哈希）
- 建议批量处理时使用并发

## 相关资源

- [InChI Trust 官网](https://www.inchi-trust.org/)
- [InChI FAQ](https://www.inchi-trust.org/inchi-faq/)
- [Indigo InChI 文档](https://lifescience.opensource.epam.com/indigo/api/)

---

💡 **提示**: InChIKey 是数据库检索的理想标识符！
