# 立体化学层实现计划

## 概述

本文档描述了 InChI 立体化学层的实现计划。立体化学层是 InChI 规范的重要组成部分，用于表示分子的三维空间构型。

## 当前状态

### 已实现 ✅
- 基本的 InChI 生成框架
- 分子式层（Hill 系统）
- 连接层（原子编号和连接）
- 氢原子层
- InChIKey 生成

### 待实现 ⏳
- 双键立体化学层 (/b)
- 四面体立体化学层 (/t)
- 对映异构体层 (/m)

## 实现计划

### 1. 双键立体化学层 (/b)

#### 功能描述
表示双键的顺反异构（cis/trans 或 E/Z）。

#### 算法步骤

1. **识别具有立体化学的双键**
   ```go
   for e := mol.edgeBegin(); e != mol.edgeEnd(); e = mol.edgeNext(e) {
       if mol.getBondOrder(e) != BOND_DOUBLE {
           continue
       }
       if mol.cis_trans.getParity(e) == 0 {
           continue // No stereochemistry
       }
       // Process this double bond
   }
   ```

2. **获取取代基**
   ```go
   var subst [4]int
   mol.getSubstituents_All(e, subst[:])
   // subst[0], subst[1]: substituents on one end
   // subst[2], subst[3]: substituents on the other end
   ```

3. **确定构型**
   ```go
   parity := mol.cis_trans.getParity(e)
   if parity == MoleculeCisTrans.CIS {
       // Same side: use '-'
   } else if parity == MoleculeCisTrans.TRANS {
       // Opposite side: use '+'
   }
   ```

4. **编码为 InChI 格式**
   - 格式：`bond_number+` 或 `bond_number-`
   - 示例：`/b3-4+` 表示键 3-4 是 trans

#### 数据结构需求

当前 `Molecule` 结构可能需要以下字段：
- `cis_trans`: 顺反异构管理器
- 方法：`getParity(bondIdx)`, `getSubstituents_All(bondIdx, subst[])`

#### 参考资料
- Indigo: `molecule_inchi.cpp`, `_printInChIComponentCisTrans()`
- IUPAC InChI Technical Manual, Section 3.4: "Double Bond Stereochemistry"
- 化学规则：E/Z 命名法（Cahn-Ingold-Prelog 优先规则）

### 2. 四面体立体化学层 (/t)

#### 功能描述
表示四面体手性中心的立体构型（R/S）。

#### 算法步骤

1. **识别手性中心**
   ```go
   for i := mol.stereocenters.begin(); i != mol.stereocenters.end(); 
       i = mol.stereocenters.next(i) {
       
       if !mol.stereocenters.isTetrahydral(i) {
           continue
       }
       
       atomIdx := mol.stereocenters.getAtomIndex(i)
       // Process this stereocenter
   }
   ```

2. **获取金字塔构型**
   ```go
   var centerType, group int
   var pyramid [4]int
   mol.stereocenters.get(atomIdx, &centerType, &group, pyramid[:])
   
   // pyramid[0-3] are the four substituents
   // pyramid[3] == -1 means only 3 explicit substituents
   ```

3. **计算奇偶性**
   - 需要根据取代基的优先级顺序计算
   - 使用规范化原子编号
   - 确定是顺时针还是逆时针

4. **编码为 InChI 格式**
   - 格式：`atom_number+` 或 `atom_number-`
   - 示例：`/t2-,3+` 表示原子 2 是 '-', 原子 3 是 '+'

#### 优先级规则（Cahn-Ingold-Prelog）

需要实现以下优先级比较：
1. 原子序数：更高的优先
2. 同位素质量：更重的优先
3. 递归比较：考虑连接的原子
4. 双键/三键：Z > E

#### 数据结构需求

当前 `Molecule` 结构需要：
- `stereocenters`: 立体中心管理器
- 方法：`isTetrahydral()`, `get()`, `getType()`

#### 参考资料
- Indigo: `molecule_inchi_component.cpp`, tetrahedral stereochemistry handling
- IUPAC InChI Technical Manual, Section 3.5: "Tetrahedral Stereochemistry"
- R/S 命名法：Cahn-Ingold-Prelog 优先规则
- 金字塔表示法：四个取代基的空间排列

### 3. 对映异构体层 (/m)

#### 功能描述
指示分子的对映异构体类型。

#### 可能的值
- `0`: 绝对立体化学（absolute）
- `1`: 相对立体化学/外消旋混合物（relative/racemic）

#### 算法步骤

1. **分析所有立体中心**
   ```go
   allAbsolute := true
   hasRelative := false
   
   for i := mol.stereocenters.begin(); i != mol.stereocenters.end(); 
       i = mol.stereocenters.next(i) {
       
       centerType := mol.stereocenters.getType(i)
       
       if centerType == MoleculeStereocenters.ATOM_AND {
           hasRelative = true
       } else if centerType == MoleculeStereocenters.ATOM_OR {
           // Special handling needed
       }
   }
   ```

2. **确定值**
   - 如果所有立体中心都是 `ATOM_ABS`，返回 "0"
   - 如果任何立体中心是 `ATOM_AND`，返回 "1"
   - `ATOM_OR` 和 `ATOM_ANY` 需要特殊处理

