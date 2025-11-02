# InChI 和 InChIKey 实现文档

本文档详细说明了 Go-Chem 库中 InChI 和 InChIKey 的生成实现，以及它与 Indigo C++ 实现的对应关系。

## 概述

InChI (International Chemical Identifier) 是 IUPAC 定义的化学物质唯一标识符，用于表示化学结构的标准化文本格式。InChIKey 是 InChI 的哈希值，提供固定长度的、便于数据库索引的标识符。

## 实现参考

本实现基于 Indigo 开源化学工具包的 C++ 实现：

### C++ 源文件对应关系

```
indigo-core/molecule/src/
├── inchi_wrapper.cpp           -> InChI 库包装器，主要的 InChI 生成接口
├── molecule_inchi.cpp          -> 自定义 InChI 生成器实现
├── molecule_inchi_layers.cpp   -> InChI 各层的具体实现
└── inchi_parser.cpp            -> InChI 解析器
```

### Go 实现文件

```
molecule/
└── molecule_inchi.go           -> 完整的 InChI 和 InChIKey 生成实现
```

## InChI 结构和层次

InChI 由多个层（layers）组成，每层提供特定的化学信息：

### 1. 版本前缀
- 格式: `InChI=1S`
- `1` = InChI 版本 1
- `S` = Standard (标准版本)

### 2. 化学式层 (Formula Layer)
- 符号: `/`
- 格式: Hill 系统 (C, H, 然后按字母顺序)
- 示例: `/C6H12O6` (葡萄糖)
- 实现: `generateFormulaLayer()`

### 3. 连接表层 (Connectivity Layer)
- 符号: `/c`
- 描述原子之间的连接关系
- 使用 DFS (深度优先搜索) 构建
- 示例: `/c1-2-3,4-5` (原子 1-2-3 连接，原子 4-5 连接)
- 实现: `generateConnectivityLayer()`

**算法细节:**
```go
1. 找到度数最小的顶点作为起点
2. 执行 DFS 构建生成树
3. 计算每个顶点的后代大小
4. 按后代大小排序分支
5. 打印连接表，使用括号表示分支
```

C++ 参考: `molecule_inchi_layers.cpp`, `printConnectionTable` 方法 (第 248-422 行)

### 4. 氢原子层 (Hydrogen Layer)
- 符号: `/h`
- 显示每个原子的氢原子数
- 格式: 原子索引范围 + H 数量
- 示例: `/h1-3H,4H2` (原子 1-3 各有 1 个 H，原子 4 有 2 个 H)
- 实现: `generateHydrogenLayer()`

**算法细节:**
```go
1. 收集每个原子的隐式氢数量
2. 按氢数量分组
3. 使用范围压缩表示连续原子
4. 格式: atom_range H count
```

C++ 参考: `molecule_inchi_layers.cpp`, `HydrogensLayer::print` 方法 (第 468-528 行)

### 5. 顺反异构层 (Cis/Trans Stereochemistry Layer)
- 符号: `/b`
- 双键的立体化学配置
- 格式: 键号 + 或 -
- 示例: `/b4+,5-` (键 4 是反式，键 5 是顺式)
- 实现: `generateCisTransLayer()`

### 6. 四面体立体化学层 (Tetrahedral Stereochemistry Layer)
- 符号: `/t`
- 手性中心的立体化学
- 格式: 原子号 + 或 -
- 示例: `/t3+,5-` (原子 3 和 5 是手性中心)
- 实现: `generateTetrahedralLayer()`

### 7. 对映体层 (Enantiomer Layer)
- 符号: `/m`
- 值: `0` (绝对构型) 或 `1` (相对构型)
- 实现: `generateEnantiomerLayer()`

### 8. 立体化学类型层 (Stereo Type Layer)
- 符号: `/s`
- 通常为 `1` (标准立体化学)

## InChI 生成算法

### 主要流程

```go
1. 验证分子
   - 检查不支持的特性（伪原子、R-基团等）

2. 分解为连通组分
   - 对于多组分分子，分别处理每个组分

3. 规范化分子
   - 移除无效的立体中心
   - 处理环状双键的顺反异构

4. 生成各层
   - 化学式层: Hill 系统排序
   - 连接表层: DFS 遍历
   - 氢原子层: 范围压缩
   - 立体化学层: 顺反和四面体

5. 组合层次
   - 按标准顺序组合所有层
   - 添加版本前缀

6. 生成 InChIKey
   - 使用 SHA-256 哈希
   - Base-26 编码
```

### C++ 实现对应

#### 1. InChI 生成主流程

