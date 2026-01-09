-- ============================================================
-- 测试数据插入脚本（改进版）
-- 用于为 Dashboard 提供真实统计数据
-- ============================================================

-- 禁用外键约束，以便安全删除数据
PRAGMA foreign_keys = OFF;

-- ============================================================
-- 清理现有测试数据
-- ============================================================
DELETE FROM inventory_logs;
DELETE FROM fission_logs;
DELETE FROM orders;
DELETE FROM appointments;
DELETE FROM schedules;
DELETE FROM physical_products;
DELETE FROM service_items;
DELETE FROM technicians;
DELETE FROM members WHERE id > 1; -- 保留 ID=1 的管理员关联会员（如果存在）

-- 重置自增 ID
DELETE FROM sqlite_sequence WHERE name IN ('members', 'technicians', 'service_items', 'physical_products', 'appointments', 'orders', 'fission_logs', 'inventory_logs');

-- 重新启用外键约束
PRAGMA foreign_keys = ON;

-- ============================================================
-- 1. 插入会员数据 (30个会员)
-- ============================================================
INSERT INTO members (created_at, updated_at, name, phone, level, yearly_total_consumption, balance, invitation_code, referrer_id) VALUES
(datetime('now', '-90 days'), datetime('now', '-90 days'), '张三', '13800001001', 'vip', 15800.00, 2000.00, 'INV001', NULL),
(datetime('now', '-85 days'), datetime('now', '-85 days'), '李四', '13800001002', 'vip', 12500.00, 1500.00, 'INV002', NULL),
(datetime('now', '-80 days'), datetime('now', '-80 days'), '王五', '13800001003', 'gold', 8600.00, 800.00, 'INV003', 1),
(datetime('now', '-75 days'), datetime('now', '-75 days'), '赵六', '13800001004', 'silver', 6200.00, 600.00, 'INV004', 1),
(datetime('now', '-70 days'), datetime('now', '-70 days'), '钱七', '13800001005', 'basic', 3200.00, 300.00, 'INV005', 2),
(datetime('now', '-65 days'), datetime('now', '-65 days'), '孙八', '13800001006', 'gold', 9800.00, 1000.00, 'INV006', 2),
(datetime('now', '-60 days'), datetime('now', '-60 days'), '周九', '13800001007', 'silver', 5500.00, 500.00, 'INV007', 3),
(datetime('now', '-55 days'), datetime('now', '-55 days'), '吴十', '13800001008', 'basic', 2800.00, 200.00, 'INV008', 3),
(datetime('now', '-50 days'), datetime('now', '-50 days'), '郑一', '13800001009', 'vip', 18000.00, 2500.00, 'INV009', NULL),
(datetime('now', '-45 days'), datetime('now', '-45 days'), '王二', '13800001010', 'gold', 7200.00, 700.00, 'INV010', 9),
(datetime('now', '-40 days'), datetime('now', '-40 days'), '陈三', '13800001011', 'silver', 4800.00, 400.00, 'INV011', 9),
(datetime('now', '-35 days'), datetime('now', '-35 days'), '林四', '13800001012', 'basic', 2200.00, 150.00, 'INV012', 10),
(datetime('now', '-30 days'), datetime('now', '-30 days'), '黄五', '13800001013', 'gold', 8800.00, 900.00, 'INV013', 1),
(datetime('now', '-28 days'), datetime('now', '-28 days'), '刘六', '13800001014', 'silver', 5200.00, 450.00, 'INV014', 2),
(datetime('now', '-26 days'), datetime('now', '-26 days'), '杨七', '13800001015', 'basic', 2600.00, 250.00, 'INV015', 3),
(datetime('now', '-24 days'), datetime('now', '-24 days'), '何八', '13800001016', 'gold', 7800.00, 750.00, 'INV016', 9),
(datetime('now', '-22 days'), datetime('now', '-22 days'), '罗九', '13800001017', 'silver', 4500.00, 350.00, 'INV017', 1),
(datetime('now', '-20 days'), datetime('now', '-20 days'), '梁十', '13800001018', 'basic', 2100.00, 180.00, 'INV018', 2),
(datetime('now', '-18 days'), datetime('now', '-18 days'), '宋一', '13800001019', 'vip', 16500.00, 2200.00, 'INV019', NULL),
(datetime('now', '-16 days'), datetime('now', '-16 days'), '唐二', '13800001020', 'gold', 8200.00, 820.00, 'INV020', 19),
(datetime('now', '-14 days'), datetime('now', '-14 days'), '许三', '13800001021', 'silver', 5800.00, 550.00, 'INV021', 19),
(datetime('now', '-12 days'), datetime('now', '-12 days'), '韩四', '13800001022', 'basic', 2900.00, 280.00, 'INV022', 20),
(datetime('now', '-10 days'), datetime('now', '-10 days'), '邓五', '13800001023', 'gold', 9200.00, 950.00, 'INV023', 1),
(datetime('now', '-8 days'), datetime('now', '-8 days'), '冯六', '13800001024', 'silver', 6100.00, 600.00, 'INV024', 2),
(datetime('now', '-6 days'), datetime('now', '-6 days'), '曹七', '13800001025', 'basic', 3100.00, 300.00, 'INV025', 9),
(datetime('now', '-5 days'), datetime('now', '-5 days'), '彭八', '13800001026', 'gold', 7500.00, 700.00, 'INV026', 19),
(datetime('now', '-4 days'), datetime('now', '-4 days'), '曾九', '13800001027', 'silver', 4200.00, 380.00, 'INV027', 1),
(datetime('now', '-3 days'), datetime('now', '-3 days'), '肖十', '13800001028', 'basic', 1800.00, 150.00, 'INV028', 9),
(datetime('now', '-2 days'), datetime('now', '-2 days'), '田一', '13800001029', 'gold', 8500.00, 850.00, 'INV029', 19),
(datetime('now', '-1 days'), datetime('now', '-1 days'), '董二', '13800001030', 'basic', 2400.00, 220.00, 'INV030', 1);

