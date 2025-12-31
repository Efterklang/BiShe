# AI Prompt: 智能养生店管理系统 (SmartSpa Admin)

### 1. 角色定义

你是一位顶级的全栈工程师，擅长构建高性能、轻量化的管理系统。你将使用 **Go (Gin)** 和 **Vue 3** 构建一套养生店后台管理系统。系统核心目标是：流程规范化、营销精准化（裂变）和决策科学化（LLM集成）。

---

### 2. 技术栈约束 (Mandatory Tech Stack)

- **后端:** Go 1.25+, Gin Framework, GORM (ORM).
- **数据库:** SQLite (单文件存储 `spa_management.db`)。
- **前端:** Vue 3 (Composition API / `<script setup>`), Vite, Pinia, Axios.
- **UI/样式:** Tailwind CSS + Daisy UI (主题：emerald)。
- **AI集成:** OpenAI 兼容 API 协议。
- **架构:** 前后端分离，RESTful API 风格，响应格式统一为：`{"code": 200, "data": {}, "msg": ""}`。

---

### 3. 数据库模型 (Database Models)

请严格按照以下逻辑定义 GORM 模型：

1.  **Member (会员):** ID, Name, Phone, Level (等级), YearlyTotalConsumption (年消费总值), InvitationCode (唯一邀请码), ReferrerID (推荐人ID).
2.  **Technician (技师):** ID, Name, Skills (JSON 存储技能标签), Status (0:空闲, 1:已预约, 2:请假), AverageRating.
3.  **ServiceItem (项目):** ID, Name, Duration (标准时长/分钟), Price, IsActive (上下架状态).
4.  **Appointment (预约单):** ID, MemberID, TechID, ServiceID, StartTime, EndTime (自动计算), Status (待服务/完成/候补/取消), OriginPrice, ActualPrice.
5.  **Schedule (排班):** ID, TechID, Date (日期), TimeSlots (已占用时间段 JSON), IsAvailable (布尔).
6.  **FissionLog (裂变日志):** ID, InviterID (邀请人), InviteeID (被邀请人), CommissionAmount (返利金额), OrderID (关联订单).

---

### 4. 核心功能逻辑描述 (Core Logic)

#### (1) 智能预约与候补调度

- **预约冲突检测:** 录入预约时，`EndTime = StartTime + Duration`。系统必须校验：
  1.  目标技师在 `Schedule` 中该日期是否可用。
  2.  `Appointment` 表中该技师在 `[StartTime, EndTime]` 范围内是否有重叠订单。
- **候补机制:** 若检测到冲突，允许创建状态为“候补中”的订单。当有订单“取消”时，触发钩子检查候补队列。

#### (2) 会员与裂变分销

- **等级自动晋升:** 每次订单完成后，更新会员的 `YearlyTotalConsumption`，并根据预设阈值自动触发 `Level` 升级。
- **社交分销:** 推荐人通过 `InvitationCode` 绑定新会员。新会员完成订单后，系统按 `ActualPrice` 的 10% 自动计算佣金并写入 `FissionLog`。

#### (3) 数据看板与 LLM 分析

- **可视化:** 前端展示：每日营收、会员增长趋势、项目销售排行（使用折线图/饼图）。
- **AI 经营决策:** 提供接口，后端聚合近 30 天营收数据、技师负载率及热门项目分布，构造结构化 Prompt 发送给 LLM，生成 Markdown 格式的经营分析周报（含项目上下架建议）。

---

### 5. 前端 UI/UX 要求

- 使用 **Daisy UI** 的 `emerald` 或 `cupcake` 主题。
- **Dashboard:** 采用响应式网格布局。
- **预约管理:** 需提供直观的时间轴或表格展示技师状态。
- **AI 交互:** 设置一个悬浮或侧边栏按钮“AI 顾问”，点击后以打字机效果展示 LLM 的经营建议。

---

### 6. 任务执行步骤 (Implementation Steps)

**Step 1 (后端初始化):** 建立 Go 项目目录，定义 `models` 结构体，初始化 SQLite 数据库及 GORM 自动迁移。
**Step 2 (核心 API):** 实现技师排班、服务项目管理、以及最核心的“带冲突检测的预约接口”。
**Step 3 (AI 服务):** 封装 LLM 调用逻辑，实现数据抽取与 Prompt 构建。
**Step 4 (前端构建):** 使用 Vite 初始化 Vue 3，配置 Axios 拦截器，使用 Tailwind 构建 Dashboard。
**Step 5 (业务闭环):** 实现会员注册、邀请码绑定、预约转候补的完整流程测试。

---

**请从 Step 1 开始，先给出项目的目录结构以及核心数据库模型（models）的代码实现。**
