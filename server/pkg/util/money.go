package util

import (
	"math"

	"server/pkg/config"
)

// RoundMoney 将金额四舍五入保留两位小数
// val: 金额（单位：元）
// 返回值: 处理后的金额（单位：元）
func RoundMoney(val float64) float64 {
	return math.Round(val*100) / 100
}

// CalculateRate 计算按比例的金额（如折扣、佣金），并保留两位小数
// amount: 原始金额（单位：元）
// rate: 比例（例如 0.8 表示八折，0.1 表示 10% 佣金）
// 返回值: 计算后的金额（单位：元）
func CalculateRate(amount float64, rate float64) float64 {
	return RoundMoney(amount * rate)
}

// ToCents 将金额转换为分（整数），用于精确计算或存储
// amount: 金额（单位：元）
// 返回值: 金额（单位：分）
func ToCents(amount float64) int64 {
	// 先四舍五入避免浮点数精度问题 (e.g. 19.999999 -> 20.00)
	return int64(math.Round(amount * 100))
}

// CentsToYuan 将分（整数）转换为元（浮点数）
// cents: 金额（单位：分）
// 返回值: 金额（单位：元）
func CentsToYuan(cents int64) float64 {
	return float64(cents) / 100.0
}

// CalculateMemberLevel 根据年度消费额计算会员等级
// consumption: 年度消费总额（单位：元）
// 返回值: 对应的会员等级字符串 (platinum/gold/silver/basic)
func CalculateMemberLevel(consumption float64) string {
	switch {
	case consumption > config.GlobalMemberUpgrade.Platinum:
		return "platinum"
	case consumption > config.GlobalMemberUpgrade.Gold:
		return "gold"
	case consumption > config.GlobalMemberUpgrade.Silver:
		return "silver"
	default:
		return "basic"
	}
}