-- ============================================================
-- 2. 插入技师数据 (8个技师)
-- ============================================================
INSERT INTO technicians (created_at, updated_at, name, skills, status, average_rating) VALUES
(datetime('now', '-120 days'), datetime('now'), '陈美丽', '["精油SPA", "全身按摩", "足疗"]', 0, 4.8),
(datetime('now', '-110 days'), datetime('now'), '王芳', '["推拿", "拔罐", "刮痧"]', 0, 4.7),
(datetime('now', '-100 days'), datetime('now'), '李娜', '["足底按摩", "艾灸", "精油护理"]', 0, 4.9),
(datetime('now', '-90 days'), datetime('now'), '刘静', '["SPA", "身体护理", "面部护理"]', 1, 4.6),
(datetime('now', '-80 days'), datetime('now'), '张丽', '["推拿", "按摩", "理疗"]', 0, 4.8),
(datetime('now', '-70 days'), datetime('now'), '赵敏', '["足疗", "精油", "艾灸"]', 0, 4.7),
(datetime('now', '-60 days'), datetime('now'), '孙婷', '["SPA", "推拿", "拔罐"]', 2, 4.5),
(datetime('now', '-50 days'), datetime('now'), '周洁', '["按摩", "足疗", "身体护理"]', 0, 4.9);

-- ============================================================
-- 3. 插入服务项目数据 (10个服务项目)
-- ============================================================
INSERT INTO service_items (created_at, updated_at, name, duration, price, is_active) VALUES
(datetime('now', '-150 days'), datetime('now'), '全身精油SPA', 90, 388.00, 1),
(datetime('now', '-150 days'), datetime('now'), '中式推拿', 60, 268.00, 1),
(datetime('now', '-150 days'), datetime('now'), '足底按摩', 45, 158.00, 1),
(datetime('now', '-150 days'), datetime('now'), '艾灸护理', 60, 218.00, 1),
(datetime('now', '-150 days'), datetime('now'), '拔罐刮痧', 45, 188.00, 1),
(datetime('now', '-150 days'), datetime('now'), '背部舒缓按摩', 60, 258.00, 1),
(datetime('now', '-150 days'), datetime('now'), '头部肩颈按摩', 45, 198.00, 1),
(datetime('now', '-150 days'), datetime('now'), '淋巴排毒', 75, 328.00, 1),
(datetime('now', '-150 days'), datetime('now'), '热石疗法', 90, 458.00, 1),
(datetime('now', '-150 days'), datetime('now'), '泰式按摩', 90, 398.00, 1);

-- ============================================================
-- 4. 插入实体商品数据 (12个商品)
-- ============================================================
INSERT INTO physical_products (created_at, updated_at, name, stock, retail_price, cost_price, description, is_active, image_url) VALUES
(datetime('now', '-100 days'), datetime('now'), '薰衣草精油 50ml', 85, 158.00, 80.00, '100%纯天然薰衣草精油，助眠放松', 1, 'https://picsum.photos/seed/prod1/300/300'),
(datetime('now', '-100 days'), datetime('now'), '玫瑰精油 30ml', 62, 288.00, 150.00, '保加利亚玫瑰精油，美容养颜', 1, 'https://picsum.photos/seed/prod2/300/300'),
(datetime('now', '-100 days'), datetime('now'), '茶树精油 50ml', 48, 128.00, 65.00, '澳洲茶树精油，抗菌消炎', 1, 'https://picsum.photos/seed/prod3/300/300'),
(datetime('now', '-100 days'), datetime('now'), '按摩精油套装', 35, 398.00, 200.00, '6种精油组合装，适合全身按摩', 1, 'https://picsum.photos/seed/prod4/300/300'),
(datetime('now', '-100 days'), datetime('now'), '足部护理霜 100g', 120, 88.00, 45.00, '深层滋养，去除足部死皮', 1, 'https://picsum.photos/seed/prod5/300/300'),
(datetime('now', '-100 days'), datetime('now'), '身体乳液 200ml', 95, 168.00, 85.00, '保湿滋润，改善肌肤干燥', 1, 'https://picsum.photos/seed/prod6/300/300'),
(datetime('now', '-100 days'), datetime('now'), '艾灸贴 10片装', 150, 58.00, 28.00, '自发热艾灸贴，缓解肌肉疲劳', 1, 'https://picsum.photos/seed/prod7/300/300'),
(datetime('now', '-100 days'), datetime('now'), '按摩油 500ml', 42, 218.00, 110.00, '专业按摩用油，润滑不油腻', 1, 'https://picsum.photos/seed/prod8/300/300'),
(datetime('now', '-100 days'), datetime('now'), '香薰蜡烛套装', 68, 128.00, 65.00, '3支装，多种香味可选', 1, 'https://picsum.photos/seed/prod9/300/300'),
(datetime('now', '-100 days'), datetime('now'), '面部精华液 30ml', 55, 328.00, 165.00, '抗衰老精华，提拉紧致', 1, 'https://picsum.photos/seed/prod10/300/300'),
(datetime('now', '-100 days'), datetime('now'), '竹炭足贴 20片', 180, 68.00, 35.00, '吸附毒素，改善睡眠', 1, 'https://picsum.photos/seed/prod11/300/300'),
(datetime('now', '-100 days'), datetime('now'), '按摩刮痧板', 75, 98.00, 50.00, '天然玉石，促进血液循环', 1, 'https://picsum.photos/seed/prod12/300/300');

