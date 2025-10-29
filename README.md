# Go-Chem - 化学分子结构处理库

这是一个用Go语言编写的化学分子结构处理库，基于Indigo C++库的设计理念重写。

## 项目概述

本项目将molecule文件夹中的C++代码系统地转换为Go版本，提供完整的分子操作、分析和文件格式支持。

## 已实现功能

### 核心分子操作 (`src/molecule.go`)
- ✅ 基本分子结构（原子、化学键、顶点）
- ✅ 原子属性管理（电荷、同位素、自由基）
- ✅ 化学键管理（单键、双键、三键、芳香键）
- ✅ 2D/3D坐标支持
- ✅ 隐式氢原子计算
- ✅ 分子量计算
- ✅ 邻居原子和化学键查询
- ✅ 分子克隆和编辑版本跟踪
- ✅ 伪原子和模板原子支持

### 元素数据 (`src/elements.go`)
- ✅ 完整的元素周期表数据
- ✅ 元素符号与原子序数转换
- ✅ 元素属性查询（族、周期、芳香性）
- ✅ 价电子和轨道计算

### 立体化学

#### 立体中心 (`src/molecule_stereocenters.go`)
- ✅ 立体中心检测和管理
- ✅ 金字塔构型表示
- ✅ 立体中心类型（ABS、OR、AND、ANY）
- ✅ 从3D坐标检测立体中心
- ✅ 从化学键方向构建立体中心
- ✅ 立体中心反转和操作

#### 顺反异构 (`src/molecule_cis_trans.go`)
- ✅ 顺反（E/Z）立体化学管理
- ✅ 几何立体化学键检测
- ✅ 从3D坐标确定构型
- ✅ 从化学键方向确定构型
- ✅ 取代基分析

### 文件格式支持

#### MOL文件 (`src/molfile_loader.go`, `src/molfile_saver.go`)
- ✅ MDL Molfile V2000格式加载
- ✅ MOL文件保存
- ✅ 原子属性解析（坐标、电荷、同位素、自由基）
- ✅ 化学键类型和立体化学
- ✅ M行属性块（CHG、ISO、RAD）
- ✅ 手性标志支持

#### SDF文件 (`src/sdf_loader.go`)
- ✅ SDF（结构数据文件）多分子加载
- ✅ 数据字段解析
- ✅ 批量分子加载

### 芳香化处理
- ✅ 芳香化算法 (`src/aromatizer.go`)
- ✅ 去芳香化算法 (`src/dearomatizer.go`)
- ✅ 苯环识别和处理

### 分子性质计算
- ✅ Lipinski五规则 (`src/lipinski.go`)
- ✅ TPSA（拓扑极性表面积）(`src/tpsa.go`)
- ✅ 总分子式生成 (`src/gross_formula.go`)
- ✅ 分子哈希 (`src/molecule_hash.go`)

### SMILES支持
- ✅ SMILES加载器 (`src/smiles_loader.go`)
- ✅ 芳香原子和化学键解析
- ✅ 环结构识别

### 其他格式
- ✅ CML（化学标记语言）基础支持 (`src/cml.go`)
- ✅ CDXML格式基础支持 (`src/cdxml.go`)

## 测试覆盖

### 分子基础测试 (`test/molecule_test.go`)
- ✅ 基本分子操作测试
- ✅ 原子属性测试
- ✅ 化学键操作测试
- ✅ 坐标处理测试
- ✅ 分子克隆测试
- ✅ 分子量计算测试
- ✅ 隐式氢测试

### MOL文件测试 (`test/molfile_test.go`)
- ✅ MOL文件加载测试
- ✅ 带电荷的分子测试
- ✅ 同位素测试
- ✅ MOL文件保存测试
- ✅ 往返测试（加载-保存-加载）
- ✅ 不同化学键类型测试
- ✅ 伪原子测试
- ✅ 立体化学测试

### 立体化学测试 (`test/stereochemistry_test.go`)
- ✅ 立体中心基础操作
- ✅ 金字塔构型测试
- ✅ 不同立体中心类型
- ✅ 立体中心检测
- ✅ 从3D坐标检测
- ✅ 顺反异构基础操作
- ✅ 顺反异构检测
- ✅ 构型字符串表示

### 现有测试优化
- ✅ 芳香化测试 (`test/aromatizer_test.go`)
- ✅ 化学基础测试 (`test/chem_test.go`)
- ✅ 总分子式测试 (`test/gross_formula_test.go`)
- ✅ 分子性质测试 (`test/properties_test.go`)
- ✅ SMILES加载器测试 (`test/smiles_loader_test.go`)

## 项目结构

```
go-chem/
├── src/                          # 源代码
│   ├── molecule.go               # 核心分子结构
│   ├── molecule_stereocenters.go # 立体中心
│   ├── molecule_cis_trans.go     # 顺反异构
│   ├── molfile_loader.go         # MOL文件加载
│   ├── molfile_saver.go          # MOL文件保存
│   ├── sdf_loader.go             # SDF文件加载
│   ├── elements.go               # 元素数据
│   ├── aromatizer.go             # 芳香化
│   ├── dearomatizer.go           # 去芳香化
│   ├── smiles_loader.go          # SMILES加载
│   ├── lipinski.go               # Lipinski规则
│   ├── tpsa.go                   # TPSA计算
│   ├── gross_formula.go          # 总分子式
│   ├── molecule_hash.go          # 分子哈希
│   └── ...                       # 其他模块
├── test/                         # 测试文件
│   ├── molecule_test.go          # 分子测试
│   ├── molfile_test.go           # MOL文件测试
│   ├── stereochemistry_test.go   # 立体化学测试
│   └── ...                       # 其他测试
├── molecule/                     # C++原始代码（参考）
│   ├── *.h                       # 头文件
│   └── src/                      # C++实现
└── go.mod                        # Go模块定义
```

