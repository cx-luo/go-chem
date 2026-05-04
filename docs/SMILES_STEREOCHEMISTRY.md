# SMILES 立体化学支持文档

## 概述

本项目现已支持 SMILES 中的 `/` 和 `\` 符号，用于表示双键的顺反异构（cis/trans stereochemistry）。

## SMILES 立体化学规则

### 基础概念

在 SMILES 中，`/` 和 `\` 符号用于指定双键两侧取代基的空间排列：

- **`/`** - 表示键的"上"方向（BOND_UP）
- **`\`** - 表示键的"下"方向（BOND_DOWN）

### 判断顺反异构

立体化学配置取决于双键两侧的方向符号：

1. **反式（trans）配置**：双键两侧使用不同的方向符号
   - `C/C=C\C` - 反式-2-丁烯
   - `Cl/C=C\Cl` - 反式-1,2-二氯乙烯
   - `ClC\C=C\Cl` - 反式-1,2-二氯乙烯（从不同端看）

2. **顺式（cis）配置**：双键两侧使用相同的方向符号
   - `C/C=C/C` - 顺式-2-丁烯
   - `Cl/C=C/Cl` - 顺式-1,2-二氯乙烯
   - `ClC\C=C/Cl` - 顺式-1,2-二氯乙烯

### 重要原则

1. **方向是相对的**：方向符号指示键相对于即将出现的双键的方向
2. **必须成对出现**：完整的立体化学信息需要双键两侧都有方向标记
3. **与分支结合**：可以在分支中使用，如 `CC(/C)=C\C`
4. **环中使用**：也支持在环结构中使用方向符号

## 实现细节

### 数据结构

在 `Bond` 结构中存储方向信息：

```go
type Bond struct {
    Beg       int // 起始原子索引
    End       int // 结束原子索引
    Order     int // 键级（BOND_SINGLE, BOND_DOUBLE等）
    Direction int // 立体化学方向（BOND_UP, BOND_DOWN）
}
```

### 解析流程

1. **识别符号**：在 SMILES 解析过程中识别 `/` 和 `\` 字符
2. **设置 pending 方向**：将方向信息暂存在 `pendingDirection` 变量中
3. **应用到键**：创建键时，将 pending 方向应用到键的 `Direction` 字段
4. **重置状态**：应用后重置 `pendingDirection`

### 关键代码片段

```go
// 识别方向符号
if ch == '/' {
    pendingDirection = BOND_UP
    i++
    continue
}
if ch == '\\' {
    pendingDirection = BOND_DOWN
    i++
    continue
}

// 应用方向到键
bondIdx := m.AddBond(lastAtom, idx, order)
if pendingDirection != 0 {
    m.SetBondDirection(bondIdx, pendingDirection)
}
pendingDirection = 0
```

## 使用示例

```go
package main

import (
    "fmt"

    "github.com/cx-luo/go-chem/molecule"
)

func main() {
    loader := molecule.SmilesLoader{}

    // 解析反式-2-丁烯
    mol, err := loader.Parse("C/C=C\\C")
    if err != nil {
        panic(err)
    }

    // 检查键的方向信息
    for i, bond := range mol.Bonds {
        if bond.Direction != 0 {
            fmt.Printf("Bond %d: direction=%d\n", i, bond.Direction)
        }
    }
}
```

## 支持的 SMILES 示例

| SMILES | 描述 | 配置 |
|--------|------|------|
| `C/C=C\C` | 2-丁烯 | trans |
| `C/C=C/C` | 2-丁烯 | cis |
| `Cl/C=C\Cl` | 1,2-二氯乙烯 | trans |
| `Cl/C=C/Cl` | 1,2-二氯乙烯 | cis |
| `C/C=C\C1=CC=CC=C1` | 苯乙烯衍生物 | trans |
| `CC(/C)=C\C` | 带分支的烯烃 | trans |

## 测试

项目包含完整的立体化学测试套件：

```bash
go test -v ./test -run TestSmilesStereochemistry
```

测试覆盖：

- 反式和顺式配置
- 不同的起始方向（`/` 和 `\`）
- 与分支结合
- 环结构中的立体化学
- 方向信息的正确存储

## 技术限制

1. 当前实现存储了键的方向信息，但尚未实现完整的 E/Z 配置计算
2. 方向信息可用于后续的立体化学分析和匹配
3. 需要双键两侧都有方向标记才能完整表示立体化学

## 未来改进

- [ ] 实现自动 E/Z 配置判断
- [ ] 支持 SMILES 输出时保留立体化学信息
- [ ] 与 `MoleculeCisTrans` 系统集成
- [ ] 支持更复杂的立体化学情况（如 allenes）

## 参考资料

- [SMILES 规范](http://opensmiles.org/opensmiles.html)
- [Daylight SMILES 教程](https://www.daylight.com/dayhtml/doc/theory/theory.smiles.html)