-- ============================================================
-- 5. 插入预约记录（手动插入确定的数据，避免随机函数问题）
-- ============================================================

-- 近30天的预约，每天3-6条
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at) VALUES
-- 30天前
(1, 1, 1, datetime('now', '-30 days', '+10 hours'), datetime('now', '-30 days', '+11 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-30 days'), datetime('now', '-30 days')),
(3, 2, 2, datetime('now', '-30 days', '+14 hours'), datetime('now', '-30 days', '+15 hours'), 'completed', 268.00, 268.00, datetime('now', '-30 days'), datetime('now', '-30 days')),
(5, 3, 3, datetime('now', '-30 days', '+16 hours'), datetime('now', '-30 days', '+16 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-30 days'), datetime('now', '-30 days')),
(7, 5, 4, datetime('now', '-30 days', '+18 hours'), datetime('now', '-30 days', '+19 hours'), 'completed', 218.00, 198.00, datetime('now', '-30 days'), datetime('now', '-30 days')),

-- 29天前
(2, 6, 1, datetime('now', '-29 days', '+10 hours'), datetime('now', '-29 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-29 days'), datetime('now', '-29 days')),
(4, 8, 3, datetime('now', '-29 days', '+14 hours'), datetime('now', '-29 days', '+14 hours 45 minutes'), 'completed', 158.00, 138.00, datetime('now', '-29 days'), datetime('now', '-29 days')),
(6, 1, 2, datetime('now', '-29 days', '+16 hours'), datetime('now', '-29 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-29 days'), datetime('now', '-29 days')),
(8, 3, 5, datetime('now', '-29 days', '+11 hours'), datetime('now', '-29 days', '+11 hours 45 minutes'), 'completed', 188.00, 188.00, datetime('now', '-29 days'), datetime('now', '-29 days')),
(10, 5, 1, datetime('now', '-29 days', '+15 hours'), datetime('now', '-29 days', '+16 hours 30 minutes'), 'completed', 388.00, 348.00, datetime('now', '-29 days'), datetime('now', '-29 days')),

-- 28天前
(1, 2, 6, datetime('now', '-28 days', '+10 hours'), datetime('now', '-28 days', '+11 hours'), 'completed', 258.00, 258.00, datetime('now', '-28 days'), datetime('now', '-28 days')),
(9, 6, 1, datetime('now', '-28 days', '+14 hours'), datetime('now', '-28 days', '+15 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-28 days'), datetime('now', '-28 days')),
(11, 8, 3, datetime('now', '-28 days', '+16 hours'), datetime('now', '-28 days', '+16 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-28 days'), datetime('now', '-28 days')),
(13, 1, 2, datetime('now', '-28 days', '+11 hours'), datetime('now', '-28 days', '+12 hours'), 'completed', 268.00, 238.00, datetime('now', '-28 days'), datetime('now', '-28 days')),

-- 27天前
(2, 5, 1, datetime('now', '-27 days', '+10 hours'), datetime('now', '-27 days', '+11 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-27 days'), datetime('now', '-27 days')),
(4, 6, 3, datetime('now', '-27 days', '+14 hours'), datetime('now', '-27 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-27 days'), datetime('now', '-27 days')),
(6, 8, 2, datetime('now', '-27 days', '+16 hours'), datetime('now', '-27 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-27 days'), datetime('now', '-27 days')),
(12, 1, 4, datetime('now', '-27 days', '+11 hours'), datetime('now', '-27 days', '+12 hours'), 'completed', 218.00, 218.00, datetime('now', '-27 days'), datetime('now', '-27 days')),

-- 26天前
(3, 5, 1, datetime('now', '-26 days', '+10 hours'), datetime('now', '-26 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-26 days'), datetime('now', '-26 days')),
(5, 6, 3, datetime('now', '-26 days', '+14 hours'), datetime('now', '-26 days', '+14 hours 45 minutes'), 'completed', 158.00, 138.00, datetime('now', '-26 days'), datetime('now', '-26 days')),
(7, 8, 2, datetime('now', '-26 days', '+16 hours'), datetime('now', '-26 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-26 days'), datetime('now', '-26 days')),
(9, 1, 8, datetime('now', '-26 days', '+11 hours'), datetime('now', '-26 days', '+12 hours 15 minutes'), 'completed', 328.00, 328.00, datetime('now', '-26 days'), datetime('now', '-26 days')),

-- 25天前
(1, 3, 1, datetime('now', '-25 days', '+10 hours'), datetime('now', '-25 days', '+11 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-25 days'), datetime('now', '-25 days')),
(11, 5, 3, datetime('now', '-25 days', '+14 hours'), datetime('now', '-25 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-25 days'), datetime('now', '-25 days')),
(13, 6, 2, datetime('now', '-25 days', '+16 hours'), datetime('now', '-25 days', '+17 hours'), 'completed', 268.00, 248.00, datetime('now', '-25 days'), datetime('now', '-25 days')),
(15, 8, 4, datetime('now', '-25 days', '+18 hours'), datetime('now', '-25 days', '+19 hours'), 'completed', 218.00, 218.00, datetime('now', '-25 days'), datetime('now', '-25 days')),
(17, 1, 7, datetime('now', '-25 days', '+13 hours'), datetime('now', '-25 days', '+13 hours 45 minutes'), 'completed', 198.00, 198.00, datetime('now', '-25 days'), datetime('now', '-25 days')),

