# API 参考文档

go-indigo 完整 API 参考文档。

## 目录

- [Molecule 包](#molecule-包)
- [Reaction 包](#reaction-包)
- [Render 包](#render-包)
- [类型定义](#类型定义)
- [常量](#常量)
- [错误处理](#错误处理)

## Molecule 包

### 类型

#### Molecule

```go
type Molecule struct {
    handle int    // Indigo 句柄
    closed bool   // 是否已关闭
}
```

分子对象，表示一个化学分子结构。

### 创建和加载

#### CreateMolecule

```go
func CreateMolecule() (*Molecule, error)
```

创建一个空的分子对象。

**返回值:**

- `*Molecule`: 新创建的分子对象
- `error`: 错误信息

**示例:**

```go
mol, err := molecule.CreateMolecule()
if err != nil {
    log.Fatal(err)
}
defer mol.Close()
```

#### LoadMoleculeFromString

```go
func LoadMoleculeFromString(data string) (*Molecule, error)
```

从字符串（SMILES 或其他格式）加载分子。

**参数:**

- `data` (string): SMILES 字符串或其他格式的分子数据

**返回值:**

- `*Molecule`: 加载的分子对象
- `error`: 错误信息

**示例:**

```go
mol, err := molecule.LoadMoleculeFromString("CCO")
if err != nil {
    log.Fatal(err)
}
defer mol.Close()
```

#### LoadMoleculeFromFile

```go
func LoadMoleculeFromFile(filename string) (*Molecule, error)
```

从文件加载分子。

**参数:**

- `filename` (string): 文件路径

**返回值:**

- `*Molecule`: 加载的分子对象
- `error`: 错误信息

**支持的格式:**

- MOL 文件 (.mol)
- SDF 文件 (.sdf)
- SMILES 文件 (.smi)

#### LoadMoleculeFromBuffer

```go
func LoadMoleculeFromBuffer(buffer []byte) (*Molecule, error)
```

从字节缓冲区加载分子。

**参数:**

- `buffer` ([]byte): 包含分子数据的字节数组

### 分子操作

#### Close

```go
func (m *Molecule) Close() error
```

关闭分子对象，释放相关资源。

**返回值:**

- `error`: 错误信息

**注意:** 应该使用 `defer mol.Close()` 确保资源释放。

#### Clone

```go
func (m *Molecule) Clone() (*Molecule, error)
```

创建分子的深拷贝。

**返回值:**

- `*Molecule`: 克隆的分子对象
- `error`: 错误信息

#### CountAtoms

```go
func (m *Molecule) CountAtoms() (int, error)
```

返回分子中的原子数量。

**返回值:**

- `int`: 原子数量
- `error`: 错误信息

#### CountBonds

```go
func (m *Molecule) CountBonds() (int, error)
```

返回分子中的化学键数量。

**返回值:**

- `int`: 化学键数量
- `error`: 错误信息

#### CountHeavyAtoms

```go
func (m *Molecule) CountHeavyAtoms() (int, error)
```

返回重原子（非氢原子）数量。

**返回值:**

- `int`: 重原子数量
- `error`: 错误信息

### 分子转换

#### Aromatize

```go
func (m *Molecule) Aromatize() error
```

将分子芳香化（识别并标记芳香环）。

#### Dearomatize

```go
func (m *Molecule) Dearomatize() error
```

将分子去芳香化（将芳香键转换为单双键）。

#### FoldHydrogens

```go
func (m *Molecule) FoldHydrogens() error
```

折叠氢原子（隐式表示氢原子）。

#### UnfoldHydrogens

```go
func (m *Molecule) UnfoldHydrogens() error
```

展开氢原子（显式添加氢原子）。

#### Layout

```go
func (m *Molecule) Layout() error
```

执行 2D 布局（计算 2D 坐标）。

#### Clean2D

```go
func (m *Molecule) Clean2d() error
```

清理 2D 结构（优化 2D 坐标）。

#### Normalize

```go
func (m *Molecule) Normalize(options string) error
```

归一化分子结构。

**参数:**

- `options` (string): 归一化选项

#### Standardize

```go
func (m *Molecule) Standardize() error
```

标准化分子结构。

#### Ionize

```go
func (m *Molecule) Ionize(pH float32, pHTolerance float32) error
```

在指定 pH 下离子化分子。

**参数:**

- `pH` (float32): pH 值
- `pHTolerance` (float32): pH 容差

### 分子属性

#### MolecularWeight

```go
func (m *Molecule) MolecularWeight() (float64, error)
```

计算分子量。

**返回值:**

- `float64`: 分子量（g/mol）
- `error`: 错误信息

#### MonoisotopicMass

```go
func (m *Molecule) MonoisotopicMass() (float64, error)
```

计算单同位素质量。

#### MostAbundantMass

```go
func (m *Molecule) MostAbundantMass() (float64, error)
```

计算最丰富同位素质量。

#### GrossFormula

```go
func (m *Molecule) GrossFormula() (string, error)
```

获取总分子式。

**返回值:**

- `string`: 分子式（如 "C2H6O"）
- `error`: 错误信息

#### MolecularFormula

```go
func (m *Molecule) MolecularFormula() (string, error)
```

获取分子式。

#### TPSA

```go
func (m *Molecule) TPSA(includeSP bool) (float64, error)
```

计算拓扑极性表面积。

**参数:**

- `includeSP` (bool): 是否包含硫和磷

**返回值:**

- `float64`: TPSA 值（Ų）
- `error`: 错误信息

#### NumRotatableBonds

```go
func (m *Molecule) NumRotatableBonds() (int, error)
```

计算可旋转键数量。

### 分子保存

#### ToSmiles

```go
func (m *Molecule) ToSmiles() (string, error)
```

转换为 SMILES 字符串。

#### ToCanonicalSmiles

```go
func (m *Molecule) ToCanonicalSmiles() (string, error)
```

转换为规范 SMILES 字符串。

#### ToMolfile

```go
func (m *Molecule) ToMolfile() (string, error)
```

转换为 MOL 文件格式字符串。

#### SaveToFile

```go
func (m *Molecule) SaveToFile(filename string) error
```

保存分子到文件。

**参数:**

- `filename` (string): 输出文件路径

### 分子构建

#### AddAtom

```go
func (m *Molecule) AddAtom(symbol string) (int, error)
```

添加原子到分子。

**参数:**

- `symbol` (string): 元素符号（如 "C", "N", "O"）

**返回值:**

- `int`: 原子句柄
- `error`: 错误信息

#### AddBond

```go
func (m *Molecule) AddBond(source int, destination int, order int) (int, error)
```

在两个原子之间添加化学键。

**参数:**

- `source` (int): 源原子句柄
- `destination` (int): 目标原子句柄
- `order` (int): 键级（BOND_SINGLE, BOND_DOUBLE, BOND_TRIPLE, BOND_AROMATIC）

**返回值:**

- `int`: 化学键句柄
- `error`: 错误信息

#### SetCharge

```go
func SetCharge(atomHandle int, charge int) error
```

设置原子电荷。

**参数:**

- `atomHandle` (int): 原子句柄
- `charge` (int): 电荷值

#### SetIsotope

```go
func SetIsotope(atomHandle int, isotope int) error
```

设置原子同位素。

**参数:**

- `atomHandle` (int): 原子句柄
- `isotope` (int): 同位素质量数

### InChI 功能

#### InitInChI

```go
func InitInChI() error
```

初始化 InChI 模块。

#### DisposeInChI

```go
func DisposeInChI() error
```

释放 InChI 模块。

#### ToInChI

```go
func (m *Molecule) ToInChI() (string, error)
```

生成 InChI 字符串。

**返回值:**

- `string`: InChI 字符串
- `error`: 错误信息

#### ToInChIKey

```go
func (m *Molecule) ToInChIKey() (string, error)
```

生成 InChIKey。

**返回值:**

- `string`: InChIKey（27字符）
- `error`: 错误信息

#### LoadFromInChI

```go
func LoadFromInChI(inchi string) (*Molecule, error)
```

从 InChI 字符串加载分子。

**参数:**

- `inchi` (string): InChI 字符串

## Reaction 包

### 类型

#### Reaction

```go
type Reaction struct {
    handle int
    closed bool
}
```

### 创建和加载

#### CreateReaction

```go
func CreateReaction() (*Reaction, error)
```

创建空反应。

#### LoadReactionFromString

```go
func LoadReactionFromString(data string) (*Reaction, error)
```

从 Reaction SMILES 加载反应。

#### LoadReactionFromFile

```go
func LoadReactionFromFile(filename string) (*Reaction, error)
```

从 RXN 文件加载反应。

### 反应操作

#### AddReactant

```go
func (r *Reaction) AddReactant(moleculeHandle int) error
```

添加反应物。

#### AddProduct

```go
func (r *Reaction) AddProduct(moleculeHandle int) error
```

添加产物。

#### AddCatalyst

```go
func (r *Reaction) AddCatalyst(moleculeHandle int) error
```

添加催化剂。

#### CountReactants

```go
func (r *Reaction) CountReactants() (int, error)
```

返回反应物数量。

#### CountProducts

```go
func (r *Reaction) CountProducts() (int, error)
```

返回产物数量。

#### Automap

```go
func (r *Reaction) Automap(mode string) error
```

自动原子映射。

**参数:**

- `mode` (string): 映射模式（"discard", "keep", "alter", "clear"）

#### SaveRxnfileToFile

```go
func (r *Reaction) SaveRxnfileToFile(filename string) error
```

保存为 RXN 文件。

## Render 包

### 初始化

#### InitRenderer

```go
func InitRenderer() error
```

初始化渲染器。

#### DisposeRenderer

```go
func DisposeRenderer() error
```

释放渲染器。

### 渲染功能

#### RenderToFile

```go
func RenderToFile(objectHandle int, filename string) error
```

渲染对象到文件。

**参数:**

- `objectHandle` (int): 分子或反应句柄
- `filename` (string): 输出文件路径

#### RenderGridToFile

```go
func RenderGridToFile(arrayHandle int, refAtoms []int, nColumns int, filename string) error
```

渲染网格到文件。

**参数:**

- `arrayHandle` (int): 分子数组句柄
- `refAtoms` ([]int): 参考原子索引
- `nColumns` (int): 列数
- `filename` (string): 输出文件路径

### 渲染选项

#### SetRenderOption

```go
func SetRenderOption(option string, value string) error
```

设置渲染选项。

#### SetRenderOptionInt

```go
func SetRenderOptionInt(option string, value int) error
```

设置整数选项。

#### SetRenderOptionFloat

```go
func SetRenderOptionFloat(option string, value float64) error
```

设置浮点数选项。

#### SetRenderOptionBool

```go
func SetRenderOptionBool(option string, value bool) error
```

设置布尔选项。

### 渲染选项结构

#### RenderOptions

```go
type RenderOptions struct {
    OutputFormat      string
    ImageWidth        int
    ImageHeight       int
    BackgroundColor   string
    BondLength        int
    RelativeThickness float64
    ShowAtomIDs       bool
    ShowBondIDs       bool
    Margins           string
    StereoStyle       string
    LabelMode         string
}
```

#### DefaultRenderOptions

```go
func DefaultRenderOptions() *RenderOptions
```

返回默认渲染选项。

## 常量

### 化学键类型

```go
const (
    BOND_SINGLE   = 1  // 单键
    BOND_DOUBLE   = 2  // 双键
    BOND_TRIPLE   = 3  // 三键
    BOND_AROMATIC = 4  // 芳香键
)
```

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
)
```

### 特殊元素

```go
const (
    ELEM_PSEUDO   = -1  // 伪原子
    ELEM_RSITE    = -2  // R-site
    ELEM_TEMPLATE = -3  // 模板原子
)
```

### 自由基类型

```go
const (
    RADICAL_SINGLET = 2
    RADICAL_DOUBLET = 3
    RADICAL_TRIPLET = 4
)
```

## 错误处理

所有可能失败的操作都返回 error。始终检查错误：

```go
mol, err := molecule.LoadMoleculeFromString("CCO")
if err != nil {
    return fmt.Errorf("加载分子失败: %w", err)
}
defer mol.Close()
```

### 常见错误

- `"molecule is closed"`: 尝试操作已关闭的分子
- `"invalid handle"`: 无效的对象句柄
- `"failed to load molecule"`: 加载分子失败
- `"failed to parse SMILES"`: SMILES 解析错误

## 最佳实践

1. **资源管理**: 总是使用 `defer obj.Close()`
2. **错误检查**: 检查所有错误返回值
3. **初始化**: 使用 InChI 前调用 `InitInChI()`
4. **并发**: Molecule 对象不是并发安全的

## 索引

快速查找 API：

- [AddAtom](#addatom)
- [AddBond](#addbond)
- [Aromatize](#aromatize)
- [Clone](#clone)
- [Close](#close)
- [CountAtoms](#countatoms)
- [CountBonds](#countbonds)
- [CreateMolecule](#createmolecule)
- [InitInChI](#initinchi)
- [LoadMoleculeFromString](#loadmoleculefromstring)
- [MolecularWeight](#molecularweight)
- [RenderToFile](#rendertofile)
- [ToInChI](#toinchi)
- [ToSmiles](#tosmiles)

---

💡 **提示**: 使用 `Ctrl+F` 或 `Cmd+F` 快速搜索API。
