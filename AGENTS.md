# BisheApp Project Guide

## 项目简介
本项目是一个基于 Go (后端) 和 Vue 3 (前端) 的应用，主要包含以下核心功能：

- **数据看板（Dashboard）**: 提供应用的实时数据展示和AI分析功能。
- **服务预约 (Service Reservation)**：提供完整的服务预约流程。
- **一级裂变分销 (Level 1 Fission Distribution)**：包含用户邀请与分销机制。

## 代码编写规范 (Coding Standards)

### 1. 前端开发规范 (Frontend)
- **目录位置**: `client/`
- **包管理器**: 必须使用 **[bun](https://bun.sh/)** 进行依赖管理和脚本运行。
- **UI 组件库**:
  - 优先使用 **[DaisyUI](https://daisyui.com/)** 
  - Icon 优先使用 Lucide Icons
- **框架特性**: 全面使用 **Vue Composition API**
- **代码格式化**: 使用 **[Biome](https://biomejs.dev/)** 进行代码格式化和 Lint 检查。

### 2. 后端开发规范 (Backend)

- **目录位置**: `server/`
- **通用逻辑**:
  - 可复用的通用业务逻辑或工具函数，必须放置在 `server/pkg/util`
- **配置管理**:
  - 所有可配置的选项（如常量、环境配置等），必须放置在 `server/pkg/config`
- **数据库模型**:
  - 所有数据库模型（如结构体、表定义等），必须放置在 `server/internal/models/`，模型如果有需要复用的访问操作，应该放置在 `server/internal/repo/` 目录下。