-- 24天前
(2, 2, 1, datetime('now', '-24 days', '+10 hours'), datetime('now', '-24 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-24 days'), datetime('now', '-24 days')),
(4, 3, 3, datetime('now', '-24 days', '+14 hours'), datetime('now', '-24 days', '+14 hours 45 minutes'), 'completed', 158.00, 138.00, datetime('now', '-24 days'), datetime('now', '-24 days')),
(6, 5, 2, datetime('now', '-24 days', '+16 hours'), datetime('now', '-24 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-24 days'), datetime('now', '-24 days')),
(8, 6, 5, datetime('now', '-24 days', '+11 hours'), datetime('now', '-24 days', '+11 hours 45 minutes'), 'completed', 188.00, 188.00, datetime('now', '-24 days'), datetime('now', '-24 days')),

-- 23天前
(10, 1, 1, datetime('now', '-23 days', '+10 hours'), datetime('now', '-23 days', '+11 hours 30 minutes'), 'completed', 388.00, 358.00, datetime('now', '-23 days'), datetime('now', '-23 days')),
(12, 8, 3, datetime('now', '-23 days', '+14 hours'), datetime('now', '-23 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-23 days'), datetime('now', '-23 days')),
(14, 2, 2, datetime('now', '-23 days', '+16 hours'), datetime('now', '-23 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-23 days'), datetime('now', '-23 days')),
(16, 3, 6, datetime('now', '-23 days', '+11 hours'), datetime('now', '-23 days', '+12 hours'), 'completed', 258.00, 258.00, datetime('now', '-23 days'), datetime('now', '-23 days')),
(18, 5, 1, datetime('now', '-23 days', '+15 hours'), datetime('now', '-23 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-23 days'), datetime('now', '-23 days')),

-- 22天前
(1, 6, 1, datetime('now', '-22 days', '+10 hours'), datetime('now', '-22 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-22 days'), datetime('now', '-22 days')),
(3, 8, 3, datetime('now', '-22 days', '+14 hours'), datetime('now', '-22 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-22 days'), datetime('now', '-22 days')),
(5, 1, 2, datetime('now', '-22 days', '+16 hours'), datetime('now', '-22 days', '+17 hours'), 'completed', 268.00, 248.00, datetime('now', '-22 days'), datetime('now', '-22 days')),
(19, 2, 9, datetime('now', '-22 days', '+11 hours'), datetime('now', '-22 days', '+12 hours 30 minutes'), 'completed', 458.00, 428.00, datetime('now', '-22 days'), datetime('now', '-22 days')),

-- 21天前
(7, 3, 1, datetime('now', '-21 days', '+10 hours'), datetime('now', '-21 days', '+11 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-21 days'), datetime('now', '-21 days')),
(9, 5, 3, datetime('now', '-21 days', '+14 hours'), datetime('now', '-21 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-21 days'), datetime('now', '-21 days')),
(11, 6, 2, datetime('now', '-21 days', '+16 hours'), datetime('now', '-21 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-21 days'), datetime('now', '-21 days')),
(13, 8, 10, datetime('now', '-21 days', '+11 hours'), datetime('now', '-21 days', '+12 hours 30 minutes'), 'completed', 398.00, 398.00, datetime('now', '-21 days'), datetime('now', '-21 days')),
(20, 1, 1, datetime('now', '-21 days', '+15 hours'), datetime('now', '-21 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-21 days'), datetime('now', '-21 days')),

-- 20天前
(2, 2, 1, datetime('now', '-20 days', '+10 hours'), datetime('now', '-20 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-20 days'), datetime('now', '-20 days')),
(4, 3, 3, datetime('now', '-20 days', '+14 hours'), datetime('now', '-20 days', '+14 hours 45 minutes'), 'completed', 158.00, 138.00, datetime('now', '-20 days'), datetime('now', '-20 days')),
(6, 5, 2, datetime('now', '-20 days', '+16 hours'), datetime('now', '-20 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-20 days'), datetime('now', '-20 days')),
(21, 6, 4, datetime('now', '-20 days', '+11 hours'), datetime('now', '-20 days', '+12 hours'), 'completed', 218.00, 218.00, datetime('now', '-20 days'), datetime('now', '-20 days')),

-- 19天前
(8, 1, 1, datetime('now', '-19 days', '+10 hours'), datetime('now', '-19 days', '+11 hours 30 minutes'), 'completed', 388.00, 358.00, datetime('now', '-19 days'), datetime('now', '-19 days')),
(10, 8, 3, datetime('now', '-19 days', '+14 hours'), datetime('now', '-19 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-19 days'), datetime('now', '-19 days')),
(12, 2, 2, datetime('now', '-19 days', '+16 hours'), datetime('now', '-19 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-19 days'), datetime('now', '-19 days')),
(14, 3, 7, datetime('now', '-19 days', '+11 hours'), datetime('now', '-19 days', '+11 hours 45 minutes'), 'completed', 198.00, 198.00, datetime('now', '-19 days'), datetime('now', '-19 days')),
(22, 5, 1, datetime('now', '-19 days', '+15 hours'), datetime('now', '-19 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-19 days'), datetime('now', '-19 days')),

-- 18天前
(1, 6, 1, datetime('now', '-18 days', '+10 hours'), datetime('now', '-18 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-18 days'), datetime('now', '-18 days')),
(3, 8, 3, datetime('now', '-18 days', '+14 hours'), datetime('now', '-18 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-18 days'), datetime('now', '-18 days')),
(5, 1, 2, datetime('now', '-18 days', '+16 hours'), datetime('now', '-18 days', '+17 hours'), 'completed', 268.00, 248.00, datetime('now', '-18 days'), datetime('now', '-18 days')),
(23, 2, 5, datetime('now', '-18 days', '+11 hours'), datetime('now', '-18 days', '+11 hours 45 minutes'), 'completed', 188.00, 188.00, datetime('now', '-18 days'), datetime('now', '-18 days')),