#### 参考资料
- Indigo: `molecule_inchi.cpp`, `generateEnantiomerLayer()`
- IUPAC InChI Technical Manual, Section 3.6: "Enantiomer Information"

## 实现优先级

### 高优先级 🔴
1. **四面体立体化学层** - 最常见的立体化学类型
   - 实现手性中心识别
   - 实现 Cahn-Ingold-Prelog 优先规则
   - 实现奇偶性计算

### 中优先级 🟡
2. **双键立体化学层** - 重要但复杂度较低
   - 实现顺反异构识别
   - 实现取代基分析

3. **对映异构体层** - 依赖前两者
   - 分析所有立体中心类型
   - 确定绝对/相对构型

## 技术挑战

### 1. 优先级规则实现

Cahn-Ingold-Prelog 规则的完整实现需要：
- 递归原子比较
- 同位素处理
- 多重键特殊处理
- 环系统处理

**建议方法**：
- 创建 `cipPriority` 包
- 实现递归比较算法
- 缓存计算结果

### 2. 规范化编号

立体化学层需要使用规范化的原子编号：
- 当前实现使用简化的排序
- 需要完整的图同构算法

**建议方法**：
- 先使用简化版本
- 后续升级为完整的规范化算法

### 3. 几何计算

需要从 3D 坐标或 2D 键方向计算立体构型：
- 向量叉积
- 四面体体积
- 平面法向量

**建议方法**：
- 创建 `geometry` 包
- 实现必要的向量运算

## 测试策略

### 单元测试

1. **CIP 优先级测试**
   ```go
   TestCIPPriority_SimpleAtoms
   TestCIPPriority_Isotopes
   TestCIPPriority_Recursive
   TestCIPPriority_DoubleBonds
   ```

2. **立体中心检测测试**
   ```go
   TestTetrahedralDetection_Simple
   TestTetrahedralDetection_Complex
   TestTetrahedralDetection_From3D
   ```

3. **InChI 生成测试**
   ```go
   TestInChI_WithChiralCenter
   TestInChI_WithCisTrans
   TestInChI_ComplexStereochemistry
   ```

### 集成测试

使用已知的 InChI 进行验证：
- (R)-乳酸
- (S)-丙氨酸
- 顺-2-丁烯
- 反-2-丁烯

### 参考数据源

- PubChem: 提供已知分子的 InChI
- ChemSpider: 验证立体化学
- RDKit: 作为参考实现

## 估算工作量

### Phase 1: 基础设施（2-3 周）
- [ ] 实现 Cahn-Ingold-Prelog 优先规则
- [ ] 实现几何计算工具
- [ ] 改进规范化算法

### Phase 2: 四面体立体化学（2 周）
- [ ] 手性中心识别
- [ ] 奇偶性计算
- [ ] InChI 编码

### Phase 3: 双键立体化学（1-2 周）
- [ ] 顺反异构识别
- [ ] 取代基分析
- [ ] InChI 编码

### Phase 4: 集成和测试（1 周）
- [ ] 对映异构体层
- [ ] 完整测试套件
- [ ] 与参考实现对比

**总计**: 约 6-8 周的开发时间

## 替代方案

如果完整实现工作量过大，可以考虑：

### 方案 A: 使用 CGO 调用官方 InChI 库
**优点**：
- 完全兼容
- 已经过充分测试
- 包含所有特性

**缺点**：
- CGO 依赖
- 跨平台编译复杂
- 性能可能略低

### 方案 B: 实现简化版本
**优点**：
- 快速实现
- 满足大多数使用场景

**缺点**：
- 不完全兼容
- 某些边缘情况可能不正确

**建议**：
- 优先实现核心功能（方案 B）
- 为高级用户提供 CGO 选项（方案 A）
- 在文档中明确说明限制

## 参考资料

### 官方文档
1. IUPAC InChI Technical Manual
   - https://www.inchi-trust.org/downloads/
   
2. InChI Algorithm Description
   - https://www.inchi-trust.org/technical-faq/

### 学术论文
1. Goodman et al. (2012): "InChI version 1, three years on"
2. Cahn, Ingold, Prelog (1966): "Specification of Molecular Chirality"

### 开源实现
1. Indigo Toolkit (C++)
   - https://github.com/epam/Indigo
   
2. RDKit (C++/Python)
   - https://www.rdkit.org/

3. OpenBabel (C++)
   - http://openbabel.org/

## 结论

立体化学层的实现是一个复杂但重要的功能。建议：

1. **短期**（1-2个月）：
   - 实现 CIP 优先规则
   - 实现四面体立体化学层
   - 添加基本测试

2. **中期**（3-6个月）：
   - 实现双键立体化学层
   - 完善测试套件
   - 性能优化

3. **长期**（6-12个月）：
   - 完整的规范化算法
   - 复杂分子处理
   - 与官方实现对标

对于大多数应用场景，当前的基础 InChI 功能已经足够。立体化学层可以作为高级功能逐步添加。

---

**更新日期**: 2025年11月1日
**状态**: 计划中
**优先级**: 中等