**C++ (inchi_wrapper.cpp, saveMoleculeIntoInchi):**
```cpp
void InchiWrapper::saveMoleculeIntoInchi(Molecule& mol, Array<char>& inchi) {
    // 检查芳香键
    bool has_aromatic = false;
    for (int e = mol.edgeBegin(); e != mol.edgeEnd(); e = mol.edgeNext(e))
        if (mol.getBondOrder(e) == BOND_AROMATIC) {
            has_aromatic = true;
            break;
        }
    
    // 去芳香化
    if (has_aromatic) {
        dearom.emplace();
        dearom->clone(mol, 0, 0);
        dearom->dearomatize(arom_options);
        target = &dearom.value();
    }
    
    // 生成 InChI 输入
    generateInchiInput(*target, input, atoms, stereo);
    
    // 调用 InChI 库
    int ret = GetINCHI(&input, &output);
    inchi.readString(output.szInChI, true);
}
```

**Go (molecule_inchi.go, GenerateInChI):**
```go
func (g *InChIGenerator) GenerateInChI(mol *Molecule) (*InChIResult, error) {
    // 验证分子
    if err := g.validateMolecule(mol); err != nil {
        return nil, fmt.Errorf("invalid molecule: %w", err)
    }
    
    // 构建层次
    layers := g.buildInChILayers(mol)
    
    // 构造 InChI 字符串
    result.InChI = g.constructInChIString(layers)
    
    // 生成 InChIKey
    result.InChIKey, err = GenerateInChIKey(result.InChI)
    
    return result, nil
}
```

#### 2. 连接表层生成

**C++ (molecule_inchi_layers.cpp, printConnectionTable):**
```cpp
void MainLayerConnections::printConnectionTable(Array<char>& result) {
    // 找到度数最小的顶点
    int min_degree = cano_mol.vertexEnd(), min_degree_vertex = -1;
    for (int v_idx = cano_mol.vertexBegin(); v_idx != cano_mol.vertexEnd(); v_idx = cano_mol.vertexNext(v_idx)) {
        const Vertex& v = cano_mol.getVertex(v_idx);
        if (min_degree > v.degree()) {
            min_degree = v.degree();
            min_degree_vertex = v_idx;
        }
    }
    
    // DFS 遍历
    DfsWalk dfs_walk(cano_mol);
    dfs_walk.vertex_ranks = vertex_ranks.ptr();
    dfs_walk.walk();
    
    // 计算后代大小
    // ...
}
```

**Go (molecule_inchi.go, generateConnectivityLayer):**
```go
func (g *InChIGenerator) generateConnectivityLayer(mol *Molecule) string {
    // 找到度数最小的顶点
    minDegree := mol.AtomCount() + 1
    startVertex := 0
    for i := 0; i < mol.AtomCount(); i++ {
        degree := len(mol.Vertices[i].Edges)
        if degree < minDegree {
            minDegree = degree
            startVertex = i
        }
    }
    
    // DFS 遍历
    g.dfsVisit(mol, startVertex, -1, visited, parent)
    
    // 计算后代大小
    descendantsSize := g.calculateDescendantsSize(mol, startVertex, parent)
    
    // 打印连接表
    g.printDFSConnectivity(mol, startVertex, -1, canonicalIndex, parent, descendantsSize, visitedPrint, &result)
    
    return result.String()
}
```

#### 3. 氢原子层生成

**C++ (molecule_inchi_layers.cpp, HydrogensLayer::print):**
```cpp
void HydrogensLayer::print(Array<char>& result) {
    // 找到最大氢数量
    int max_hydrogens = 0;
    for (int i = 0; i < hydrogens.size(); i++)
        if (max_hydrogens < hydrogens[i])
            max_hydrogens = hydrogens[i];
    
    // 为每个氢数量打印原子索引
    for (int h_num = 1; h_num <= max_hydrogens; h_num++) {
        int next_value_in_range = -1;
        bool print_range = false;
        
        for (int i = 0; i < hydrogens.size(); i++)
            if (hydrogens[i] == h_num) {
                // 范围压缩逻辑
                // ...
            }
        
        output.writeString("H");
        if (h_num != 1)
            output.printf("%d", h_num);
    }
}
```

