-- Dashboard 测试数据 - 简化版
PRAGMA foreign_keys = OFF;

-- 清理数据
DELETE FROM inventory_logs;
DELETE FROM fission_logs;
DELETE FROM orders;
DELETE FROM appointments;
DELETE FROM schedules;
DELETE FROM physical_products;
DELETE FROM service_items;
DELETE FROM technicians;
DELETE FROM members WHERE phone LIKE '138%';
DELETE FROM sqlite_sequence;

PRAGMA foreign_keys = ON;

-- 会员（30个）
INSERT INTO members (created_at, updated_at, name, phone, level, yearly_total_consumption, balance, invitation_code, referrer_id) VALUES
(datetime('now', '-90 days'), datetime('now'), '张三', '13800001001', 'vip', 15800, 2000, 'INV001', NULL),
(datetime('now', '-85 days'), datetime('now'), '李四', '13800001002', 'vip', 12500, 1500, 'INV002', NULL),
(datetime('now', '-80 days'), datetime('now'), '王五', '13800001003', 'gold', 8600, 800, 'INV003', 1),
(datetime('now', '-75 days'), datetime('now'), '赵六', '13800001004', 'silver', 6200, 600, 'INV004', 1),
(datetime('now', '-70 days'), datetime('now'), '钱七', '13800001005', 'basic', 3200, 300, 'INV005', 2),
(datetime('now', '-65 days'), datetime('now'), '孙八', '13800001006', 'gold', 9800, 1000, 'INV006', 2),
(datetime('now', '-60 days'), datetime('now'), '周九', '13800001007', 'silver', 5500, 500, 'INV007', 3),
(datetime('now', '-55 days'), datetime('now'), '吴十', '13800001008', 'basic', 2800, 200, 'INV008', 3),
(datetime('now', '-50 days'), datetime('now'), '郑一', '13800001009', 'vip', 18000, 2500, 'INV009', NULL),
(datetime('now', '-45 days'), datetime('now'), '王二', '13800001010', 'gold', 7200, 700, 'INV010', 9),
(datetime('now', '-40 days'), datetime('now'), '陈三', '13800001011', 'silver', 4800, 400, 'INV011', 9),
(datetime('now', '-35 days'), datetime('now'), '林四', '13800001012', 'basic', 2200, 150, 'INV012', 10),
(datetime('now', '-30 days'), datetime('now'), '黄五', '13800001013', 'gold', 8800, 900, 'INV013', 1),
(datetime('now', '-28 days'), datetime('now'), '刘六', '13800001014', 'silver', 5200, 450, 'INV014', 2),
(datetime('now', '-26 days'), datetime('now'), '杨七', '13800001015', 'basic', 2600, 250, 'INV015', 3),
(datetime('now', '-24 days'), datetime('now'), '何八', '13800001016', 'gold', 7800, 750, 'INV016', 9),
(datetime('now', '-22 days'), datetime('now'), '罗九', '13800001017', 'silver', 4500, 350, 'INV017', 1),
(datetime('now', '-20 days'), datetime('now'), '梁十', '13800001018', 'basic', 2100, 180, 'INV018', 2),
(datetime('now', '-18 days'), datetime('now'), '宋一', '13800001019', 'vip', 16500, 2200, 'INV019', NULL),
(datetime('now', '-16 days'), datetime('now'), '唐二', '13800001020', 'gold', 8200, 820, 'INV020', 19),
(datetime('now', '-14 days'), datetime('now'), '许三', '13800001021', 'silver', 5800, 550, 'INV021', 19),
(datetime('now', '-12 days'), datetime('now'), '韩四', '13800001022', 'basic', 2900, 280, 'INV022', 20),
(datetime('now', '-10 days'), datetime('now'), '邓五', '13800001023', 'gold', 9200, 950, 'INV023', 1),
(datetime('now', '-8 days'), datetime('now'), '冯六', '13800001024', 'silver', 6100, 600, 'INV024', 2),
(datetime('now', '-6 days'), datetime('now'), '曹七', '13800001025', 'basic', 3100, 300, 'INV025', 9),
(datetime('now', '-5 days'), datetime('now'), '彭八', '13800001026', 'gold', 7500, 700, 'INV026', 19),
(datetime('now', '-4 days'), datetime('now'), '曾九', '13800001027', 'silver', 4200, 380, 'INV027', 1),
(datetime('now', '-3 days'), datetime('now'), '肖十', '13800001028', 'basic', 1800, 150, 'INV028', 9),
(datetime('now', '-2 days'), datetime('now'), '田一', '13800001029', 'gold', 8500, 850, 'INV029', 19),
(datetime('now', '-1 days'), datetime('now'), '董二', '13800001030', 'basic', 2400, 220, 'INV030', 1);

