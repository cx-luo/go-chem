# 使用 GetStdINCHI 实现 InChI 生成

## 概述

`molecule_inchi.go` 现在支持使用 InChI API 的 `GetStdINCHI` 方法来生成 InChI，这是比 `MakeINCHIFromMolfileText` 更直接的方法。

## 实现方式

### 主要方法

1. **`ToInChI()`** - 主入口方法
   - 优先使用 `GetStdINCHI`（直接 API 调用）
   - 如果失败，回退到 `MakeINCHIFromMolfileText`（通过 Molfile）

2. **`toInChIUsingGetStdINCHI()`** - 使用 GetStdINCHI 的实现
   - 从 Indigo 分子提取所有原子信息
   - 构建 `inchi_Input` 结构
   - 调用 `GetStdINCHI` API

### 数据提取流程

```go
1. 获取原子数量
2. 为每个原子提取：
   - 元素符号 (Symbol)
   - 坐标 (x, y, z) - 当前设为 0，InChI 使用 0D 立体化学
   - 电荷 (Charge)
   - 自由基 (Radical) - 转换 Indigo 格式到 InChI 格式
   - 同位素 (Isotope)
   - 隐式氢 (Implicit Hydrogens)
   - 邻居原子和键类型 (通过遍历所有键)
3. 构建 inchi_Input 结构
4. 调用 GetStdINCHI
5. 提取结果
```

## 优势

### 使用 GetStdINCHI 的优势

1. **更直接**：不需要经过 Molfile 格式转换
2. **更高效**：减少中间格式转换的开销
3. **更灵活**：可以直接控制所有原子和键的属性
4. **更准确**：直接传递结构信息，避免格式转换中的信息丢失

### 回退机制

如果 `GetStdINCHI` 失败（例如由于复杂的结构或缺少某些信息），会自动回退到 `MakeINCHIFromMolfileText`，确保兼容性。

## 实现细节

### 类型转换

- **Indigo Radical → InChI Radical**: 
  - Indigo 使用 2,3,4 (对应 SINGLET, DOUBLET, TRIPLET)
  - InChI 使用 1,2,3
  - 转换：`inchi_radical = indigo_radical - 1`

- **坐标处理**:
  - 当前实现将坐标设为 (0, 0, 0)
  - InChI 会自动使用 0D 立体化学
  - 如果需要 2D/3D 坐标，可以从 Molfile 解析或使用 Layout

### 邻居提取

由于 Indigo 不提供直接的邻居迭代 API，实现通过以下方式获取：
1. 遍历所有键
2. 检查每个键的源和目标原子
3. 匹配到当前原子后，记录邻居和键类型

## 当前限制

1. **坐标**：当前设为 (0,0,0)，使用 0D 立体化学
2. **立体化学**：未实现 2D 立体化学标记（bond_stereo）
3. **0D 立体化学**：未实现 `inchi_Stereo0D` 结构
4. **性能**：邻居提取需要遍历所有键，对于大分子可能较慢

## 未来改进

1. **坐标提取**：从 Molfile 或使用 Layout 获取真实坐标
2. **立体化学支持**：实现 2D 和 0D 立体化学标记
3. **性能优化**：缓存原子索引映射，减少重复查找
4. **完整实现**：支持所有 InChI 选项和特性

## 使用示例

```go
// 基本使用（自动选择最佳方法）
mol, _ := molecule.LoadMoleculeFromString("CCO")
defer mol.Close()

inchi, err := mol.ToInChI()
if err != nil {
    log.Fatal(err)
}
fmt.Println("InChI:", inchi)

// 获取详细信息
result, err := mol.ToInChIWithInfo()
if err != nil {
    log.Fatal(err)
}
fmt.Println("InChI:", result.InChI)
fmt.Println("Key:", result.Key)
fmt.Println("Warning:", result.Warning)
```

## API 对比

| 方法 | 输入 | 优势 | 劣势 |
|------|------|------|------|
| `GetStdINCHI` | `inchi_Input` (原子/键结构) | 直接、高效、灵活 | 需要构建复杂结构 |
| `MakeINCHIFromMolfileText` | Molfile 字符串 | 简单、易用 | 需要格式转换 |

当前实现优先使用 `GetStdINCHI`，失败时回退到 `MakeINCHIFromMolfileText`，兼顾性能和兼容性。

