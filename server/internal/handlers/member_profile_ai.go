package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"server/internal/ai"
	"server/internal/db"
	"server/internal/models"
	"server/internal/response"

	"github.com/gin-gonic/gin"
)

func GenerateMemberAIProfile(c *gin.Context) {
	idStr := c.Param("id")
	memberID64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid member id", nil))
		return
	}
	memberID := uint(memberID64)

	var member models.Member
	if err := db.DB.First(&member, memberID).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Member not found", err.Error()))
		return
	}

	type Totals struct {
		TotalSpent    float64
		OrderCount    int64
		ServiceSpent  float64
		ServiceCount  int64
		PhysicalSpent float64
		PhysicalCount int64
	}
	var totals Totals
	db.DB.Model(&models.Order{}).Where("member_id = ?", memberID).Count(&totals.OrderCount)
	db.DB.Model(&models.Order{}).Where("member_id = ?", memberID).Select("COALESCE(SUM(paid_amount), 0)").Scan(&totals.TotalSpent)
	db.DB.Model(&models.Order{}).Where("member_id = ? AND order_type = ?", memberID, "service").Count(&totals.ServiceCount)
	db.DB.Model(&models.Order{}).Where("member_id = ? AND order_type = ?", memberID, "service").Select("COALESCE(SUM(paid_amount), 0)").Scan(&totals.ServiceSpent)
	db.DB.Model(&models.Order{}).Where("member_id = ? AND order_type = ?", memberID, "physical").Count(&totals.PhysicalCount)
	db.DB.Model(&models.Order{}).Where("member_id = ? AND order_type = ?", memberID, "physical").Select("COALESCE(SUM(paid_amount), 0)").Scan(&totals.PhysicalSpent)

	var lastOrderAt time.Time
	db.DB.Model(&models.Order{}).Where("member_id = ?", memberID).Select("MAX(created_at)").Scan(&lastOrderAt)

	type TopItem struct {
		Name        string  `json:"name"`
		OrderCount  int64   `json:"order_count"`
		TotalAmount float64 `json:"total_amount"`
	}
	var topServices []TopItem
	db.DB.Table("orders").
		Select("service_products.name as name, COUNT(orders.id) as order_count, COALESCE(SUM(orders.paid_amount), 0) as total_amount").
		Joins("JOIN appointments ON appointments.id = orders.appointment_id").
		Joins("JOIN service_products ON service_products.id = appointments.service_id").
		Where("orders.member_id = ? AND orders.order_type = ?", memberID, "service").
		Group("service_products.name").
		Order("order_count DESC").
		Limit(5).
		Scan(&topServices)

	var topProducts []TopItem
	db.DB.Table("orders").
		Select("physical_products.name as name, COUNT(orders.id) as order_count, COALESCE(SUM(orders.paid_amount), 0) as total_amount").
		Joins("JOIN inventory_logs ON inventory_logs.id = orders.inventory_log_id").
		Joins("JOIN physical_products ON physical_products.id = inventory_logs.product_id").
		Where("orders.member_id = ? AND orders.order_type = ?", memberID, "physical").
		Group("physical_products.name").
		Order("order_count DESC").
		Limit(5).
		Scan(&topProducts)

	var sb strings.Builder
	sb.WriteString("请基于以下会员数据，生成一份 Markdown 格式的用户画像与运营建议。\n")
	sb.WriteString("输出要求：\n- 结构清晰（分标题）\n- 画像总结、消费习惯、偏好推断、风险点、可执行建议（复购/唤醒/加购/拉新）\n- 结尾给出 3 条可落地的门店动作\n\n")
	sb.WriteString(fmt.Sprintf("会员信息：\n- ID: %d\n- 姓名: %s\n- 手机: %s\n- 等级: %s\n- 余额: %.2f\n- 年消费累计: %.2f\n\n", member.ID, member.Name, member.Phone, member.Level, member.Balance, member.YearlyTotalConsumption))
	sb.WriteString(fmt.Sprintf("订单概览：\n- 总订单数: %d\n- 总消费: %.2f\n- 服务订单: %d 单 / %.2f\n- 商品订单: %d 单 / %.2f\n", totals.OrderCount, totals.TotalSpent, totals.ServiceCount, totals.ServiceSpent, totals.PhysicalCount, totals.PhysicalSpent))
	if !lastOrderAt.IsZero() {
		sb.WriteString(fmt.Sprintf("- 最近一次下单: %s\n", lastOrderAt.Format("2006-01-02 15:04:05")))
	}
	sb.WriteString("\n")

	if len(topServices) > 0 {
		sb.WriteString("Top 服务项目：\n")
		for _, it := range topServices {
			sb.WriteString(fmt.Sprintf("- %s: %d 单 / %.2f\n", it.Name, it.OrderCount, it.TotalAmount))
		}
		sb.WriteString("\n")
	}
	if len(topProducts) > 0 {
		sb.WriteString("Top 商品：\n")
		for _, it := range topProducts {
			sb.WriteString(fmt.Sprintf("- %s: %d 单 / %.2f\n", it.Name, it.OrderCount, it.TotalAmount))
		}
		sb.WriteString("\n")
	}

	systemPrompt := "你是一位资深的会员运营与CRM分析顾问，擅长从消费记录中提炼用户画像与可执行营销策略。"
	client := ai.NewLLMClient()
	profile, err := client.GenerateText(systemPrompt, sb.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to generate AI profile: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"profile": profile,
	}, ""))
}

