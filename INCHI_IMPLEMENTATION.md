# InChI 和 InChIKey 实现文档

## 概述

本文档描述了 go-chem 项目中 InChI (IUPAC International Chemical Identifier) 和 InChIKey 的实现方法。该实现基于 Indigo C++ 源码和 IUPAC InChI 官方规范。

## InChI 简介

InChI 是由国际纯粹与应用化学联合会 (IUPAC) 开发的化学物质标准化文本标识符。它提供了一种独特且标准的方式来表示化学结构信息。

### InChI 的层级结构

InChI 采用分层设计，每一层提供特定类型的化学信息：

1. **分子式层 (Formula Layer)** - `/`
   - 使用 Hill 系统排序：C，H，然后其他元素按字母顺序
   - 示例：`/C6H12O6/`

2. **连接层 (Connectivity Layer)** - `/c`
   - 描述原子之间的连接关系
   - 使用规范化的原子编号
   - 示例：`/c1-2-3-4-5-6/`

3. **氢原子层 (Hydrogen Layer)** - `/h`
   - 描述隐式氢原子的分布
   - 示例：`/h1-2H,3H2/`

4. **双键立体化学层 (Double Bond Stereochemistry Layer)** - `/b`
   - 描述双键的顺反异构 (cis/trans)
   - 示例：`/b3-4+/` (trans) 或 `/b3-4-/` (cis)

5. **四面体立体化学层 (Tetrahedral Stereochemistry Layer)** - `/t`
   - 描述手性中心的构型
   - 示例：`/t2-,3+/`

6. **对映异构体层 (Enantiomer Layer)** - `/m`
   - 指示对映异构体信息
   - 0 = 绝对构型，1 = 相对构型

7. **立体化学类型层 (Stereo Type Layer)** - `/s`
   - 1 = 标准立体化学

### InChI 示例

```
InChI=1S/C6H12O6/c7-1-2-3(8)4(9)5(10)6(11)12-2/h2-11H,1H2/t2-,3-,4+,5-,6-/m1/s1
```

这是葡萄糖（glucose）的 InChI 表示。

## InChIKey 简介

InChIKey 是 InChI 的定长哈希表示，长度固定为 27 个字符（不包括连字符），便于数据库索引和搜索。

### InChIKey 格式

```
XXXXXXXXXXXXXX-YYYYYYYYY-ZZ
```

- **X (14 字符)**: 连接性哈希块
- **Y (9-10 字符)**: 立体化学哈希块
- **Z (1-2 字符)**: 版本和质子化标志

### InChIKey 算法

1. 分离 InChI 的主要部分和立体化学部分
2. 对每部分进行 SHA-256 哈希
3. 将哈希值的前 N 位转换为 Base26 编码（A-Z）
4. 添加版本标志（S = Standard InChI, N = Non-standard）

### InChIKey 示例

```
WQZGKKKJIJFFOK-GASJEMHNSA-N
```

这是葡萄糖的 InChIKey。

## 实现细节

### 核心文件

- `molecule/molecule_inchi.go`: 主要实现文件，包含 InChI 和 InChIKey 生成逻辑
- `test/inchi_test.go`: 测试文件，验证实现的正确性

### 主要类型和函数

#### InChIGenerator

```go
type InChIGenerator struct {
    prefix  string        // InChI 版本前缀，默认 "InChI=1S"
    options InChIOptions  // 生成选项
}
```

主要方法：
- `GenerateInChI(mol *Molecule) (*InChIResult, error)`: 从分子生成 InChI
- `SetOptions(options InChIOptions)`: 设置生成选项

#### InChIResult

```go
type InChIResult struct {
    InChI    string   // 生成的 InChI 字符串
    InChIKey string   // 生成的 InChIKey
    AuxInfo  string   // 辅助信息
    Warnings []string // 警告消息
    Log      []string // 日志消息
}
```

#### 关键函数

1. **`GenerateInChI(mol *Molecule) (*InChIResult, error)`**
   - 主入口函数
   - 验证分子结构
   - 构建各个 InChI 层
   - 生成最终的 InChI 字符串和 InChIKey