-- 17天前
(7, 3, 1, datetime('now', '-17 days', '+10 hours'), datetime('now', '-17 days', '+11 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-17 days'), datetime('now', '-17 days')),
(9, 5, 3, datetime('now', '-17 days', '+14 hours'), datetime('now', '-17 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-17 days'), datetime('now', '-17 days')),
(11, 6, 2, datetime('now', '-17 days', '+16 hours'), datetime('now', '-17 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-17 days'), datetime('now', '-17 days')),
(13, 8, 8, datetime('now', '-17 days', '+11 hours'), datetime('now', '-17 days', '+12 hours 15 minutes'), 'completed', 328.00, 328.00, datetime('now', '-17 days'), datetime('now', '-17 days')),
(24, 1, 1, datetime('now', '-17 days', '+15 hours'), datetime('now', '-17 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-17 days'), datetime('now', '-17 days')),

-- 16天前
(15, 2, 1, datetime('now', '-16 days', '+10 hours'), datetime('now', '-16 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-16 days'), datetime('now', '-16 days')),
(17, 3, 3, datetime('now', '-16 days', '+14 hours'), datetime('now', '-16 days', '+14 hours 45 minutes'), 'completed', 158.00, 138.00, datetime('now', '-16 days'), datetime('now', '-16 days')),
(19, 5, 2, datetime('now', '-16 days', '+16 hours'), datetime('now', '-16 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-16 days'), datetime('now', '-16 days')),
(25, 6, 4, datetime('now', '-16 days', '+11 hours'), datetime('now', '-16 days', '+12 hours'), 'completed', 218.00, 218.00, datetime('now', '-16 days'), datetime('now', '-16 days')),

-- 15天前
(2, 1, 1, datetime('now', '-15 days', '+10 hours'), datetime('now', '-15 days', '+11 hours 30 minutes'), 'completed', 388.00, 358.00, datetime('now', '-15 days'), datetime('now', '-15 days')),
(4, 8, 3, datetime('now', '-15 days', '+14 hours'), datetime('now', '-15 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-15 days'), datetime('now', '-15 days')),
(6, 2, 2, datetime('now', '-15 days', '+16 hours'), datetime('now', '-15 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-15 days'), datetime('now', '-15 days')),
(8, 3, 6, datetime('now', '-15 days', '+11 hours'), datetime('now', '-15 days', '+12 hours'), 'completed', 258.00, 258.00, datetime('now', '-15 days'), datetime('now', '-15 days')),
(26, 5, 1, datetime('now', '-15 days', '+15 hours'), datetime('now', '-15 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-15 days'), datetime('now', '-15 days')),

-- 14天前
(10, 6, 1, datetime('now', '-14 days', '+10 hours'), datetime('now', '-14 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-14 days'), datetime('now', '-14 days')),
(12, 8, 3, datetime('now', '-14 days', '+14 hours'), datetime('now', '-14 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-14 days'), datetime('now', '-14 days')),
(14, 1, 2, datetime('now', '-14 days', '+16 hours'), datetime('now', '-14 days', '+17 hours'), 'completed', 268.00, 248.00, datetime('now', '-14 days'), datetime('now', '-14 days')),
(27, 2, 7, datetime('now', '-14 days', '+11 hours'), datetime('now', '-14 days', '+11 hours 45 minutes'), 'completed', 198.00, 198.00, datetime('now', '-14 days'), datetime('now', '-14 days')),

-- 13天前
(16, 3, 1, datetime('now', '-13 days', '+10 hours'), datetime('now', '-13 days', '+11 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-13 days'), datetime('now', '-13 days')),
(18, 5, 3, datetime('now', '-13 days', '+14 hours'), datetime('now', '-13 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-13 days'), datetime('now', '-13 days')),
(20, 6, 2, datetime('now', '-13 days', '+16 hours'), datetime('now', '-13 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-13 days'), datetime('now', '-13 days')),
(22, 8, 5, datetime('now', '-13 days', '+11 hours'), datetime('now', '-13 days', '+11 hours 45 minutes'), 'completed', 188.00, 188.00, datetime('now', '-13 days'), datetime('now', '-13 days')),
(28, 1, 1, datetime('now', '-13 days', '+15 hours'), datetime('now', '-13 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-13 days'), datetime('now', '-13 days')),

-- 12天前
(1, 2, 1, datetime('now', '-12 days', '+10 hours'), datetime('now', '-12 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-12 days'), datetime('now', '-12 days')),
(3, 3, 3, datetime('now', '-12 days', '+14 hours'), datetime('now', '-12 days', '+14 hours 45 minutes'), 'completed', 158.00, 138.00, datetime('now', '-12 days'), datetime('now', '-12 days')),
(5, 5, 2, datetime('now', '-12 days', '+16 hours'), datetime('now', '-12 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-12 days'), datetime('now', '-12 days')),
(29, 6, 4, datetime('now', '-12 days', '+11 hours'), datetime('now', '-12 days', '+12 hours'), 'completed', 218.00, 218.00, datetime('now', '-12 days'), datetime('now', '-12 days')),

-- 11天前
(7, 1, 1, datetime('now', '-11 days', '+10 hours'), datetime('now', '-11 days', '+11 hours 30 minutes'), 'completed', 388.00, 358.00, datetime('now', '-11 days'), datetime('now', '-11 days')),
(9, 8, 3, datetime('now', '-11 days', '+14 hours'), datetime('now', '-11 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-11 days'), datetime('now', '-11 days')),
(11, 2, 2, datetime('now', '-11 days', '+16 hours'), datetime('now', '-11 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-11 days'), datetime('now', '-11 days')),
(13, 3, 9, datetime('now', '-11 days', '+11 hours'), datetime('now', '-11 days', '+12 hours 30 minutes'), 'completed', 458.00, 428.00, datetime('now', '-11 days'), datetime('now', '-11 days')),
(30, 5, 1, datetime('now', '-11 days', '+15 hours'), datetime('now', '-11 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-11 days'), datetime('now', '-11 days')),