-- 技师（8个）
INSERT INTO technicians (created_at, updated_at, name, skills, status, average_rating) VALUES
(datetime('now', '-120 days'), datetime('now'), '陈美丽', '["精油SPA","全身按摩","足疗"]', 0, 4.8),
(datetime('now', '-110 days'), datetime('now'), '王芳', '["推拿","拔罐","刮痧"]', 0, 4.7),
(datetime('now', '-100 days'), datetime('now'), '李娜', '["足底按摩","艾灸","精油护理"]', 0, 4.9),
(datetime('now', '-90 days'), datetime('now'), '刘静', '["SPA","身体护理","面部护理"]', 1, 4.6),
(datetime('now', '-80 days'), datetime('now'), '张丽', '["推拿","按摩","理疗"]', 0, 4.8),
(datetime('now', '-70 days'), datetime('now'), '赵敏', '["足疗","精油","艾灸"]', 0, 4.7),
(datetime('now', '-60 days'), datetime('now'), '孙婷', '["SPA","推拿","拔罐"]', 2, 4.5),
(datetime('now', '-50 days'), datetime('now'), '周洁', '["按摩","足疗","身体护理"]', 0, 4.9);

-- 服务项目（10个）
INSERT INTO service_items (created_at, updated_at, name, duration, price, is_active) VALUES
(datetime('now', '-150 days'), datetime('now'), '全身精油SPA', 90, 388, 1),
(datetime('now', '-150 days'), datetime('now'), '中式推拿', 60, 268, 1),
(datetime('now', '-150 days'), datetime('now'), '足底按摩', 45, 158, 1),
(datetime('now', '-150 days'), datetime('now'), '艾灸护理', 60, 218, 1),
(datetime('now', '-150 days'), datetime('now'), '拔罐刮痧', 45, 188, 1),
(datetime('now', '-150 days'), datetime('now'), '背部舒缓按摩', 60, 258, 1),
(datetime('now', '-150 days'), datetime('now'), '头部肩颈按摩', 45, 198, 1),
(datetime('now', '-150 days'), datetime('now'), '淋巴排毒', 75, 328, 1),
(datetime('now', '-150 days'), datetime('now'), '热石疗法', 90, 458, 1),
(datetime('now', '-150 days'), datetime('now'), '泰式按摩', 90, 398, 1);

-- 商品（12个）
INSERT INTO physical_products (created_at, updated_at, name, stock, retail_price, cost_price, description, is_active, image_url) VALUES
(datetime('now', '-100 days'), datetime('now'), '薰衣草精油50ml', 85, 158, 80, '纯天然精油', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '玫瑰精油30ml', 62, 288, 150, '美容养颜', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '茶树精油50ml', 48, 128, 65, '抗菌消炎', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '按摩精油套装', 35, 398, 200, '6种精油组合', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '足部护理霜100g', 120, 88, 45, '深层滋养', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '身体乳液200ml', 95, 168, 85, '保湿滋润', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '艾灸贴10片装', 150, 58, 28, '缓解疲劳', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '按摩油500ml', 42, 218, 110, '专业按摩用', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '香薰蜡烛套装', 68, 128, 65, '3支装', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '面部精华液30ml', 55, 328, 165, '抗衰老', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '竹炭足贴20片', 180, 68, 35, '改善睡眠', 1, ''),
(datetime('now', '-100 days'), datetime('now'), '按摩刮痧板', 75, 98, 50, '天然玉石', 1, '');

