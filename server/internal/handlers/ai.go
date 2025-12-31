package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"smartspa-admin/internal/ai"
	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"
	"smartspa-admin/internal/response"

	"github.com/gin-gonic/gin"
)

// GenerateAIReport 聚合数据并调用 LLM 生成经营分析报告
func GenerateAIReport(c *gin.Context) {
	// 1. 聚合近 30 天的数据
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	// 1.1 营收数据
	var totalRevenue float64
	var orderCount int64
	db.DB.Model(&models.Appointment{}).
		Where("status = ? AND end_time >= ?", "completed", thirtyDaysAgo).
		Count(&orderCount).
		Select("COALESCE(SUM(actual_price), 0)").
		Row().
		Scan(&totalRevenue)

	// 1.2 新增会员
	var newMemberCount int64
	db.DB.Model(&models.Member{}).
		Where("created_at >= ?", thirtyDaysAgo).
		Count(&newMemberCount)

	// 1.3 热门项目 (Top 5)
	type ServiceStat struct {
		Name  string
		Count int
	}
	var topServices []ServiceStat
	db.DB.Model(&models.Appointment{}).
		Select("service_items.name, count(appointments.id) as count").
		Joins("join service_items on service_items.id = appointments.service_id").
		Where("appointments.status = ? AND appointments.end_time >= ?", "completed", thirtyDaysAgo).
		Group("service_items.name").
		Order("count desc").
		Limit(5).
		Scan(&topServices)

	// 1.4 技师表现 (Top 5 by order count)
	type TechStat struct {
		Name  string
		Count int
	}
	var topTechs []TechStat
	db.DB.Model(&models.Appointment{}).
		Select("technicians.name, count(appointments.id) as count").
		Joins("join technicians on technicians.id = appointments.tech_id").
		Where("appointments.status = ? AND appointments.end_time >= ?", "completed", thirtyDaysAgo).
		Group("technicians.name").
		Order("count desc").
		Limit(5).
		Scan(&topTechs)

	// 2. 构建 Prompt
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("以下是店铺近 30 天的经营数据：\n"))
	sb.WriteString(fmt.Sprintf("- 总营收: ¥%.2f\n", totalRevenue))
	sb.WriteString(fmt.Sprintf("- 完成订单数: %d\n", orderCount))
	sb.WriteString(fmt.Sprintf("- 新增会员数: %d\n", newMemberCount))

	sb.WriteString("\n热门服务项目 (Top 5):\n")
	for _, s := range topServices {
		sb.WriteString(fmt.Sprintf("- %s: %d 单\n", s.Name, s.Count))
	}

	sb.WriteString("\n技师工作量排行 (Top 5):\n")
	for _, t := range topTechs {
		sb.WriteString(fmt.Sprintf("- %s: %d 单\n", t.Name, t.Count))
	}

	sb.WriteString("\n请根据以上数据，分析店铺经营状况，指出亮点与不足，并给出针对性的营销建议和项目调整建议。")

	// 3. 调用 AI 服务
	client := ai.NewLLMClient()
	report, err := client.GenerateAnalysis(sb.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to generate AI report: "+err.Error(), nil))
		return
	}

	// 4. 返回结果
	c.JSON(http.StatusOK, response.Success(gin.H{
		"report": report,
		"raw_data": gin.H{
			"revenue":      totalRevenue,
			"new_members":  newMemberCount,
			"top_services": topServices,
		},
	}, "Report generated successfully"))
}
