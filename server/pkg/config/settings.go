// Package config 提供项目全局可共享的业务配置管理
package config

import "time"

type BusinessHours struct {
	OpenHour     int            // 开门小时，24小时制（如10表示上午10点）
	OpenMinute   int            // 开门分钟（如0表示整点）
	CloseHour    int            // 打烊小时，24小时制（如22表示晚上10点）
	CloseMinute  int            // 打烊分钟（如0表示整点）
	SlotInterval time.Duration  // 时间槽生成间隔，如30*time.Minute表示每30分钟一个时段
	TimeLocation *time.Location // 营业时间对应的时区，统一时间计算标准
}

var GlobalBusinessHours = BusinessHours{
	OpenHour:     10,
	OpenMinute:   0,
	CloseHour:    22,
	CloseMinute:  0,
	SlotInterval: 30 * time.Minute, // 替换你代码里硬编码的30分钟间隔
	TimeLocation: time.UTC,         // 兼容你原有代码的时区，后续可改成Asia/Shanghai
}

type MemberDiscountConfig struct {
	Platinum float64 // 白金会员折扣率 (e.g. 0.8 for 20% off)
	Gold     float64 // 金牌会员折扣率
	Silver   float64 // 银牌会员折扣率
	Basic    float64 // 普通会员折扣率 (usually 1.0)
}

var GlobalMemberDiscount = MemberDiscountConfig{
	Platinum: 0.8,
	Gold:     0.9,
	Silver:   0.95,
	Basic:    1.0,
}

type CommissionConfig struct {
	ReferralRate float64 // 推荐佣金比例 (e.g. 0.1 for 10%)
}

var GlobalCommission = CommissionConfig{
	ReferralRate: 0.1,
}

type MemberUpgradeThresholds struct {
	Platinum float64 // 白金会员升级阈值（年度消费额，单位：元）
	Gold     float64 // 金卡会员升级阈值（年度消费额，单位：元）
	Silver   float64 // 银卡会员升级阈值（年度消费额，单位：元）
}

var GlobalMemberUpgrade = MemberUpgradeThresholds{
	Platinum: 10000, // 10,000元
	Gold:     5000,  // 5,000元
	Silver:   1000,  // 1,000元
}
