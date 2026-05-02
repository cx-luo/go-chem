# InChI 和 InChIKey 实现摘要

## 概述

本项目实现了基于 Indigo C++ 代码的 Go 语言 InChI 和 InChIKey 生成器。

## 主要功能

### 1. InChI 生成

InChI (International Chemical Identifier) 是 IUPAC 定义的化学物质标准标识符。

**示例:**
```go
result, err := molecule.GetInChIFromSMILES("CCO") // 乙醇
// InChI: InChI=1S/C2H6O/c1-2-3/h3H,2H2,1H3
```

### 2. InChIKey 生成

InChIKey 是 InChI 的固定长度哈希表示。

**格式:** `XXXXXXXXXXXXXX-YYYYYYYYY-ZZ`
- 前 14 字符: 连接性哈希
- 中间 9 字符: 立体化学哈希
- 后 2 字符: 版本标志

**示例:**
```go
key, err := molecule.GenerateInChIKey("InChI=1S/CH4/h1H4")
// InChIKey: VNWKTOKETHGBQD-UHFFFAOYSA-N
```

## 实现细节

### 核心算法

1. **化学式层** (Hill 系统)
   - C, H 优先，其余按字母排序
   - 参考: `molecule_inchi_layers.cpp`, `MainLayerFormula::printFormula`

2. **连接表层** (DFS 遍历)
   - 找度数最小顶点作起点
   - 深度优先搜索构建生成树
   - 按后代大小排序分支
   - 参考: `molecule_inchi_layers.cpp`, `printConnectionTable` (248-422 行)

3. **氢原子层** (范围压缩)
   - 按氢数量分组
   - 连续原子索引压缩为范围
   - 参考: `molecule_inchi_layers.cpp`, `HydrogensLayer::print` (468-528 行)

4. **立体化学层**
   - 顺反异构: 双键构型 (+/-)
   - 四面体: 手性中心奇偶性
   - 参考: `molecule_inchi_layers.cpp`, 各立体化学类

5. **InChIKey 生成**
   - SHA-256 哈希主结构和立体化学
   - Base-26 编码 (A-Z)
   - 参考: `inchi_wrapper.cpp`, `InChIKey` (705-730 行)

### C++ 到 Go 的对应

| C++ 方法 | Go 函数 | 功能 |
|---------|---------|------|
| `InchiWrapper::saveMoleculeIntoInchi` | `GenerateInChI` | 主 InChI 生成 |
| `InchiWrapper::InChIKey` | `GenerateInChIKey` | InChIKey 生成 |
| `MainLayerFormula::printFormula` | `generateFormulaLayer` | 化学式层 |
| `MainLayerConnections::printConnectionTable` | `generateConnectivityLayer` | 连接表层 |
| `HydrogensLayer::print` | `generateHydrogenLayer` | 氢原子层 |
| `CisTransStereochemistryLayer::print` | `generateCisTransLayer` | 顺反异构层 |
| `TetrahedralStereochemistryLayer::print` | `generateTetrahedralLayer` | 四面体立体化学层 |

## 使用方法

### 基本用法

```go
import "github.com/cx-luo/go-chem/molecule"

// 从 SMILES 生成 InChI
result, err := molecule.GetInChIFromSMILES("CCO")
if err != nil {
    log.Fatal(err)
}

fmt.Println("InChI:", result.InChI)
fmt.Println("InChIKey:", result.InChIKey)
```

### 高级用法

```go
// 创建自定义生成器
generator := molecule.NewInChIGenerator()
generator.SetPrefix("InChI=1S")

// 设置选项
generator.SetOptions(molecule.InChIOptions{
    FixedH:  true,  // 包含氢原子层
    RecMet:  false,
    AuxInfo: false,
    SNon:    false,
})

// 解析分子
loader := molecule.SmilesLoader{}
mol, _ := loader.Parse("CCO")

// 生成 InChI
result, err := generator.GenerateInChI(mol)
```

### 工具函数

```go
// 验证 InChI
valid := molecule.ValidateInChI("InChI=1S/CH4/h1H4")

// 比较 InChI
cmp := molecule.CompareInChI(inchi1, inchi2)

// Base64 编码/解码
encoded := molecule.Base64EncodeInChI(inchi)
decoded, _ := molecule.Base64DecodeInChI(encoded)
```

## 运行示例

```bash
# 运行示例程序
go run examples/inchi_example.go

# 运行测试
go test -v ./test -run TestInChI
```

## 改进点

相比原始实现，本次改进包括：

### 1. 更准确的连接表生成
- ✅ 基于 DFS 的算法
- ✅ 后代大小计算
- ✅ 分支排序
- ✅ 回边处理

### 2. 改进的氢原子层
- ✅ 范围压缩
- ✅ 按氢数量分组
- ✅ 格式化输出

### 3. 标准的 InChIKey 算法
- ✅ SHA-256 哈希
- ✅ Base-26 编码
- ✅ 正确的层分离
- ✅ 标志位处理

### 4. 完整的文档
- ✅ 算法详解
- ✅ C++ 对应关系
- ✅ 使用示例
- ✅ 中英文文档

## 当前状态

| 功能 | 状态 |
|------|------|
| 基本 InChI 生成 | ✅ 完成 |
| InChIKey 生成 | ✅ 完成 |
| 化学式层 | ✅ 完成 |
| 连接表层 | ✅ 完成 |
| 氢原子层 | ✅ 完成 |
| 立体化学层 | ✅ 完成 |
| 多组分支持 | ⏳ 部分完成 |
| InChI 解析 | ⏳ 待实现 |
| 规范化编号 | ⚠️ 简化版本 |

## 已知限制

1. **规范化编号**: 使用简化算法，完整版需要图自同构
2. **立体化学**: 奇偶性计算简化，需要 CIP 规则
3. **多组分**: 基本支持，可能需要改进
4. **性能**: 未优化，可能比 C++ 慢

## 未来计划

1. [ ] 实现完整的规范化编号算法
2. [ ] 改进立体化学处理 (CIP 规则)
3. [ ] 实现 InChI 解析
4. [ ] 添加更多测试用例
5. [ ] 性能优化

## 参考资料

- **Indigo C++ 源码**: indigo-core/molecule/src/
- **IUPAC InChI 规范**: https://www.inchi-trust.org/
- **详细文档**: 见 INCHI_IMPLEMENTATION.md

## 测试

测试覆盖以下场景：
- 简单分子 (甲烷、乙醇)
- 芳香环 (苯)
- 立体化学 (手性中心、双键)
- 多组分分子
- 边界情况

运行测试:
```bash
go test -v ./test -run TestInChI
```

## 贡献

欢迎提交 PR 改进实现！重点关注：
- 规范化编号算法
- 立体化学处理
- 性能优化
- 测试用例

## 许可证

Apache License 2.0 (与 Indigo 保持一致)
