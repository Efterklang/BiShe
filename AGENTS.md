# BisheApp Project Guide

## 项目简介
本项目是一个基于 Go (后端) 和 Vue 3 (前端) 的应用，主要包含以下核心功能：

- **数据看板（Dashboard）**: 提供应用的实时数据展示和AI分析功能。
- **服务预约 (Service Reservation)**：提供完整的服务预约流程。
- **一级裂变分销 (Level 1 Fission Distribution)**：包含用户邀请与分销机制。
- **订单管理与营销数据**：以 `orders` 作为统一订单事实表，支持服务/商品订单与营销指标提取。

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

## 项目注意事项

### 订单与佣金
- 统一订单表：`orders` 为订单事实表（`order_type=service|physical`），金额字段使用 decimal 存储，业务计算建议走 `server/pkg/util` 的金额工具。
- 商品销售下单：库存出库 `action_type=sale` 时必须携带 `member_id` 才会创建订单（否则会返回 400）。
- 裂变佣金：`fission_logs` 不再维护 `order_id` 字段；看板/列表的佣金展示以 `orders.commission_amount` 与 `orders.inviter_id` 为准。

### Dashboard 数据口径
- 近30天营收趋势、热门项目排行、实体商品销售概览均从 `orders` 聚合统计（服务订单/商品订单分别汇总）。
- 营销数据接口：`GET /api/dashboard/marketing` 支持按时间范围、订单类型、会员等级筛选，并返回汇总与日/周/月聚合序列。

### 测试与开发
- 后端测试：在 `server/` 下运行 `go test ./...`。
- 前端依赖与脚本：必须使用 bun（见前端规范）。
