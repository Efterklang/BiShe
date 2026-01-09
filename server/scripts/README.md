# 测试数据初始化脚本使用指南

## 概述

本目录包含用于为 SPA 管理系统生成测试数据的 SQL 脚本，帮助 Dashboard 显示真实的统计数据。

## 文件说明

- `seed_data.sql` - 主数据初始化脚本，包含会员、技师、服务项目、商品、预约、订单等测试数据

## 数据概览

脚本会插入以下测试数据：

- **30个会员** - 包含不同等级（VIP、Gold、Silver、Basic）和推荐关系
- **8个技师** - 包含不同技能和状态
- **10个服务项目** - 各种 SPA 和按摩服务
- **12个实体商品** - 精油、护理产品等
- **约150条预约记录** - 近30天的预约数据
- **约230条订单记录** - 服务订单 + 商品订单
- **约100条裂变记录** - 推荐佣金记录
- **约62条库存变动记录** - 入库、出库、销售记录

## 使用方法

### 方法一：使用 sqlite3 命令行工具

```bash
# 进入 server 目录
cd server

# 执行 SQL 脚本
sqlite3 spa_management.db < scripts/seed_data.sql
```

### 方法二：通过 Go 代码执行（推荐用于开发环境）

可以在 `server/internal/db/database.go` 中添加一个初始化函数：

```go
func SeedTestData(db *gorm.DB) error {
    sqlFile, err := os.ReadFile("scripts/seed_data.sql")
    if err != nil {
        return err
    }
    
    return db.Exec(string(sqlFile)).Error
}
```

然后在需要时调用该函数。

### 方法三：分步执行

如果只想插入特定类型的数据，可以手动复制 `seed_data.sql` 中对应的 INSERT 语句单独执行。

## 数据特点

### 1. 会员数据
- 包含推荐关系网络（referrer_id）
- 不同的会员等级和消费金额
- 覆盖近90天的注册时间

### 2. 预约和订单数据
- **时间分布**: 近30天均匀分布
- **状态**: 全部为已完成（completed）
- **价格**: 包含原价和实付价（部分有折扣）
- **每天平均5-6条预约记录**

### 3. 裂变数据
- 基于会员推荐关系自动生成
- 佣金比例为订单金额的 5%
- 与订单关联

### 4. 库存数据
- 包含商品销售导致的库存变动
- 包含批量入库记录
- 记录变动前后的库存数量

## Dashboard 统计说明

插入数据后，Dashboard 将显示以下真实统计：

### 1. 今日营收
- 统计今天所有已完成订单的实付金额总和
- 计算相比昨日的增长率

### 2. 新增会员
- 统计今天新注册的会员数量
- 显示本月累计新增会员数

### 3. 技师负载率
- 基于今日完成的预约数和技师总数计算
- 假设每个技师每天工作8小时

### 4. 近30天营收趋势
- 按日期统计每天的营收
- 柱状图显示趋势变化

### 5. 热门项目排行
- 统计各服务项目的订单数和总营收
- 显示 TOP 5 最受欢迎的服务

### 6. 裂变达人榜
- 统计每个会员邀请的人数
- 显示累计佣金收入
- 按邀请人数排序

## 清理数据

如果需要清理测试数据，取消 `seed_data.sql` 开头的注释：

```sql
-- 清理现有数据（可选，小心使用）
DELETE FROM orders;
DELETE FROM inventory_logs;
DELETE FROM appointments;
DELETE FROM fission_logs;
DELETE FROM physical_products;
DELETE FROM service_items;
DELETE FROM technicians;
DELETE FROM members;
```

**⚠️ 警告**: 这会删除所有相关数据，请谨慎操作！

## 自定义数据

### 修改数据量

如果需要更多或更少的数据，可以调整以下部分：

1. **会员数量**: 修改 `INSERT INTO members` 语句
2. **预约数量**: 修改 `LIMIT 120` 为其他数值
3. **商品订单数量**: 修改 `LIMIT 80` 为其他数值

### 修改时间范围

脚本使用 SQLite 的 `datetime('now', '-N days')` 函数生成时间。修改数字 N 可以调整时间范围。

例如，要生成近60天的数据：
```sql
datetime('now', '-60 days')
```

## 验证数据

执行脚本后，可以运行以下 SQL 验证数据是否正确插入：

```sql
SELECT '会员数量: ' || COUNT(*) FROM members;
SELECT '技师数量: ' || COUNT(*) FROM technicians;
SELECT '服务项目数量: ' || COUNT(*) FROM service_items;
SELECT '商品数量: ' || COUNT(*) FROM physical_products;
SELECT '预约数量: ' || COUNT(*) FROM appointments;
SELECT '订单数量: ' || COUNT(*) FROM orders;
SELECT '裂变记录数量: ' || COUNT(*) FROM fission_logs;
SELECT '库存变动记录数量: ' || COUNT(*) FROM inventory_logs;
```

或者直接访问 Dashboard 页面查看统计数据是否正常显示。

## 注意事项

1. **ID 冲突**: 脚本中使用了固定的 ID，如果数据库已有数据可能会冲突。建议在新数据库中执行。

2. **外键约束**: 数据插入顺序很重要，必须先插入被引用的表（如 members, technicians）再插入引用表（如 appointments）。

3. **随机数据**: 部分数据使用 SQLite 的 `RANDOM()` 函数生成，每次执行结果会有所不同。

4. **时区**: 脚本使用 SQLite 的 `datetime('now')` 函数，使用的是 UTC 时间。如需本地时间，可改为 `datetime('now', 'localtime')`。

## 故障排除

### 问题：执行脚本时出现外键约束错误

**解决方案**: 确保按顺序执行，或临时禁用外键约束：
```sql
PRAGMA foreign_keys = OFF;
-- 执行插入语句
PRAGMA foreign_keys = ON;
```

### 问题：Dashboard 显示数据为 0

**解决方案**: 
1. 检查数据是否成功插入（运行验证 SQL）
2. 检查后端 API 是否正常运行
3. 检查浏览器控制台是否有 API 调用错误
4. 确认后端路由已正确注册

### 问题：营收趋势图没有显示

**解决方案**:
1. 确保 orders 表中有 `status = 'completed'` 的记录
2. 确保订单的 `created_at` 在近30天内
3. 检查前端是否成功调用 `/api/dashboard/revenue-trend` API

## 生产环境注意

⚠️ **这些脚本仅用于开发和测试环境！**

在生产环境中：
- 不要执行清理数据的 DELETE 语句
- 不要使用固定 ID 插入数据
- 建议使用数据迁移工具而不是直接执行 SQL
- 做好数据备份

## 联系支持

如有问题，请查看：
- 主项目 README.md
- API 文档
- 或联系开发团队
