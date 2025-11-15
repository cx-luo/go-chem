# 环境配置指南

本指南详细说明如何配置 go-indigo 的开发和运行环境。

## 目录

- [系统要求](#系统要求)
- [Windows 配置](#windows-配置)
- [Linux 配置](#linux-配置)
- [macOS 配置](#macos-配置)
- [CGO 配置详解](#cgo-配置详解)
- [IDE 配置](#ide-配置)
- [常见问题](#常见问题)

## 系统要求

### 必需组件

1. **Go 1.20 或更高版本**

   ```bash
   go version  # 检查版本
   ```

2. **CGO 支持**

   ```bash
   # 检查 CGO 是否启用
   go env CGO_ENABLED  # 应显示 "1"
   ```

3. **C 编译器**
   - Windows: MinGW-w64 或 MSVC
   - Linux: GCC
   - macOS: Clang (Xcode Command Line Tools)

### Indigo 库

项目已包含预编译的 Indigo 库：

```
3rd/
├── windows-x86_64/    # Windows 64位
├── windows-i386/      # Windows 32位
├── linux-x86_64/      # Linux 64位
├── linux-aarch64/     # Linux ARM64
├── darwin-x86_64/     # macOS Intel
└── darwin-aarch64/    # macOS Apple Silicon
```

## Windows 配置

### 1. 安装 Go

从 [golang.org](https://golang.org/dl/) 下载并安装 Go。

### 2. 安装 MinGW-w64

推荐使用 [MSYS2](https://www.msys2.org/)：

```powershell
# 安装 MSYS2 后，在 MSYS2 终端中执行：
pacman -S mingw-w64-x86_64-gcc
```

### 3. 设置环境变量

#### 临时设置（命令提示符）

```cmd
set CGO_ENABLED=1
set CGO_CFLAGS=-ID:/path/to/go-indigo/3rd
set CGO_LDFLAGS=-LD:/path/to/go-indigo/3rd/windows-x86_64
set PATH=%PATH%;D:/path/to/go-indigo/3rd/windows-x86_64
```

#### 临时设置（PowerShell）

```powershell
$env:CGO_ENABLED="1"
$env:CGO_CFLAGS="-ID:/path/to/go-indigo/3rd"
$env:CGO_LDFLAGS="-LD:/path/to/go-indigo/3rd/windows-x86_64"
$env:PATH="$env:PATH;D:/path/to/go-indigo/3rd/windows-x86_64"
```

#### 永久设置

1. 右键"此电脑" → 属性 → 高级系统设置 → 环境变量
2. 在"系统变量"中添加：
   - `CGO_ENABLED` = `1`
   - `CGO_CFLAGS` = `-ID:\path\to\go-indigo\3rd`
   - `CGO_LDFLAGS` = `-LD:\path\to\go-indigo\3rd\windows-x86_64`
3. 在 `PATH` 中添加：
   - `D:\path\to\go-indigo\3rd\windows-x86_64`

### 4. 验证配置

```cmd
cd go-indigo
go test ./test/molecule/... -v
```

### 常见问题（Windows）

#### 问题：找不到 DLL

**错误信息:**

```
exit status 0xc0000135
```

**解决方案:**

```cmd
# 确保 DLL 目录在 PATH 中
set PATH=%PATH%;D:/path/to/go-indigo/3rd/windows-x86_64

# 或者复制 DLL 到可执行文件目录
copy 3rd\windows-x86_64\*.dll .
```

#### 问题：CGO 编译失败

**错误信息:**

```
gcc: command not found
```

**解决方案:**

```cmd
# 确保 MinGW-w64 在 PATH 中
set PATH=%PATH%;C:\msys64\mingw64\bin
```

## Linux 配置

### 1. 安装 Go

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang

# 或从官网下载最新版本
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 2. 安装编译工具

```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# CentOS/RHEL
sudo yum groupinstall "Development Tools"

# Arch Linux
sudo pacman -S base-devel
```

### 3. 设置环境变量

#### 临时设置

```bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/3rd"
export CGO_LDFLAGS="-L$(pwd)/3rd/linux-x86_64"
export LD_LIBRARY_PATH="$(pwd)/3rd/linux-x86_64:$LD_LIBRARY_PATH"
```

#### 永久设置

编辑 `~/.bashrc` 或 `~/.zshrc`:

```bash
# go-indigo 配置
export CGO_ENABLED=1
export CGO_CFLAGS="-I$HOME/go-indigo/3rd"
export CGO_LDFLAGS="-L$HOME/go-indigo/3rd/linux-x86_64"
export LD_LIBRARY_PATH="$HOME/go-indigo/3rd/linux-x86_64:$LD_LIBRARY_PATH"
```

然后：

```bash
source ~/.bashrc  # 或 source ~/.zshrc
```

### 4. 验证配置

```bash
cd go-indigo
go test ./test/molecule/... -v
```

### 常见问题（Linux）

#### 问题：找不到共享库

**错误信息:**

```
error while loading shared libraries: libindigo.so: cannot open shared object file
```

**解决方案:**

```bash
# 方案 1: 设置 LD_LIBRARY_PATH
export LD_LIBRARY_PATH="$(pwd)/3rd/linux-x86_64:$LD_LIBRARY_PATH"

# 方案 2: 添加到系统库路径（需要 root）
sudo cp 3rd/linux-x86_64/*.so /usr/local/lib/
sudo ldconfig

# 方案 3: 创建符号链接
sudo ln -s $(pwd)/3rd/linux-x86_64/libindigo.so /usr/local/lib/
```

#### 问题：权限不足

```bash
# 确保库文件有执行权限
chmod +x 3rd/linux-x86_64/*.so
```

## macOS 配置

### 1. 安装 Xcode Command Line Tools

```bash
xcode-select --install
```

### 2. 安装 Go

```bash
# 使用 Homebrew
brew install go

# 或从官网下载
# https://golang.org/dl/
```

### 3. 设置环境变量

#### 临时设置

```bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/3rd"

# 根据你的Mac型号选择：
# Intel Mac:
export CGO_LDFLAGS="-L$(pwd)/3rd/darwin-x86_64"
export DYLD_LIBRARY_PATH="$(pwd)/3rd/darwin-x86_64:$DYLD_LIBRARY_PATH"

# Apple Silicon (M1/M2):
export CGO_LDFLAGS="-L$(pwd)/3rd/darwin-aarch64"
export DYLD_LIBRARY_PATH="$(pwd)/3rd/darwin-aarch64:$DYLD_LIBRARY_PATH"
```

#### 永久设置

编辑 `~/.zshrc`:

```bash
# go-indigo 配置
export CGO_ENABLED=1
export CGO_CFLAGS="-I$HOME/go-indigo/3rd"

# 根据架构选择一个：
export CGO_LDFLAGS="-L$HOME/go-indigo/3rd/darwin-aarch64"  # M1/M2
# export CGO_LDFLAGS="-L$HOME/go-indigo/3rd/darwin-x86_64"  # Intel

export DYLD_LIBRARY_PATH="$HOME/go-indigo/3rd/darwin-aarch64:$DYLD_LIBRARY_PATH"
```

### 4. 验证配置

```bash
cd go-indigo
go test ./test/molecule/... -v
```

### 常见问题（macOS）

#### 问题：安全策略阻止库加载

**错误信息:**

```
cannot be opened because the developer cannot be verified
```

**解决方案:**

```bash
# 移除隔离属性
xattr -d com.apple.quarantine 3rd/darwin-*/*.dylib

# 或在系统偏好设置中允许
```

#### 问题：架构不匹配

确保使用正确的库：

```bash
# 检查架构
uname -m
# x86_64 -> 使用 darwin-x86_64
# arm64  -> 使用 darwin-aarch64
```

## CGO 配置详解

### CGO 环境变量

#### CGO_ENABLED

```bash
# 启用 CGO（必需）
export CGO_ENABLED=1

# 禁用 CGO
export CGO_ENABLED=0
```

#### CGO_CFLAGS

指定 C 编译器标志，主要用于包含头文件：

```bash
# 单个路径
export CGO_CFLAGS="-I/path/to/include"

# 多个路径
export CGO_CFLAGS="-I/path/to/include1 -I/path/to/include2"

# 额外的编译选项
export CGO_CFLAGS="-I/path/to/include -O2 -Wall"
```

#### CGO_LDFLAGS

指定链接器标志，用于查找库文件：

```bash
# 库路径
export CGO_LDFLAGS="-L/path/to/lib"

# 链接特定库
export CGO_LDFLAGS="-L/path/to/lib -lindigo -lindigo-inchi"

# Linux rpath（运行时库路径）
export CGO_LDFLAGS="-L/path/to/lib -lindigo -Wl,-rpath,/path/to/lib"
```

### 代码中的 CGO 指令

在 Go 代码中使用 `#cgo` 指令：

```go
/*
#cgo CFLAGS: -I${SRCDIR}/../3rd

// 平台特定配置
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../3rd/windows-x86_64 -lindigo
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../3rd/linux-x86_64 -lindigo -Wl,-rpath,${SRCDIR}/../3rd/linux-x86_64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../3rd/darwin-aarch64 -lindigo

#include "indigo.h"
*/
import "C"
```

### 构建标签

使用构建标签控制不同平台的编译：

```go
//go:build windows
// +build windows

package mypackage

// Windows 特定代码
```

## IDE 配置

### VS Code

#### 1. 安装插件

- Go (官方)
- C/C++ (用于 CGO)

#### 2. 配置 settings.json

```json
{
    "go.toolsEnvVars": {
        "CGO_ENABLED": "1",
        "CGO_CFLAGS": "-I${workspaceFolder}/3rd",
        "CGO_LDFLAGS": "-L${workspaceFolder}/3rd/windows-x86_64"
    },
    "go.testEnvVars": {
        "PATH": "${env:PATH};${workspaceFolder}/3rd/windows-x86_64"
    }
}
```

### GoLand

#### 1. 配置 Go Modules

File → Settings → Go → Go Modules

- 启用 Go Modules
- 设置 Environment: `CGO_ENABLED=1`

#### 2. 配置运行配置

Run → Edit Configurations → Go Build

- Environment:

  ```
  CGO_ENABLED=1;
  CGO_CFLAGS=-I$ProjectFileDir$/3rd;
  CGO_LDFLAGS=-L$ProjectFileDir$/3rd/windows-x86_64;
  PATH=$PATH$;$ProjectFileDir$/3rd/windows-x86_64
  ```

### Vim/Neovim

使用 vim-go 插件，在 `.bashrc` 中设置环境变量即可。

## 测试配置

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./test/molecule/...

# 详细输出
go test -v ./test/molecule/...

# 指定测试
go test -v -run TestLoadMoleculeFromString ./test/molecule/
```

### 测试环境变量

```bash
# 设置测试超时
go test -timeout 5m ./...

# 并行测试
go test -parallel 4 ./...

# 竞态检测
go test -race ./...
```

## 部署配置

### 构建可执行文件

```bash
# 普通构建
go build -o myapp

# 优化构建（减小体积）
go build -ldflags="-s -w" -o myapp
```

### 分发应用

确保包含必需的库文件：

```bash
# Windows
myapp.exe
indigo.dll
indigo-inchi.dll
indigo-renderer.dll
msvcp140.dll
vcruntime140.dll

# Linux
myapp
libindigo.so
libindigo-inchi.so
libindigo-renderer.so

# macOS
myapp
libindigo.dylib
libindigo-inchi.dylib
libindigo-renderer.dylib
```

## 常见问题汇总

### Q: CGO_ENABLED=0 错误

**问题:** 使用交叉编译时 CGO 被禁用。

**解决:**

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build
```

### Q: 头文件找不到

**问题:** `fatal error: indigo.h: No such file or directory`

**解决:**

```bash
export CGO_CFLAGS="-I$(pwd)/3rd"
```

### Q: 链接库找不到

**问题:** `undefined reference to 'indigoCreateMolecule'`

**解决:**

```bash
# 检查库文件是否存在
ls 3rd/linux-x86_64/libindigo.so

# 设置链接路径
export CGO_LDFLAGS="-L$(pwd)/3rd/linux-x86_64 -lindigo"
```

### Q: 运行时找不到动态库

**问题:**

- Windows: `exit status 0xc0000135`
- Linux: `error while loading shared libraries`
- macOS: `Library not loaded`

**解决:**

```bash
# Windows
set PATH=%PATH%;path\to\3rd\windows-x86_64

# Linux
export LD_LIBRARY_PATH=path/to/3rd/linux-x86_64:$LD_LIBRARY_PATH

# macOS
export DYLD_LIBRARY_PATH=path/to/3rd/darwin-aarch64:$DYLD_LIBRARY_PATH
```

## 获取帮助

如果配置遇到问题：

1. 检查 [常见问题](FAQ.md)
2. 在 GitHub 创建 Issue
3. 发送邮件至 <chengxiang.luo@foxmail.com>

---

💡 **提示**: 推荐使用脚本自动配置环境变量，避免每次手动设置！