2. **`generateFormulaLayer(mol *Molecule) string`**
   - 生成分子式层
   - 使用 Hill 系统排序
   - 统计原子数（包括隐式氢）

3. **`generateConnectivityLayer(mol *Molecule) string`**
   - 生成连接层
   - 使用规范化原子编号
   - 采用 BFS/DFS 遍历构建连接字符串

4. **`generateHydrogenLayer(mol *Molecule) string`**
   - 生成氢原子层
   - 列出每个原子的隐式氢数量

5. **`GenerateInChIKey(inchi string) (string, error)`**
   - 从 InChI 字符串生成 InChIKey
   - 使用 SHA-256 哈希
   - Base26 编码

6. **`GetInChIFromSMILES(smiles string) (*InChIResult, error)`**
   - 便捷函数，直接从 SMILES 生成 InChI

## 算法实现

### 1. 分子式层生成

根据 Hill 系统规则：

```go
// 1. 碳 (C) 优先
// 2. 氢 (H) 其次
// 3. 其他元素按符号字母顺序
```

算法步骤：
1. 统计每种元素的原子数
2. 加上隐式氢原子数
3. 按 Hill 系统排序
4. 格式化为字符串（数量 > 1 时显示数字）

### 2. 连接层生成

算法步骤：
1. 创建规范化原子编号（基于原子序数、连接度等）
2. 对每个连通分量：
   - 从最高优先级原子开始
   - BFS/DFS 遍历所有连接的原子
   - 记录访问顺序
3. 多个分量用分号分隔

### 3. 氢原子层生成

算法步骤：
1. 遍历所有原子
2. 对有隐式氢的原子，记录 `原子编号,氢数量H`
3. 格式：`1,2H,3H2` 表示原子1有2个H，原子3有2个H

### 4. InChIKey 生成

算法步骤：
1. 移除 "InChI=" 前缀
2. 分离主要部分和立体化学部分
   - 主要部分：从开头到 `/t` 或 `/m` 或 `/s` 之前
   - 立体化学部分：从 `/t` 或 `/m` 或 `/s` 开始
3. 对每部分计算 SHA-256 哈希
4. 将哈希值转换为 Base26 编码：
   - 连接性块：14 字符
   - 立体化学块：9 字符
5. 添加版本标志：
   - "SA" = Standard InChI, no protonation
   - "N" = Non-standard

## 规范化和排序规则

### 原子规范化排序

基于以下优先级：
1. 原子序数（原子量）
2. 连接度（邻接原子数）
3. 键序和（单键=1，双键=2，三键=3，芳香键=1.5）
4. 环成员关系
5. 立体化学信息

### Hill 系统

化学式中元素的标准排序系统：
- C（碳）优先
- H（氢）其次
- 其他元素按元素符号字母顺序
- 每个元素后跟原子数（数量为1时省略）

示例：
- `CH4` （甲烷）
- `C2H6O` （乙醇）
- `C6H12O6` （葡萄糖）
- `H2O` （水，无碳时从H开始）
- `H2SO4` （硫酸）

## 参考文献

1. **IUPAC InChI Trust**
   - 官方网站: https://www.inchi-trust.org/
   - 技术文档: https://www.inchi-trust.org/downloads/

2. **InChI Technical Manual**
   - 详细描述了 InChI 的算法和规范
   - 提供了各层的详细说明

3. **Goodman et al. (2012)**
   - "InChI version 1, three years on"
   - Journal of Cheminformatics, 4:29
   - DOI: 10.1186/1758-2946-4-29
   - 描述了 InChI 的发展和应用

4. **Indigo Toolkit**
   - 开源化学信息学工具包
   - GitHub: https://github.com/epam/Indigo
   - 提供了 C++ 参考实现

5. **Hill System**
   - Hill, E. A. (1900). "On a system of indexing chemical literature; adopted by the classification division of the U.S. Patent Office"
   - Journal of the American Chemical Society, 22(8): 478-494

## 使用示例

### 基本使用

```go
import "github.com/cx-luo/go-chem/molecule"

// 从 SMILES 生成 InChI
result, err := molecule.GetInChIFromSMILES("CC(=O)O") // 乙酸
if err != nil {
    log.Fatal(err)
}

fmt.Println("InChI:", result.InChI)
fmt.Println("InChIKey:", result.InChIKey)
```

