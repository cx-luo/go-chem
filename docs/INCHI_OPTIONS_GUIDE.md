# InChI Options 使用指南

## 概述

本指南说明如何正确使用 InChI 选项，基于官方 InChI API 规范 (`inchi_api.h`)。

## Options 格式规范

### 基本规则

根据 `inchi_api.h` 第 803 行的说明：

```
Valid options for GetINCHI:
(use - instead of / for O.S. other than MS Windows)
```

**格式要求**：
1. 选项之间用**空格**分隔
2. **Windows** 系统：每个选项前用 `/`
3. **Linux/其他** 系统：每个选项前用 `-`

### 示例

#### Windows
```go
generator.SetOptions("/FixedH /RecMet /AuxNone")
```

#### Linux/macOS
```go
generator.SetOptions("-FixedH -RecMet -AuxNone")
```

#### 自动格式化（推荐）

我们的实现会自动根据操作系统添加正确的前缀，因此您可以简单地写：

```go
generator.SetOptions("FixedH RecMet AuxNone")
```

代码会自动转换为：
- Windows: `/FixedH /RecMet /AuxNone`
- Linux: `-FixedH -RecMet -AuxNone`

## 可用选项

### 结构感知选项（兼容标准 InChI）

| 选项 | 说明 |
|------|------|
| `NEWPSOFF` | 不使用新的伪立体化学感知 |
| `DoNotAddH` | 不自动添加氢原子 |
| `SNon` | 包含省略的未定义/未知立体化学 |

### 立体化学解释选项（生成非标准 InChI）

| 选项 | 说明 |
|------|------|
| `SRel` | 相对立体化学 |
| `SRac` | 外消旋立体化学 |
| `SUCF` | 使用手性标志 |
| `ChiralFlagON` | 启用手性标志 |
| `ChiralFlagOFF` | 禁用手性标志 |

### InChI 创建选项（生成非标准 InChI）

| 选项 | 说明 |
|------|------|
| `SUU` | 允许不寻常的价态 |
| `SLUUD` | 包含不寻常的价态在断开层 |
| `FixedH` | 固定氢层（保留氢的位置） |
| `RecMet` | 重连接金属（处理金属配合物） |
| `KET` | 酮-烯醇互变异构 |
| `15T` | 1,5-氢转移互变异构 |

### 其他选项

| 选项 | 说明 |
|------|------|
| `AuxNone` | 省略辅助信息（默认包含） |
| `Wnumber` | 设置超时（秒），W0 = 无限制 |
| `WMnumber` | 设置超时（毫秒），WM0 = 无限制 |
| `WarnOnEmptyStructure` | 对空结构发出警告 |
| `SaveOpt` | 保存自定义 InChI 创建选项 |

## 标准 vs 非标准 InChI

### 标准 InChI

不使用任何 InChI 创建选项或立体化学修改选项时，生成标准 InChI：

```go
generator := molecule.NewInChIGeneratorCGO()
// 不设置选项，或只使用结构感知选项
generator.SetOptions("") // 或 "SNon" 或 "NEWPSOFF"
result, _ := generator.GenerateInChI(mol)
// 生成：InChI=1S/...
```

### 非标准 InChI

使用以下任一选项会生成非标准 InChI：
- `SUU`, `SLUUD`
- `FixedH`, `RecMet`
- `KET`, `15T`
- `SRel`, `SRac`, `SUCF`

```go
generator := molecule.NewInChIGeneratorCGO()
generator.SetOptions("FixedH RecMet")
result, _ := generator.GenerateInChI(mol)
// 生成：InChI=1/...（注意不是 1S）
```

## 使用示例

### 基本用法

```go
import "github.com/cx-luo/go-chem/molecule"

// 标准 InChI
gen := molecule.NewInChIGeneratorCGO()
result, err := gen.GenerateInChI(mol)
```

### 使用单个选项

```go
// 固定氢层
gen := molecule.NewInChIGeneratorCGO()
gen.SetOptions("FixedH")
result, err := gen.GenerateInChI(mol)
```

### 使用多个选项

```go
// 多个选项（空格分隔）
gen := molecule.NewInChIGeneratorCGO()
gen.SetOptions("FixedH RecMet AuxNone")
result, err := gen.GenerateInChI(mol)
```

### 显式指定前缀（不推荐，会自动处理）

```go
// Windows
gen.SetOptions("/FixedH /RecMet")

// Linux
gen.SetOptions("-FixedH -RecMet")

// 推荐：让代码自动处理
gen.SetOptions("FixedH RecMet")
```

## 常见用例

### 1. 标准 InChI（默认）

```go
gen := molecule.NewInChIGeneratorCGO()
// 不设置选项
result, _ := gen.GenerateInChI(mol)
fmt.Println(result.InChI) // InChI=1S/...
```

### 2. 包含辅助信息（默认已包含）

```go
gen := molecule.NewInChIGeneratorCGO()
result, _ := gen.GenerateInChI(mol)
fmt.Println(result.AuxInfo) // AuxInfo=1/...
```

### 3. 省略辅助信息