## 使用示例

### 创建分子

```go
import "go-chem/src"

// 创建乙醇分子 (CH3CH2OH)
mol := src.NewMolecule()
mol.Name = "Ethanol"

// 添加原子
c1 := mol.AddAtom(src.ELEM_C)
c2 := mol.AddAtom(src.ELEM_C)
o := mol.AddAtom(src.ELEM_O)

// 添加化学键
mol.AddBond(c1, c2, src.BOND_SINGLE)
mol.AddBond(c2, o, src.BOND_SINGLE)

// 设置坐标
mol.SetAtomXYZ(c1, 0.0, 0.0, 0.0)
mol.SetAtomXYZ(c2, 1.5, 0.0, 0.0)
mol.SetAtomXYZ(o, 2.0, 1.0, 0.0)
```

### 加载MOL文件

```go
import (
    "os"
    "go-chem/src"
)

file, _ := os.Open("molecule.mol")
defer file.Close()

loader := src.NewMolfileLoader(file)
mol, err := loader.LoadMolecule()
if err != nil {
    panic(err)
}

fmt.Printf("Loaded molecule: %s\n", mol.Name)
fmt.Printf("Atoms: %d, Bonds: %d\n", mol.AtomCount(), mol.BondCount())
```

### 保存MOL文件

```go
import (
    "os"
    "go-chem/src"
)

file, _ := os.Create("output.mol")
defer file.Close()

saver := src.NewMolfileSaver(file)
err := saver.SaveMolecule(mol)
if err != nil {
    panic(err)
}
```

### 计算分子性质

```go
// 计算分子量
mw := mol.CalcMolecularWeight()
fmt.Printf("Molecular weight: %.2f\n", mw)

// 计算Lipinski规则
lipinski := src.NewLipinskiCalculator()
lipinski.Calculate(mol)
fmt.Printf("HBD: %d, HBA: %d\n", lipinski.HBD, lipinski.HBA)

// 计算TPSA
tpsa := src.CalculateTPSA(mol)
fmt.Printf("TPSA: %.2f\n", tpsa)
```

### 立体化学

```go
// 创建立体中心管理器
stereo := src.NewMoleculeStereocenters()

// 从3D坐标检测立体中心
stereo.BuildFrom3DCoordinates(mol)

// 检查原子是否为立体中心
if stereo.Exists(atomIdx) {
    center, _ := stereo.Get(atomIdx)
    fmt.Printf("Atom %d is a stereocenter, type: %d\n", atomIdx, center.Type)
}
```

## 运行测试

```bash
# 运行所有测试
go test ./test/...

# 运行特定测试
go test ./test/ -run TestMoleculeBasics

# 带详细输出
go test -v ./test/...
```

### 分子指纹 (`src/molecule_fingerprint.go`)
- ✅ 基于路径的指纹（类似Daylight）
- ✅ ECFP（Extended Connectivity Fingerprints）支持
- ✅ ECFP2, ECFP4, ECFP6多种半径
- ✅ Tanimoto相似度计算
- ✅ Dice系数
- ✅ Cosine相似度
- ✅ Hamming距离和欧氏距离
- ✅ 十六进制字符串转换

### 子结构匹配 (`src/molecule_substructure_matcher.go`)
- ✅ 完整的子图同构算法
- ✅ 递归回溯搜索
- ✅ 原子和化学键匹配
- ✅ 查找所有匹配
- ✅ 查找第一个匹配（快速模式）
- ✅ 匹配计数
- ✅ 便捷函数接口
- ✅ 最大公共子结构框架

### S-Groups支持 (`src/molecule_sgroups.go`)
- ✅ 通用S-Group（GEN）
- ✅ 数据S-Group（DAT）
- ✅ 超原子/缩写（SUP）
- ✅ 结构重复单元（SRU）
- ✅ 多重组（MUL）
- ✅ 聚合物S-Group（MON, MER, COP等）
- ✅ 括号和显示选项
- ✅ S-Group层次结构管理
- ✅ 原子和化学键移除时的更新

## 待实现功能（未来计划）

根据原始C++代码，以下功能可在未来版本中实现：

- 🔲 Morgan指纹变体
- 🔲 完整的SMARTS模式匹配
- 🔲 互变异构体生成和匹配
- 🔲 InChI生成和解析
- 🔲 更多文件格式（RDF、RXN、V3000等）
- 🔲 3D构象生成
- 🔲 力场和能量最小化

## 技术特点

1. **纯Go实现**：无CGO依赖，易于跨平台部署
2. **惰性计算**：分子属性按需计算并缓存
3. **编辑跟踪**：自动跟踪分子修改
4. **内存高效**：使用切片和映射优化内存使用
5. **类型安全**：利用Go的类型系统确保正确性
6. **完整测试**：全面的单元测试覆盖

## 性能考虑

- 原子和化学键使用索引而非指针，避免GC压力
- 缓存常用属性（连接性、隐式氢等）
- 使用编辑版本号避免不必要的重新计算
- 向量和几何计算使用内联函数

## 贡献

欢迎贡献代码、报告问题或提出改进建议！

## 许可证

本项目基于Apache License 2.0许可证（与Indigo toolkit相同）。

## 致谢

本项目的设计和API受到[EPAM Indigo toolkit](https://github.com/epam/Indigo)的启发。