-- 10天前
(15, 6, 1, datetime('now', '-10 days', '+10 hours'), datetime('now', '-10 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-10 days'), datetime('now', '-10 days')),
(17, 8, 3, datetime('now', '-10 days', '+14 hours'), datetime('now', '-10 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-10 days'), datetime('now', '-10 days')),
(19, 1, 2, datetime('now', '-10 days', '+16 hours'), datetime('now', '-10 days', '+17 hours'), 'completed', 268.00, 248.00, datetime('now', '-10 days'), datetime('now', '-10 days')),
(21, 2, 10, datetime('now', '-10 days', '+11 hours'), datetime('now', '-10 days', '+12 hours 30 minutes'), 'completed', 398.00, 398.00, datetime('now', '-10 days'), datetime('now', '-10 days')),

-- 9天前
(2, 3, 1, datetime('now', '-9 days', '+10 hours'), datetime('now', '-9 days', '+11 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-9 days'), datetime('now', '-9 days')),
(4, 5, 3, datetime('now', '-9 days', '+14 hours'), datetime('now', '-9 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-9 days'), datetime('now', '-9 days')),
(6, 6, 2, datetime('now', '-9 days', '+16 hours'), datetime('now', '-9 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-9 days'), datetime('now', '-9 days')),
(23, 8, 8, datetime('now', '-9 days', '+11 hours'), datetime('now', '-9 days', '+12 hours 15 minutes'), 'completed', 328.00, 328.00, datetime('now', '-9 days'), datetime('now', '-9 days')),
(25, 1, 1, datetime('now', '-9 days', '+15 hours'), datetime('now', '-9 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-9 days'), datetime('now', '-9 days')),

-- 8天前
(8, 2, 1, datetime('now', '-8 days', '+10 hours'), datetime('now', '-8 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-8 days'), datetime('now', '-8 days')),
(10, 3, 3, datetime('now', '-8 days', '+14 hours'), datetime('now', '-8 days', '+14 hours 45 minutes'), 'completed', 158.00, 138.00, datetime('now', '-8 days'), datetime('now', '-8 days')),
(12, 5, 2, datetime('now', '-8 days', '+16 hours'), datetime('now', '-8 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-8 days'), datetime('now', '-8 days')),
(24, 6, 4, datetime('now', '-8 days', '+11 hours'), datetime('now', '-8 days', '+12 hours'), 'completed', 218.00, 218.00, datetime('now', '-8 days'), datetime('now', '-8 days')),

-- 7天前
(14, 1, 1, datetime('now', '-7 days', '+10 hours'), datetime('now', '-7 days', '+11 hours 30 minutes'), 'completed', 388.00, 358.00, datetime('now', '-7 days'), datetime('now', '-7 days')),
(16, 8, 3, datetime('now', '-7 days', '+14 hours'), datetime('now', '-7 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-7 days'), datetime('now', '-7 days')),
(18, 2, 2, datetime('now', '-7 days', '+16 hours'), datetime('now', '-7 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-7 days'), datetime('now', '-7 days')),
(20, 3, 6, datetime('now', '-7 days', '+11 hours'), datetime('now', '-7 days', '+12 hours'), 'completed', 258.00, 258.00, datetime('now', '-7 days'), datetime('now', '-7 days')),
(26, 5, 1, datetime('now', '-7 days', '+15 hours'), datetime('now', '-7 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-7 days'), datetime('now', '-7 days')),

-- 6天前
(22, 6, 1, datetime('now', '-6 days', '+10 hours'), datetime('now', '-6 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-6 days'), datetime('now', '-6 days')),
(1, 8, 3, datetime('now', '-6 days', '+14 hours'), datetime('now', '-6 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-6 days'), datetime('now', '-6 days')),
(3, 1, 2, datetime('now', '-6 days', '+16 hours'), datetime('now', '-6 days', '+17 hours'), 'completed', 268.00, 248.00, datetime('now', '-6 days'), datetime('now', '-6 days')),
(27, 2, 7, datetime('now', '-6 days', '+11 hours'), datetime('now', '-6 days', '+11 hours 45 minutes'), 'completed', 198.00, 198.00, datetime('now', '-6 days'), datetime('now', '-6 days')),

-- 5天前
(5, 3, 1, datetime('now', '-5 days', '+10 hours'), datetime('now', '-5 days', '+11 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-5 days'), datetime('now', '-5 days')),
(7, 5, 3, datetime('now', '-5 days', '+14 hours'), datetime('now', '-5 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-5 days'), datetime('now', '-5 days')),
(9, 6, 2, datetime('now', '-5 days', '+16 hours'), datetime('now', '-5 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-5 days'), datetime('now', '-5 days')),
(11, 8, 5, datetime('now', '-5 days', '+11 hours'), datetime('now', '-5 days', '+11 hours 45 minutes'), 'completed', 188.00, 188.00, datetime('now', '-5 days'), datetime('now', '-5 days')),
(28, 1, 1, datetime('now', '-5 days', '+15 hours'), datetime('now', '-5 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-5 days'), datetime('now', '-5 days')),