```go
gen := molecule.NewInChIGeneratorCGO()
gen.SetOptions("AuxNone")
result, _ := gen.GenerateInChI(mol)
fmt.Println(result.AuxInfo) // 空或很短
```

### 4. 固定氢位置

```go
gen := molecule.NewInChIGeneratorCGO()
gen.SetOptions("FixedH")
result, _ := gen.GenerateInChI(mol)
// 生成包含 /f 层的 InChI
```

### 5. 处理金属配合物

```go
gen := molecule.NewInChIGeneratorCGO()
gen.SetOptions("RecMet")
result, _ := gen.GenerateInChI(mol)
// 生成包含重连接金属信息的 InChI
```

### 6. 设置超时

```go
gen := molecule.NewInChIGeneratorCGO()
gen.SetOptions("W30") // 30 秒超时
result, _ := gen.GenerateInChI(mol)
```

### 7. 相对立体化学

```go
gen := molecule.NewInChIGeneratorCGO()
gen.SetOptions("SRel")
result, _ := gen.GenerateInChI(mol)
// 用于相对立体化学的分子
```

## 实现细节

### 自动前缀处理

```go
// normalizeInChIOptions 函数会：
// 1. 移除现有的 / 或 - 前缀
// 2. 根据 runtime.GOOS 添加正确的前缀
// 3. 用空格连接所有选项

func normalizeInChIOptions(options string) string {
    if options == "" {
        return ""
    }
    
    parts := strings.Fields(options)
    var normalized []string
    
    for _, part := range parts {
        // 移除现有前缀
        part = strings.TrimPrefix(part, "/")
        part = strings.TrimPrefix(part, "-")
        
        // 添加正确前缀
        var prefixed string
        if runtime.GOOS == "windows" {
            prefixed = "/" + part
        } else {
            prefixed = "-" + part
        }
        normalized = append(normalized, prefixed)
    }
    
    return strings.Join(normalized, " ")
}
```

## 测试

### 运行选项测试

```bash
# 测试选项格式化
CGO_ENABLED=1 go test -v ./test -run TestInChIOptionsNormalization

# 测试标准选项
CGO_ENABLED=1 go test -v ./test -run TestInChIStandardOptions

# 测试 AuxInfo 选项
CGO_ENABLED=1 go test -v ./test -run TestInChIOptionsWithAuxInfo
```

## 注意事项

### 1. 标准兼容性

⚠️ **重要**：使用以下选项会生成**非标准 InChI**：
- `FixedH`
- `RecMet`
- `SUU`
- `SLUUD`
- `KET`
- `15T`
- `SRel`
- `SRac`
- `SUCF`

非标准 InChI 可能无法被其他软件正确识别。

### 2. 选项大小写

选项名称是**大小写敏感**的：
- ✅ `FixedH`（正确）
- ❌ `fixedh`（错误）
- ❌ `FIXEDH`（错误）

### 3. 空格分隔

选项必须用**空格**分隔，不能用逗号或其他字符：
- ✅ `"FixedH RecMet"`（正确）
- ❌ `"FixedH,RecMet"`（错误）
- ❌ `"FixedH;RecMet"`（错误）

### 4. 前缀自动处理

我们的实现会自动处理前缀，您不需要手动添加：
- ✅ `"FixedH"`（推荐）
- ✅ `"/FixedH"`（Windows，会自动处理）
- ✅ `"-FixedH"`（Linux，会自动处理）

## 参考资料

### 官方文档

- **inchi_api.h**: 第 800-834 行
- **InChI Trust**: https://www.inchi-trust.org/
- **InChI Technical Manual**: https://www.inchi-trust.org/downloads/

### 相关代码

- `molecule/molecule_inchi_cgo.go`: CGO 实现
- `test/inchi_options_test.go`: 选项测试
- `examples/inchi_cgo_example.go`: 使用示例

## 常见问题

### Q: 我应该使用哪些选项？

**A**: 对于大多数情况，不使用任何选项（标准 InChI）是最好的选择。只有在特定需求时才使用特殊选项。

### Q: FixedH 和 RecMet 有什么区别？

**A**: 
- `FixedH`: 保留氢原子的位置信息，适用于互变异构体
- `RecMet`: 重连接金属原子，适用于金属配合物

### Q: 如何知道生成的是标准还是非标准 InChI？

**A**: 查看 InChI 字符串：
- 标准: `InChI=1S/...`（有 'S'）
- 非标准: `InChI=1/...`（没有 'S'）

### Q: 可以组合多个选项吗？

**A**: 可以，用空格分隔：
```go
gen.SetOptions("FixedH RecMet AuxNone")
```

### Q: Windows 和 Linux 的选项有区别吗？

**A**: 只有前缀不同（`/` vs `-`），但我们的实现会自动处理，您不需要关心。

## 总结

- ✅ 使用空格分隔选项
- ✅ 让代码自动处理前缀（推荐）
- ✅ 了解哪些选项会生成非标准 InChI
- ✅ 根据实际需求选择选项
- ✅ 大多数情况使用标准 InChI（不设置选项）