-- 预约记录（近30天，每天4-5条）
-- 这里插入120条预约
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 2, 2, 2, datetime('now', '-30 days', '+12 hours'), datetime('now', '-30 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-30 days'), datetime('now', '-30 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 3, 3, 3, datetime('now', '-30 days', '+14 hours'), datetime('now', '-30 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-30 days'), datetime('now', '-30 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 4, 4, 4, datetime('now', '-30 days', '+16 hours'), datetime('now', '-30 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-30 days'), datetime('now', '-30 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 5, 5, 5, datetime('now', '-30 days', '+18 hours'), datetime('now', '-30 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-30 days'), datetime('now', '-30 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 1, 2, 1, datetime('now', '-29 days', '+12 hours'), datetime('now', '-29 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-29 days'), datetime('now', '-29 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 2, 3, 2, datetime('now', '-29 days', '+14 hours'), datetime('now', '-29 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-29 days'), datetime('now', '-29 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 3, 4, 3, datetime('now', '-29 days', '+16 hours'), datetime('now', '-29 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-29 days'), datetime('now', '-29 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 4, 5, 4, datetime('now', '-29 days', '+18 hours'), datetime('now', '-29 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-29 days'), datetime('now', '-29 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 30, 2, 10, datetime('now', '-28 days', '+12 hours'), datetime('now', '-28 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-28 days'), datetime('now', '-28 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 1, 3, 1, datetime('now', '-28 days', '+14 hours'), datetime('now', '-28 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-28 days'), datetime('now', '-28 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 2, 4, 2, datetime('now', '-28 days', '+16 hours'), datetime('now', '-28 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-28 days'), datetime('now', '-28 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 3, 5, 3, datetime('now', '-28 days', '+18 hours'), datetime('now', '-28 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-28 days'), datetime('now', '-28 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 29, 2, 9, datetime('now', '-27 days', '+12 hours'), datetime('now', '-27 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-27 days'), datetime('now', '-27 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 30, 3, 10, datetime('now', '-27 days', '+14 hours'), datetime('now', '-27 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-27 days'), datetime('now', '-27 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 1, 4, 1, datetime('now', '-27 days', '+16 hours'), datetime('now', '-27 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-27 days'), datetime('now', '-27 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 2, 5, 2, datetime('now', '-27 days', '+18 hours'), datetime('now', '-27 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-27 days'), datetime('now', '-27 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 28, 2, 8, datetime('now', '-26 days', '+12 hours'), datetime('now', '-26 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-26 days'), datetime('now', '-26 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 29, 3, 9, datetime('now', '-26 days', '+14 hours'), datetime('now', '-26 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-26 days'), datetime('now', '-26 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 30, 4, 10, datetime('now', '-26 days', '+16 hours'), datetime('now', '-26 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-26 days'), datetime('now', '-26 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 1, 5, 1, datetime('now', '-26 days', '+18 hours'), datetime('now', '-26 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-26 days'), datetime('now', '-26 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 27, 2, 7, datetime('now', '-25 days', '+12 hours'), datetime('now', '-25 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-25 days'), datetime('now', '-25 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 28, 3, 8, datetime('now', '-25 days', '+14 hours'), datetime('now', '-25 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-25 days'), datetime('now', '-25 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 29, 4, 9, datetime('now', '-25 days', '+16 hours'), datetime('now', '-25 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-25 days'), datetime('now', '-25 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 30, 5, 10, datetime('now', '-25 days', '+18 hours'), datetime('now', '-25 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-25 days'), datetime('now', '-25 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 26, 2, 6, datetime('now', '-24 days', '+12 hours'), datetime('now', '-24 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-24 days'), datetime('now', '-24 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 27, 3, 7, datetime('now', '-24 days', '+14 hours'), datetime('now', '-24 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-24 days'), datetime('now', '-24 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 28, 4, 8, datetime('now', '-24 days', '+16 hours'), datetime('now', '-24 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-24 days'), datetime('now', '-24 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 29, 5, 9, datetime('now', '-24 days', '+18 hours'), datetime('now', '-24 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-24 days'), datetime('now', '-24 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 25, 2, 5, datetime('now', '-23 days', '+12 hours'), datetime('now', '-23 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-23 days'), datetime('now', '-23 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 26, 3, 6, datetime('now', '-23 days', '+14 hours'), datetime('now', '-23 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-23 days'), datetime('now', '-23 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 27, 4, 7, datetime('now', '-23 days', '+16 hours'), datetime('now', '-23 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-23 days'), datetime('now', '-23 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 28, 5, 8, datetime('now', '-23 days', '+18 hours'), datetime('now', '-23 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-23 days'), datetime('now', '-23 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 24, 2, 4, datetime('now', '-22 days', '+12 hours'), datetime('now', '-22 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-22 days'), datetime('now', '-22 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 25, 3, 5, datetime('now', '-22 days', '+14 hours'), datetime('now', '-22 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-22 days'), datetime('now', '-22 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 26, 4, 6, datetime('now', '-22 days', '+16 hours'), datetime('now', '-22 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-22 days'), datetime('now', '-22 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 27, 5, 7, datetime('now', '-22 days', '+18 hours'), datetime('now', '-22 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-22 days'), datetime('now', '-22 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 23, 2, 3, datetime('now', '-21 days', '+12 hours'), datetime('now', '-21 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-21 days'), datetime('now', '-21 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 24, 3, 4, datetime('now', '-21 days', '+14 hours'), datetime('now', '-21 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-21 days'), datetime('now', '-21 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 25, 4, 5, datetime('now', '-21 days', '+16 hours'), datetime('now', '-21 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-21 days'), datetime('now', '-21 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 26, 5, 6, datetime('now', '-21 days', '+18 hours'), datetime('now', '-21 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-21 days'), datetime('now', '-21 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 22, 2, 2, datetime('now', '-20 days', '+12 hours'), datetime('now', '-20 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-20 days'), datetime('now', '-20 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 23, 3, 3, datetime('now', '-20 days', '+14 hours'), datetime('now', '-20 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-20 days'), datetime('now', '-20 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 24, 4, 4, datetime('now', '-20 days', '+16 hours'), datetime('now', '-20 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-20 days'), datetime('now', '-20 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 25, 5, 5, datetime('now', '-20 days', '+18 hours'), datetime('now', '-20 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-20 days'), datetime('now', '-20 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 21, 2, 1, datetime('now', '-19 days', '+12 hours'), datetime('now', '-19 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-19 days'), datetime('now', '-19 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 22, 3, 2, datetime('now', '-19 days', '+14 hours'), datetime('now', '-19 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-19 days'), datetime('now', '-19 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 23, 4, 3, datetime('now', '-19 days', '+16 hours'), datetime('now', '-19 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-19 days'), datetime('now', '-19 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 24, 5, 4, datetime('now', '-19 days', '+18 hours'), datetime('now', '-19 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-19 days'), datetime('now', '-19 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 20, 2, 10, datetime('now', '-18 days', '+12 hours'), datetime('now', '-18 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-18 days'), datetime('now', '-18 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 21, 3, 1, datetime('now', '-18 days', '+14 hours'), datetime('now', '-18 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-18 days'), datetime('now', '-18 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 22, 4, 2, datetime('now', '-18 days', '+16 hours'), datetime('now', '-18 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-18 days'), datetime('now', '-18 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 23, 5, 3, datetime('now', '-18 days', '+18 hours'), datetime('now', '-18 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-18 days'), datetime('now', '-18 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 19, 2, 9, datetime('now', '-17 days', '+12 hours'), datetime('now', '-17 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-17 days'), datetime('now', '-17 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 20, 3, 10, datetime('now', '-17 days', '+14 hours'), datetime('now', '-17 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-17 days'), datetime('now', '-17 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 21, 4, 1, datetime('now', '-17 days', '+16 hours'), datetime('now', '-17 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-17 days'), datetime('now', '-17 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 22, 5, 2, datetime('now', '-17 days', '+18 hours'), datetime('now', '-17 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-17 days'), datetime('now', '-17 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 18, 2, 8, datetime('now', '-16 days', '+12 hours'), datetime('now', '-16 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-16 days'), datetime('now', '-16 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 19, 3, 9, datetime('now', '-16 days', '+14 hours'), datetime('now', '-16 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-16 days'), datetime('now', '-16 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 20, 4, 10, datetime('now', '-16 days', '+16 hours'), datetime('now', '-16 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-16 days'), datetime('now', '-16 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 21, 5, 1, datetime('now', '-16 days', '+18 hours'), datetime('now', '-16 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-16 days'), datetime('now', '-16 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 17, 2, 7, datetime('now', '-15 days', '+12 hours'), datetime('now', '-15 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-15 days'), datetime('now', '-15 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 18, 3, 8, datetime('now', '-15 days', '+14 hours'), datetime('now', '-15 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-15 days'), datetime('now', '-15 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 19, 4, 9, datetime('now', '-15 days', '+16 hours'), datetime('now', '-15 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-15 days'), datetime('now', '-15 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 20, 5, 10, datetime('now', '-15 days', '+18 hours'), datetime('now', '-15 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-15 days'), datetime('now', '-15 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 16, 2, 6, datetime('now', '-14 days', '+12 hours'), datetime('now', '-14 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-14 days'), datetime('now', '-14 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 17, 3, 7, datetime('now', '-14 days', '+14 hours'), datetime('now', '-14 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-14 days'), datetime('now', '-14 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 18, 4, 8, datetime('now', '-14 days', '+16 hours'), datetime('now', '-14 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-14 days'), datetime('now', '-14 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 19, 5, 9, datetime('now', '-14 days', '+18 hours'), datetime('now', '-14 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-14 days'), datetime('now', '-14 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 15, 2, 5, datetime('now', '-13 days', '+12 hours'), datetime('now', '-13 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-13 days'), datetime('now', '-13 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 16, 3, 6, datetime('now', '-13 days', '+14 hours'), datetime('now', '-13 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-13 days'), datetime('now', '-13 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 17, 4, 7, datetime('now', '-13 days', '+16 hours'), datetime('now', '-13 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-13 days'), datetime('now', '-13 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 18, 5, 8, datetime('now', '-13 days', '+18 hours'), datetime('now', '-13 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-13 days'), datetime('now', '-13 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 14, 2, 4, datetime('now', '-12 days', '+12 hours'), datetime('now', '-12 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-12 days'), datetime('now', '-12 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 15, 3, 5, datetime('now', '-12 days', '+14 hours'), datetime('now', '-12 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-12 days'), datetime('now', '-12 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 16, 4, 6, datetime('now', '-12 days', '+16 hours'), datetime('now', '-12 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-12 days'), datetime('now', '-12 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 17, 5, 7, datetime('now', '-12 days', '+18 hours'), datetime('now', '-12 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-12 days'), datetime('now', '-12 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 13, 2, 3, datetime('now', '-11 days', '+12 hours'), datetime('now', '-11 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-11 days'), datetime('now', '-11 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 14, 3, 4, datetime('now', '-11 days', '+14 hours'), datetime('now', '-11 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-11 days'), datetime('now', '-11 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 15, 4, 5, datetime('now', '-11 days', '+16 hours'), datetime('now', '-11 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-11 days'), datetime('now', '-11 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 16, 5, 6, datetime('now', '-11 days', '+18 hours'), datetime('now', '-11 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-11 days'), datetime('now', '-11 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 12, 2, 2, datetime('now', '-10 days', '+12 hours'), datetime('now', '-10 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-10 days'), datetime('now', '-10 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 13, 3, 3, datetime('now', '-10 days', '+14 hours'), datetime('now', '-10 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-10 days'), datetime('now', '-10 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 14, 4, 4, datetime('now', '-10 days', '+16 hours'), datetime('now', '-10 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-10 days'), datetime('now', '-10 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 15, 5, 5, datetime('now', '-10 days', '+18 hours'), datetime('now', '-10 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-10 days'), datetime('now', '-10 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 11, 2, 1, datetime('now', '-9 days', '+12 hours'), datetime('now', '-9 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-9 days'), datetime('now', '-9 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 12, 3, 2, datetime('now', '-9 days', '+14 hours'), datetime('now', '-9 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-9 days'), datetime('now', '-9 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 13, 4, 3, datetime('now', '-9 days', '+16 hours'), datetime('now', '-9 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-9 days'), datetime('now', '-9 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 14, 5, 4, datetime('now', '-9 days', '+18 hours'), datetime('now', '-9 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-9 days'), datetime('now', '-9 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 10, 2, 10, datetime('now', '-8 days', '+12 hours'), datetime('now', '-8 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-8 days'), datetime('now', '-8 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 11, 3, 1, datetime('now', '-8 days', '+14 hours'), datetime('now', '-8 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-8 days'), datetime('now', '-8 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 12, 4, 2, datetime('now', '-8 days', '+16 hours'), datetime('now', '-8 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-8 days'), datetime('now', '-8 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 13, 5, 3, datetime('now', '-8 days', '+18 hours'), datetime('now', '-8 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-8 days'), datetime('now', '-8 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 9, 2, 9, datetime('now', '-7 days', '+12 hours'), datetime('now', '-7 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-7 days'), datetime('now', '-7 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 10, 3, 10, datetime('now', '-7 days', '+14 hours'), datetime('now', '-7 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-7 days'), datetime('now', '-7 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 11, 4, 1, datetime('now', '-7 days', '+16 hours'), datetime('now', '-7 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-7 days'), datetime('now', '-7 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 12, 5, 2, datetime('now', '-7 days', '+18 hours'), datetime('now', '-7 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-7 days'), datetime('now', '-7 days') FROM service_items WHERE id = 2;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 8, 2, 8, datetime('now', '-6 days', '+12 hours'), datetime('now', '-6 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-6 days'), datetime('now', '-6 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 9, 3, 9, datetime('now', '-6 days', '+14 hours'), datetime('now', '-6 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-6 days'), datetime('now', '-6 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 10, 4, 10, datetime('now', '-6 days', '+16 hours'), datetime('now', '-6 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-6 days'), datetime('now', '-6 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 11, 5, 1, datetime('now', '-6 days', '+18 hours'), datetime('now', '-6 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-6 days'), datetime('now', '-6 days') FROM service_items WHERE id = 1;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 7, 2, 7, datetime('now', '-5 days', '+12 hours'), datetime('now', '-5 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-5 days'), datetime('now', '-5 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 8, 3, 8, datetime('now', '-5 days', '+14 hours'), datetime('now', '-5 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-5 days'), datetime('now', '-5 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 9, 4, 9, datetime('now', '-5 days', '+16 hours'), datetime('now', '-5 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-5 days'), datetime('now', '-5 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 10, 5, 10, datetime('now', '-5 days', '+18 hours'), datetime('now', '-5 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-5 days'), datetime('now', '-5 days') FROM service_items WHERE id = 10;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 6, 2, 6, datetime('now', '-4 days', '+12 hours'), datetime('now', '-4 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-4 days'), datetime('now', '-4 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 7, 3, 7, datetime('now', '-4 days', '+14 hours'), datetime('now', '-4 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-4 days'), datetime('now', '-4 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 8, 4, 8, datetime('now', '-4 days', '+16 hours'), datetime('now', '-4 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-4 days'), datetime('now', '-4 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 9, 5, 9, datetime('now', '-4 days', '+18 hours'), datetime('now', '-4 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-4 days'), datetime('now', '-4 days') FROM service_items WHERE id = 9;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 5, 2, 5, datetime('now', '-3 days', '+12 hours'), datetime('now', '-3 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-3 days'), datetime('now', '-3 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 6, 3, 6, datetime('now', '-3 days', '+14 hours'), datetime('now', '-3 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-3 days'), datetime('now', '-3 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 7, 4, 7, datetime('now', '-3 days', '+16 hours'), datetime('now', '-3 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-3 days'), datetime('now', '-3 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 8, 5, 8, datetime('now', '-3 days', '+18 hours'), datetime('now', '-3 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-3 days'), datetime('now', '-3 days') FROM service_items WHERE id = 8;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 4, 2, 4, datetime('now', '-2 days', '+12 hours'), datetime('now', '-2 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-2 days'), datetime('now', '-2 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 5, 3, 5, datetime('now', '-2 days', '+14 hours'), datetime('now', '-2 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-2 days'), datetime('now', '-2 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 6, 4, 6, datetime('now', '-2 days', '+16 hours'), datetime('now', '-2 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-2 days'), datetime('now', '-2 days') FROM service_items WHERE id = 6;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 7, 5, 7, datetime('now', '-2 days', '+18 hours'), datetime('now', '-2 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-2 days'), datetime('now', '-2 days') FROM service_items WHERE id = 7;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 3, 2, 3, datetime('now', '-1 days', '+12 hours'), datetime('now', '-1 days', '+12 hours', '+90 minutes'), 'completed', price, price - 10, datetime('now', '-1 days'), datetime('now', '-1 days') FROM service_items WHERE id = 3;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 4, 3, 4, datetime('now', '-1 days', '+14 hours'), datetime('now', '-1 days', '+14 hours', '+90 minutes'), 'completed', price, price - 20, datetime('now', '-1 days'), datetime('now', '-1 days') FROM service_items WHERE id = 4;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 5, 4, 5, datetime('now', '-1 days', '+16 hours'), datetime('now', '-1 days', '+16 hours', '+90 minutes'), 'completed', price, price - 30, datetime('now', '-1 days'), datetime('now', '-1 days') FROM service_items WHERE id = 5;
INSERT INTO appointments (member_id, tech_id, service_id, start_time, end_time, status, origin_price, actual_price, created_at, updated_at)
SELECT 6, 5, 6, datetime('now', '-1 days', '+18 hours'), datetime('now', '-1 days', '+18 hours', '+90 minutes'), 'completed', price, price - 40, datetime('now', '-1 days'), datetime('now', '-1 days') FROM service_items WHERE id = 6;

-- 为预约创建订单
INSERT INTO orders (member_id, appointment_id, type, quantity, total_amount, actual_paid, status, payment_method, created_at, updated_at)
SELECT member_id, id, 'service', 1, actual_price, actual_price, 'completed', 
  CASE (id % 3) WHEN 0 THEN 'cash' WHEN 1 THEN 'wechat' ELSE 'alipay' END,
  created_at, updated_at
FROM appointments WHERE status = 'completed';

-- 商品订单（30条）
INSERT INTO orders (member_id, product_id, type, quantity, total_amount, actual_paid, status, payment_method, created_at, updated_at)
SELECT 
  (id % 30) + 1,
  id,
  'product',
  2,
  retail_price * 2,
  retail_price * 2,
  'completed',
  CASE (id % 3) WHEN 0 THEN 'cash' WHEN 1 THEN 'wechat' ELSE 'alipay' END,
  datetime('now', '-' || (id * 2) || ' days'),
  datetime('now', '-' || (id * 2) || ' days')
FROM physical_products LIMIT 12;

-- 裂变记录
INSERT INTO fission_logs (inviter_id, invitee_id, commission_amount, order_id, created_at, updated_at)
SELECT m.referrer_id, m.id, o.actual_paid * 0.05, o.id, o.created_at, o.updated_at
FROM members m
JOIN orders o ON o.member_id = m.id
WHERE m.referrer_id IS NOT NULL AND o.status = 'completed';

-- 库存变动
INSERT INTO inventory_logs (product_id, operator_id, change_amount, action_type, before_stock, after_stock, remark, created_at, updated_at)
SELECT id, 1, 100, 'restock', 0, 100, '初始库存', datetime('now', '-30 days'), datetime('now', '-30 days')
FROM physical_products;