-- 4天前
(13, 2, 1, datetime('now', '-4 days', '+10 hours'), datetime('now', '-4 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-4 days'), datetime('now', '-4 days')),
(15, 3, 3, datetime('now', '-4 days', '+14 hours'), datetime('now', '-4 days', '+14 hours 45 minutes'), 'completed', 158.00, 138.00, datetime('now', '-4 days'), datetime('now', '-4 days')),
(17, 5, 2, datetime('now', '-4 days', '+16 hours'), datetime('now', '-4 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-4 days'), datetime('now', '-4 days')),
(29, 6, 4, datetime('now', '-4 days', '+11 hours'), datetime('now', '-4 days', '+12 hours'), 'completed', 218.00, 218.00, datetime('now', '-4 days'), datetime('now', '-4 days')),

-- 3天前
(19, 1, 1, datetime('now', '-3 days', '+10 hours'), datetime('now', '-3 days', '+11 hours 30 minutes'), 'completed', 388.00, 358.00, datetime('now', '-3 days'), datetime('now', '-3 days')),
(21, 8, 3, datetime('now', '-3 days', '+14 hours'), datetime('now', '-3 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-3 days'), datetime('now', '-3 days')),
(23, 2, 2, datetime('now', '-3 days', '+16 hours'), datetime('now', '-3 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-3 days'), datetime('now', '-3 days')),
(25, 3, 9, datetime('now', '-3 days', '+11 hours'), datetime('now', '-3 days', '+12 hours 30 minutes'), 'completed', 458.00, 428.00, datetime('now', '-3 days'), datetime('now', '-3 days')),
(30, 5, 1, datetime('now', '-3 days', '+15 hours'), datetime('now', '-3 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-3 days'), datetime('now', '-3 days')),

-- 2天前
(2, 6, 1, datetime('now', '-2 days', '+10 hours'), datetime('now', '-2 days', '+11 hours 30 minutes'), 'completed', 388.00, 388.00, datetime('now', '-2 days'), datetime('now', '-2 days')),
(4, 8, 3, datetime('now', '-2 days', '+14 hours'), datetime('now', '-2 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-2 days'), datetime('now', '-2 days')),
(6, 1, 2, datetime('now', '-2 days', '+16 hours'), datetime('now', '-2 days', '+17 hours'), 'completed', 268.00, 248.00, datetime('now', '-2 days'), datetime('now', '-2 days')),
(26, 2, 10, datetime('now', '-2 days', '+11 hours'), datetime('now', '-2 days', '+12 hours 30 minutes'), 'completed', 398.00, 398.00, datetime('now', '-2 days'), datetime('now', '-2 days')),

-- 1天前
(8, 3, 1, datetime('now', '-1 days', '+10 hours'), datetime('now', '-1 days', '+11 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-1 days'), datetime('now', '-1 days')),
(10, 5, 3, datetime('now', '-1 days', '+14 hours'), datetime('now', '-1 days', '+14 hours 45 minutes'), 'completed', 158.00, 158.00, datetime('now', '-1 days'), datetime('now', '-1 days')),
(12, 6, 2, datetime('now', '-1 days', '+16 hours'), datetime('now', '-1 days', '+17 hours'), 'completed', 268.00, 268.00, datetime('now', '-1 days'), datetime('now', '-1 days')),
(14, 8, 8, datetime('now', '-1 days', '+11 hours'), datetime('now', '-1 days', '+12 hours 15 minutes'), 'completed', 328.00, 328.00, datetime('now', '-1 days'), datetime('now', '-1 days')),
(27, 1, 1, datetime('now', '-1 days', '+15 hours'), datetime('now', '-1 days', '+16 hours 30 minutes'), 'completed', 388.00, 368.00, datetime('now', '-1 days'), datetime('now', '-1 days'));

-- ============================================================
-- 6. 插入订单记录 (服务订单)
-- ============================================================

-- 为每个已完成的预约创建服务订单
INSERT INTO orders (member_id, appointment_id, product_id, type, quantity, total_amount, actual_paid, status, payment_method, created_at, updated_at)
SELECT
    member_id,
    id as appointment_id,
    NULL as product_id,
    'service' as type,
    1 as quantity,
    actual_price as total_amount,
    actual_price as actual_paid,
    'completed' as status,
    CASE (id % 3)
        WHEN 0 THEN 'cash'
        WHEN 1 THEN 'wechat'
        ELSE 'alipay'
    END as payment_method,
    created_at,
    updated_at
FROM appointments
WHERE status = 'completed';