### 使用 InChIGenerator

```go
// 解析 SMILES
loader := molecule.SmilesLoader{}
mol, err := loader.Parse("c1ccccc1") // 苯
if err != nil {
    log.Fatal(err)
}

// 创建 InChI 生成器
generator := molecule.NewInChIGenerator()

// 设置选项（可选）
generator.SetOptions(molecule.InChIOptions{
    FixedH: true,  // 包含氢原子层
})

// 生成 InChI
result, err := generator.GenerateInChI(mol)
if err != nil {
    log.Fatal(err)
}

fmt.Println("InChI:", result.InChI)
fmt.Println("InChIKey:", result.InChIKey)
```

### 验证和比较 InChI

```go
// 验证 InChI
inchi := "InChI=1S/C6H6/c1-2-4-6-5-3-1/h1-6H"
valid := molecule.ValidateInChI(inchi)
fmt.Println("Valid:", valid)

// 比较两个 InChI
inchi1 := "InChI=1S/CH4/h1H4"
inchi2 := "InChI=1S/C2H6/c1-2/h1-2H3"
cmp := molecule.CompareInChI(inchi1, inchi2)
// cmp: 0 = 相同, -1 = inchi1 < inchi2, 1 = inchi1 > inchi2
```

## 当前实现状态

### 已实现的功能

- ✅ 基本框架和类型定义
- ✅ 分子式层生成（Hill 系统）
- ✅ 连接层生成（简化版规范化）
- ✅ 氢原子层生成
- ✅ InChIKey 生成（SHA-256 + Base26）
- ✅ SMILES 到 InChI 转换
- ✅ InChI 验证和比较
- ✅ Base64 编码/解码
- ✅ 测试套件

### 待完善的功能

- ⏳ 完整的规范化排序（需要实现图同构算法）
- ⏳ 双键立体化学层（cis/trans）
- ⏳ 四面体立体化学层（R/S 手性）
- ⏳ 对映异构体层
- ⏳ InChI 解析（反向操作）
- ⏳ 多组分分子处理
- ⏳ 同位素标记

### 限制和注意事项

1. **规范化**: 当前实现使用简化的规范化算法。完整的规范化需要图同构算法（automorphism search），这在 Indigo C++ 实现中有完整的实现。

2. **立体化学**: 立体化学层（/b, /t, /m）需要复杂的几何和拓扑分析，当前版本返回空字符串。

3. **多组分**: 当前实现假设单一连通分子。多组分分子（用 `.` 分隔的 SMILES）需要额外的处理逻辑。

4. **与官方 InChI 的兼容性**: 由于使用了简化算法，生成的 InChI 可能与官方 InChI 库的结果略有不同。对于需要完全兼容性的应用，建议通过 CGO 调用官方 InChI C 库。

## 性能优化建议

1. **缓存**: 对于重复计算的分子，可以缓存 InChI 结果
2. **并行处理**: 批量处理时可以使用 goroutine 并行生成
3. **内存池**: 使用 sync.Pool 重用临时缓冲区
4. **预计算**: 预先计算并缓存原子属性（连接度、价态等）

## 扩展方向

1. **集成官方 InChI 库**: 使用 CGO 调用官方 InChI C 库以获得完整功能和完全兼容性
2. **InChI 解析**: 实现从 InChI 字符串重建分子结构
3. **增强的立体化学**: 完整实现所有立体化学层
4. **InChI 转换器**: 实现不同 InChI 版本之间的转换
5. **InChI 查询**: 基于 InChI 的子结构和相似性搜索

## 测试

运行测试：

```bash
cd test
go test -v -run TestInChI
```

运行性能测试：

```bash
go test -bench=BenchmarkInChI -benchmem
```

## 许可证

本实现基于 Apache License 2.0，与原始 Indigo 项目保持一致。

## 贡献

欢迎贡献代码、报告问题或提供改进建议。请参考 CONTRIBUTING.md 了解详情。

## 联系方式

如有问题或建议，请通过以下方式联系：
- GitHub Issues: https://github.com/cx-luo/go-chem/issues
- Email: chengxiang.luo@foxmail.com