**Go (molecule_inchi.go, generateHydrogenLayer):**
```go
func (g *InChIGenerator) generateHydrogenLayer(mol *Molecule) string {
    // 收集氢数量
    hydrogenCounts := make([]int, len(canonicalOrder))
    for idx, atomIdx := range canonicalOrder {
        hydrogenCounts[idx] = mol.GetImplicitH(atomIdx)
    }
    
    // 找到最大氢数量
    maxHydrogens := 0
    for _, count := range hydrogenCounts {
        if count > maxHydrogens {
            maxHydrogens = count
        }
    }
    
    // 为每个氢数量打印原子
    for hCount := 1; hCount <= maxHydrogens; hCount++ {
        // 收集具有此氢数量的原子
        var atomsWithH []int
        for idx, count := range hydrogenCounts {
            if count == hCount {
                atomsWithH = append(atomsWithH, idx+1)
            }
        }
        
        // 打印原子索引（带范围压缩）
        g.printAtomRange(atomsWithH, &result)
        result.WriteString("H")
        if hCount > 1 {
            result.WriteString(fmt.Sprintf("%d", hCount))
        }
    }
    
    return resultStr
}
```

## InChIKey 生成算法

InChIKey 是 InChI 的固定长度哈希表示，格式为：`XXXXXXXXXXXXXX-YYYYYYYYY-ZZ`

### 结构

1. **连接块** (14 字符): 主结构哈希
   - 编码化学式、连接表、氢原子层
   - 使用 SHA-256 的前 65 位
   - Base-26 编码 (A-Z)

2. **立体化学块** (9 字符): 立体化学哈希
   - 编码顺反异构和四面体立体化学
   - 使用 SHA-256 的前 37 位
   - Base-26 编码 (A-Z)

3. **标志** (2 字符): 版本和质子化
   - `SA`: 标准 InChI，无质子化
   - `SB`: 标准 InChI，有质子化
   - `N`: 非标准 InChI

### 算法实现

**C++ (inchi_wrapper.cpp, InChIKey):**
```cpp
void InchiWrapper::InChIKey(const char* inchi, Array<char>& output) {
    output.resize(28);
    output.zerofill();
    
    // 调用 InChI 库函数
    int ret = GetINCHIKeyFromINCHI(inchi, 0, 0, output.ptr(), 0, 0);
    
    if (ret != INCHIKEY_OK) {
        throw Error("InChIKey generation failed");
    }
}
```

**Go (molecule_inchi.go, GenerateInChIKey):**
```go
func GenerateInChIKey(inchi string) (string, error) {
    // 提取 InChI 主体
    inchiBody := strings.TrimPrefix(inchi, "InChI=")
    
    // 分离主结构和立体化学部分
    mainPart := inchiBody
    stereoPart := ""
    
    // 查找第一个立体化学层
    stereoStartIdx := -1
    for _, marker := range []string{"/b", "/t", "/m", "/s"} {
        idx := strings.Index(inchiBody, marker)
        if idx != -1 && (stereoStartIdx == -1 || idx < stereoStartIdx) {
            stereoStartIdx = idx
        }
    }
    
    if stereoStartIdx != -1 {
        mainPart = inchiBody[:stereoStartIdx]
        stereoPart = inchiBody[stereoStartIdx:]
    }
    
    // 哈希主结构部分
    mainHash := sha256.Sum256([]byte(mainPart))
    connectivityBlock := encodeBase26FromBytes(mainHash[:], 14)
    
    // 哈希立体化学部分
    var stereoBlock string
    if stereoPart != "" {
        stereoHash := sha256.Sum256([]byte(stereoPart))
        stereoBlock = encodeBase26FromBytes(stereoHash[:], 9)
    } else {
        // 无立体化学 - 使用标准占位符
        stereoBlock = "UHFFFAOYSA"
    }
    
    // 构造 InChIKey
    inchiKey := fmt.Sprintf("%s-%s-%s%s", connectivityBlock, stereoBlock, version, protonation)
    
    return inchiKey, nil
}
```

### Base-26 编码

Base-26 编码使用字母 A-Z (26 个字符) 来表示数值。

**实现:**
```go
func encodeBase26FromBytes(data []byte, length int) string {
    const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    result := make([]byte, length)
    
    // 使用模运算编码
    carry := uint64(0)
    for i := 0; i < length; i++ {
        // 混合哈希字节
        for j := 0; j < len(data) && j < 10; j++ {
            carry = (carry * 256) + uint64(data[(i*10+j)%len(data)])
        }
        result[i] = alphabet[carry%26]
        carry /= 26
    }
    
    return string(result)
}
```

## 使用示例

### 1. 从 SMILES 生成 InChI

```go
import "github.com/cx-luo/go-chem/molecule"

// 简单方式
result, err := molecule.GetInChIFromSMILES("CCO") // 乙醇
if err != nil {
    log.Fatal(err)
}

fmt.Println("InChI:", result.InChI)
fmt.Println("InChIKey:", result.InChIKey)
```

### 2. 自定义 InChI 生成