-- 插入商品订单（近30天，手动插入）
INSERT INTO orders (member_id, appointment_id, product_id, type, quantity, total_amount, actual_paid, status, payment_method, created_at, updated_at) VALUES
(1, NULL, 1, 'product', 2, 316.00, 316.00, 'completed', 'wechat', datetime('now', '-28 days'), datetime('now', '-28 days')),
(2, NULL, 2, 'product', 1, 288.00, 288.00, 'completed', 'alipay', datetime('now', '-27 days'), datetime('now', '-27 days')),
(3, NULL, 5, 'product', 3, 264.00, 264.00, 'completed', 'cash', datetime('now', '-26 days'), datetime('now', '-26 days')),
(4, NULL, 7, 'product', 5, 290.00, 290.00, 'completed', 'wechat', datetime('now', '-25 days'), datetime('now', '-25 days')),
(5, NULL, 11, 'product', 4, 272.00, 272.00, 'completed', 'alipay', datetime('now', '-24 days'), datetime('now', '-24 days')),
(6, NULL, 3, 'product', 2, 256.00, 256.00, 'completed', 'cash', datetime('now', '-23 days'), datetime('now', '-23 days')),
(7, NULL, 6, 'product', 1, 168.00, 168.00, 'completed', 'wechat', datetime('now', '-22 days'), datetime('now', '-22 days')),
(8, NULL, 9, 'product', 2, 256.00, 256.00, 'completed', 'alipay', datetime('now', '-21 days'), datetime('now', '-21 days')),
(9, NULL, 1, 'product', 1, 158.00, 158.00, 'completed', 'cash', datetime('now', '-20 days'), datetime('now', '-20 days')),
(10, NULL, 4, 'product', 1, 398.00, 398.00, 'completed', 'wechat', datetime('now', '-19 days'), datetime('now', '-19 days')),
(11, NULL, 8, 'product', 1, 218.00, 218.00, 'completed', 'alipay', datetime('now', '-18 days'), datetime('now', '-18 days')),
(12, NULL, 12, 'product', 2, 196.00, 196.00, 'completed', 'cash', datetime('now', '-17 days'), datetime('now', '-17 days')),
(13, NULL, 2, 'product', 1, 288.00, 288.00, 'completed', 'wechat', datetime('now', '-16 days'), datetime('now', '-16 days')),
(14, NULL, 5, 'product', 2, 176.00, 176.00, 'completed', 'alipay', datetime('now', '-15 days'), datetime('now', '-15 days')),
(15, NULL, 10, 'product', 1, 328.00, 328.00, 'completed', 'cash', datetime('now', '-14 days'), datetime('now', '-14 days')),
(16, NULL, 1, 'product', 2, 316.00, 316.00, 'completed', 'wechat', datetime('now', '-13 days'), datetime('now', '-13 days')),
(17, NULL, 7, 'product', 3, 174.00, 174.00, 'completed', 'alipay', datetime('now', '-12 days'), datetime('now', '-12 days')),
(18, NULL, 11, 'product', 6, 408.00, 408.00, 'completed', 'cash', datetime('now', '-11 days'), datetime('now', '-11 days')),
(19, NULL, 3, 'product', 1, 128.00, 128.00, 'completed', 'wechat', datetime('now', '-10 days'), datetime('now', '-10 days')),
(20, NULL, 6, 'product', 2, 336.00, 336.00, 'completed', 'alipay', datetime('now', '-9 days'), datetime('now', '-9 days')),
(21, NULL, 9, 'product', 1, 128.00, 128.00, 'completed', 'cash', datetime('now', '-8 days'), datetime('now', '-8 days')),
(22, NULL, 4, 'product', 1, 398.00, 398.00, 'completed', 'wechat', datetime('now', '-7 days'), datetime('now', '-7 days')),
(23, NULL, 8, 'product', 2, 436.00, 436.00, 'completed', 'alipay', datetime('now', '-6 days'), datetime('now', '-6 days')),
(24, NULL, 12, 'product', 3, 294.00, 294.00, 'completed', 'cash', datetime('now', '-5 days'), datetime('now', '-5 days')),
(25, NULL, 1, 'product', 1, 158.00, 158.00, 'completed', 'wechat', datetime('now', '-4 days'), datetime('now', '-4 days')),
(26, NULL, 2, 'product', 1, 288.00, 288.00, 'completed', 'alipay', datetime('now', '-3 days'), datetime('now', '-3 days')),
(27, NULL, 5, 'product', 2, 176.00, 176.00, 'completed', 'cash', datetime('now', '-2 days'), datetime('now', '-2 days')),
(28, NULL, 10, 'product', 1, 328.00, 328.00, 'completed', 'wechat', datetime('now', '-1 days'), datetime('now', '-1 days'));

-- ============================================================
-- 7. 插入裂变记录 (佣金记录)
-- ============================================================
INSERT INTO fission_logs (inviter_id, invitee_id, commission_amount, order_id, created_at, updated_at)
SELECT
    m.referrer_id as inviter_id,
    m.id as invitee_id,
    o.actual_paid * 0.05 as commission_amount,
    o.id as order_id,
    o.created_at,
    o.updated_at
FROM members m
JOIN orders o ON o.member_id = m.id
WHERE m.referrer_id IS NOT NULL
AND o.status = 'completed';

-- ============================================================
-- 8. 插入库存变动记录（批量入库）
-- ============================================================
INSERT INTO inventory_logs (product_id, operator_id, change
_amount, action_type, before_stock, after_stock, order_id, remark, created_at, updated_at) VALUES
(1, 1, 100, 'restock', 0, 100, NULL, '供应商补货', datetime('now', '-20 days'), datetime('now', '-20 days')),
(2, 1, 80, 'restock', 0, 80, NULL, '供应商补货', datetime('now', '-20 days'), datetime('now', '-20 days')),
(3, 1, 60, 'restock', 0, 60, NULL, '供应商补货', datetime('now', '-20 days'), datetime('now', '-20 days')),
(4, 1, 50, 'restock', 0, 50, NULL, '供应商补货', datetime('now', '-18 days'), datetime('now', '-18 days')),
(5, 1, 150, 'restock', 0, 150, NULL, '供应商补货', datetime('now', '-18 days'), datetime('now', '-18 days')),
(6, 1, 120, 'restock', 0, 120, NULL, '供应商补货', datetime('now', '-15 days'), datetime('now', '-15 days')),
(7, 1, 200, 'restock', 0, 200, NULL, '供应商补货', datetime('now', '-15 days'), datetime('now', '-15 days')),
(8, 1, 60, 'restock', 0, 60, NULL, '供应商补货', datetime('now', '-12 days'), datetime('now', '-12 days')),
(9, 1, 80, 'restock', 0, 80, NULL, '供应商补货', datetime('now', '-12 days'), datetime('now', '-12 days')),
(10, 1, 70, 'restock', 0, 70, NULL, '供应商补货', datetime('now', '-10 days'), datetime('now', '-10 days')),
(11, 1, 200, 'restock', 0, 200, NULL, '供应商补货', datetime('now', '-8 days'), datetime('now', '-8 days')),
(12, 1, 90, 'restock', 0, 90, NULL, '供应商补货', datetime('now', '-8 days'), datetime('now', '-8 days'));
