# Molecule Access Improvements

## 概述

针对用户提出的问题，我们对 reaction 包中的分子访问方法进行了全面优化，使得从反应中获取的分子对象更易于使用。

## 问题分析

### 原有问题

1. `GetReactant(index)` 等方法返回 `int` 类型的 handle，而不是 `*Molecule` 对象
2. 使用者需要手动将 handle 转换为 Molecule 对象，增加了使用难度
3. 缺少批量获取分子的便捷方法

### 解决方案

我们添加了以下改进：

1. **molecule 包**: 添加了公开的 `NewMoleculeFromHandle()` 方法
2. **reaction 包**: 添加了返回 `*Molecule` 的新方法
3. **reaction 包**: 添加了批量获取的便捷方法

## 新增方法

### molecule 包

#### NewMoleculeFromHandle

```go
func NewMoleculeFromHandle(handle int) (*Molecule, error)
```

从 Indigo handle 创建 Molecule 对象。这对于从反应或其他来源获取的 handle 特别有用。

**注意**: Molecule 对象会接管 handle 的所有权，并在 Close() 时释放它。

### reaction 包

#### 单个分子访问方法

```go
func (r *Reaction) GetReactantMolecule(index int) (*Molecule, error)
func (r *Reaction) GetProductMolecule(index int) (*Molecule, error)
func (r *Reaction) GetCatalystMolecule(index int) (*Molecule, error)
```

通过索引获取单个反应物、产物或催化剂分子，返回完整的 `*Molecule` 对象。

**使用示例**:

```go
// 获取第一个反应物
reactant, err := rxn.GetReactantMolecule(0)
if err != nil {
    log.Fatal(err)
}
defer reactant.Close()

// 现在可以使用所有 Molecule 方法
smiles, _ := reactant.ToSmiles()
formula, _ := reactant.GrossFormula()
mass, _ := reactant.MolecularWeight()
```

#### 批量访问方法

```go
func (r *Reaction) GetAllReactants() ([]*Molecule, error)
func (r *Reaction) GetAllProducts() ([]*Molecule, error)
func (r *Reaction) GetAllCatalysts() ([]*Molecule, error)
```

一次性获取所有反应物、产物或催化剂，返回 `[]*Molecule` 切片。

**使用示例**:

```go
// 获取所有反应物
reactants, err := rxn.GetAllReactants()
if err != nil {
    log.Fatal(err)
}

// 处理每个反应物
for i, mol := range reactants {
    smiles, _ := mol.ToSmiles()
    fmt.Printf("Reactant %d: %s\n", i, smiles)
    mol.Close() // 记得关闭！
}
```

## 技术细节

### Handle 克隆

所有新方法都会自动克隆分子的 handle，避免所有权冲突：

```go
clonedHandle := int(C.indigoClone(C.int(molHandle)))
```

这确保了：
1. 原始反应对象不受影响
2. 返回的 Molecule 对象可以独立管理
3. 避免重复释放导致的内存问题

### 资源管理

返回的 `*Molecule` 对象使用 `runtime.SetFinalizer` 自动管理内存，但仍然建议显式调用 `Close()`:

```go
mol, _ := rxn.GetReactantMolecule(0)
defer mol.Close() // 推荐做法
```

## 向后兼容性

### 保留的原有方法

我们保留了原有的返回 `int` handle 的方法，以保持向后兼容：

```go
func (r *Reaction) GetReactant(index int) (int, error)
func (r *Reaction) GetProduct(index int) (int, error)
func (r *Reaction) GetCatalyst(index int) (int, error)
```

### 推荐使用新方法

虽然原有方法仍然可用，但我们强烈推荐使用新的返回 `*Molecule` 的方法，因为它们：
- 更易于使用
- 类型安全
- 自动管理资源
- 提供完整的 Molecule API

## 使用场景

### 1. 分析反应质量平衡

```go
reactants, _ := rxn.GetAllReactants()
products, _ := rxn.GetAllProducts()

var reactantMass, productMass float32
for _, mol := range reactants {
    mass, _ := mol.MolecularWeight()
    reactantMass += mass
    mol.Close()
}

for _, mol := range products {
    mass, _ := mol.MolecularWeight()
    productMass += mass
    mol.Close()
}

fmt.Printf("Mass balance: %.2f → %.2f\n", reactantMass, productMass)
```

### 2. 处理反应中的分子

```go
// 获取所有反应物并进行芳香化
reactants, _ := rxn.GetAllReactants()
for _, mol := range reactants {
    mol.Aromatize()
    smiles, _ := mol.ToCanonicalSmiles()
    fmt.Println(smiles)
    mol.Close()
}
```

### 3. 生成反应报告

```go
fmt.Println("Reaction Analysis:")
fmt.Println("\nReactants:")

reactants, _ := rxn.GetAllReactants()
for i, mol := range reactants {
    smiles, _ := mol.ToSmiles()
    formula, _ := mol.GrossFormula()
    mass, _ := mol.MolecularWeight()
    atoms, _ := mol.CountAtoms()
    
    fmt.Printf("%d. %s\n", i+1, smiles)
    fmt.Printf("   Formula: %s, MW: %.2f, Atoms: %d\n", formula, mass, atoms)
    
    mol.Close()
}
```

## 测试

新方法包含完整的测试覆盖：

- `test/reaction/reaction_molecules_test.go` - 单元测试
- `examples/reaction/reaction_molecules.go` - 使用示例

运行测试：

```bash
go test ./test/reaction/reaction_molecules_test.go -v
```

运行示例：

```bash
cd examples/reaction
go run reaction_molecules.go
```

## 性能考虑

### Handle 克隆的开销

所有方法都会克隆 handle，这会有轻微的性能开销。如果需要极致性能且了解 Indigo API，可以使用原有的返回 handle 的方法。

### 批量访问的优势

`GetAllReactants()` 等方法比多次调用 `GetReactantMolecule()` 更高效，因为它们只创建一次迭代器。

## 迁移指南

### 从旧代码迁移

**旧方式**:
```go
handle, _ := rxn.GetReactant(0)
// 需要手动处理 handle
```

**新方式**:
```go
mol, _ := rxn.GetReactantMolecule(0)
defer mol.Close()
// 直接使用 Molecule 方法
smiles, _ := mol.ToSmiles()
```

## 总结

这些改进使得从反应中访问和操作分子变得更加直观和安全。新的 API 设计遵循 Go 的最佳实践，提供了更好的类型安全性和资源管理。

## 相关文件

- `molecule/molecule.go` - 添加了 `NewMoleculeFromHandle()`
- `reaction/reaction_helpers.go` - 添加了所有新的分子访问方法
- `examples/reaction/reaction_molecules.go` - 完整的使用示例
- `examples/reaction/README.md` - 更新了文档
- `test/reaction/reaction_molecules_test.go` - 完整的测试套件