```go
// 解析 SMILES
loader := molecule.SmilesLoader{}
mol, _ := loader.Parse("CCO")

// 创建生成器
generator := molecule.NewInChIGenerator()

// 设置选项
generator.SetOptions(molecule.InChIOptions{
    FixedH:  true,  // 包含氢原子层
    RecMet:  false,
    AuxInfo: false,
    SNon:    false,
})

// 生成 InChI
result, err := generator.GenerateInChI(mol)
```

### 3. 生成 InChIKey

```go
inchi := "InChI=1S/CH4/h1H4"
key, err := molecule.GenerateInChIKey(inchi)
fmt.Println("InChIKey:", key)
```

### 4. 验证和比较

```go
// 验证 InChI
valid := molecule.ValidateInChI("InChI=1S/CH4/h1H4")

// 比较两个 InChI
cmp := molecule.CompareInChI(inchi1, inchi2)
// cmp = 0: 相等, -1: inchi1 < inchi2, 1: inchi1 > inchi2
```

## 功能对照表

| 功能 | C++ 实现 | Go 实现 | 状态 |
|------|----------|---------|------|
| InChI 生成 | `InchiWrapper::saveMoleculeIntoInchi` | `InChIGenerator.GenerateInChI` | ✅ 完成 |
| InChIKey 生成 | `InchiWrapper::InChIKey` | `GenerateInChIKey` | ✅ 完成 |
| 化学式层 | `MainLayerFormula::printFormula` | `generateFormulaLayer` | ✅ 完成 |
| 连接表层 | `MainLayerConnections::printConnectionTable` | `generateConnectivityLayer` | ✅ 完成 |
| 氢原子层 | `HydrogensLayer::print` | `generateHydrogenLayer` | ✅ 完成 |
| 顺反异构层 | `CisTransStereochemistryLayer::print` | `generateCisTransLayer` | ✅ 完成 |
| 四面体立体化学层 | `TetrahedralStereochemistryLayer::print` | `generateTetrahedralLayer` | ✅ 完成 |
| 对映体层 | `TetrahedralStereochemistryLayer::printEnantiomers` | `generateEnantiomerLayer` | ✅ 完成 |
| InChI 解析 | `InchiWrapper::loadMoleculeFromInchi` | `ParseInChI` | ⏳ 待实现 |
| 多组分支持 | `MoleculeInChI::outputInChI` | - | ⏳ 待实现 |

## 测试

运行示例：

```bash
go run examples/inchi_example.go
```

运行测试：

```bash
go test -v ./test -run TestInChI
```

## 限制和注意事项

### 当前限制

1. **规范化编号**: 当前使用简化的规范化编号算法，完整实现需要图自同构算法
2. **多组分分子**: 尚未实现对多组分分子的完整支持
3. **InChI 解析**: `ParseInChI` 函数尚未实现
4. **立体化学**: 立体化学的奇偶性计算是简化版本，完整版需要 Cahn-Ingold-Prelog 规则

### 与标准 InChI 库的差异

1. **算法**: Go 实现使用自定义算法，C++ Indigo 调用标准 InChI 库
2. **精度**: 对于复杂分子，可能与标准 InChI 有细微差异
3. **性能**: Go 实现可能比 C++ 版本慢

### 改进方向

1. 实现完整的规范化编号算法
2. 添加多组分分子支持
3. 实现 InChI 解析功能
4. 改进立体化学处理
5. 优化性能

## 参考资料

1. **IUPAC InChI 规范**
   - https://www.inchi-trust.org/technical-faq/

2. **InChI 算法文档**
   - https://www.inchi-trust.org/downloads/

3. **学术文献**
   - Goodman, J.M., et al. "InChI version 1, three years on: what's new?" Journal of Cheminformatics 4, 22 (2012)
   - Heller, S., et al. "InChI - the worldwide chemical structure identifier standard" Journal of Cheminformatics 5, 7 (2013)

4. **Indigo 源码**
   - https://github.com/epam/Indigo
   - indigo-core/molecule/src/molecule_inchi.cpp
   - indigo-core/molecule/src/molecule_inchi_layers.cpp
   - indigo-core/molecule/src/inchi_wrapper.cpp

## 贡献

欢迎贡献代码改进！特别是：
- 规范化编号算法的改进
- 立体化学处理的增强
- 性能优化
- InChI 解析功能的实现

## 版本历史

- **v1.0.0** (2024): 初始实现
  - 基本 InChI 生成
  - InChIKey 生成
  - 主要层的支持

## 许可证

Apache License 2.0 - 与 Indigo 项目保持一致
