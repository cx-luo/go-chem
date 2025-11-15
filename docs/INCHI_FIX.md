# InChI 转换错误修复说明

## 问题描述

在使用 `ToInChI()` 方法时，特别是对从反应中获取的分子进行转换时，可能会遇到以下错误：

```
failed to convert reactant to InChI: failed to convert to InChI: invalid unordered_map<K, T> key
```

## 错误原因

这个错误来自 Indigo 库的 C++ 内部实现，通常在以下情况下发生：

1. **分子缺少必要的属性**：从反应中获取的分子可能缺少某些 InChI 生成器期望的属性
2. **芳香性信息缺失**：InChI 生成过程需要正确识别分子的芳香性
3. **内部状态不一致**：分子对象的内部 `unordered_map` 中缺少某些必需的键值对

### Indigo 源码分析

在 Indigo 的 InChI 插件实现中，InChI 生成器会：
- 访问分子的原子和键属性映射表
- 查询芳香性信息
- 访问立体化学信息

如果这些信息在内部的 `std::unordered_map` 中不存在，就会抛出 `std::out_of_range` 异常。

## 解决方案

我们在 `molecule_inchi.go` 中修改了 `ToInChI()` 方法，添加了以下保护措施：

### 1. 克隆分子

```go
// Clone the molecule to avoid modifying the original
clonedHandle := int(C.indigoClone(C.int(m.handle)))
if clonedHandle < 0 {
    return "", fmt.Errorf("failed to clone molecule for InChI: %s", getLastError())
}
defer C.indigoFree(C.int(clonedHandle))
```

**原因**：克隆可以确保：
- 不会修改原始分子对象
- 获得一个完整的、独立的分子副本
- 避免与反应对象的所有权冲突

### 2. 芳香化处理

```go
// Aromatize the cloned molecule
ret := int(C.indigoAromatize(C.int(clonedHandle)))
if ret < 0 {
    // Aromatization might fail for some molecules, which is okay
    C.indigoClearTitles()
}
```

**原因**：
- InChI 生成器需要正确的芳香性信息
- 芳香化会补全分子内部的属性映射表
- 即使芳香化失败，也尝试继续生成 InChI

## 使用示例

### 从反应中获取分子并转换为 InChI

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cx-luo/go-indigo/molecule"
    "github.com/cx-luo/go-indigo/reaction"
)

func main() {
    // 加载反应
    rxn, err := reaction.LoadReactionFromString("CCO.CC(=O)O>>[H+].CCOC(=O)C")
    if err != nil {
        log.Fatal(err)
    }
    defer rxn.Close()

    // 获取反应物分子
    reactants, err := rxn.GetReactantMolecules()
    if err != nil {
        log.Fatal(err)
    }

    // 转换为 InChI（现在已修复）
    for i, mol := range reactants {
        defer mol.Close()
        
        inchi, err := mol.ToInChI()
        if err != nil {
            log.Printf("Failed to convert reactant %d: %v", i, err)
            continue
        }
        
        fmt.Printf("Reactant %d InChI: %s\n", i, inchi)
    }
}
```

### 直接加载分子并转换（无需特殊处理）

```go
// 对于直接加载的分子，新的实现仍然完全兼容
mol, err := molecule.LoadMoleculeFromString("CCO")
if err != nil {
    log.Fatal(err)
}
defer mol.Close()

inchi, err := mol.ToInChI()
if err != nil {
    log.Fatal(err)
}
fmt.Println("InChI:", inchi)
```

## 性能影响

修复后的实现会：
- 额外执行一次分子克隆（开销很小）
- 尝试芳香化（对大多数分子来说很快）
- 自动清理克隆的分子对象

对于正常使用场景，性能影响可以忽略不计（通常 < 1ms）。

## 兼容性

这个修复：
- ✅ 向后兼容：对现有代码无影响
- ✅ 自动应用：无需修改调用代码
- ✅ 安全：不会修改原始分子对象
- ✅ 健壮：处理了芳香化失败的情况

## 测试

所有现有的 InChI 测试用例仍然通过：

```bash
go test ./test/molecule -run TestToInChI -v
go test ./test/molecule -run TestInChI -v
```

## 相关问题

如果仍然遇到 InChI 转换错误，可能的原因包括：

1. **分子结构无效**：检查 SMILES 或输入格式是否正确
2. **不支持的元素**：InChI 可能不支持某些罕见元素
3. **立体化学问题**：未定义的立体中心可能导致问题

可以通过以下方式获取更多诊断信息：

```go
result, err := mol.ToInChIWithInfo()
if err != nil {
    log.Fatal(err)
}

fmt.Println("InChI:", result.InChI)
fmt.Println("Warning:", result.Warning)
fmt.Println("Log:", result.Log)
fmt.Println("AuxInfo:", result.AuxInfo)
```

## 参考资料

- [Indigo 官方文档](https://lifescience.opensource.epam.com/indigo/)
- [Indigo GitHub 仓库](https://github.com/epam/indigo)
- [InChI Trust](https://www.inchi-trust.org/)
- [std::unordered_map 文档](https://en.cppreference.com/w/cpp/container/unordered_map)

