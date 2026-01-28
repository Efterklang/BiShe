package repo

import (
	"smartspa-admin/internal/db"
	"smartspa-admin/internal/models"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

type ScheduleRepo struct{}

var Schedule = &ScheduleRepo{}

// GetSchedules 获取排班列表
func (r *ScheduleRepo) GetSchedules(startDate, endDate time.Time, techIDs []string) ([]models.Schedule, error) {
	var schedules []models.Schedule
	query := db.DB.Model(&models.Schedule{}).
		Preload("Technician").
		Where("date >= ? AND date <= ?", datatypes.Date(startDate), datatypes.Date(endDate))

	if len(techIDs) > 0 {
		query = query.Where("tech_id IN ?", techIDs)
	}

	if err := query.Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// BatchUpsertSchedules 批量更新排班
func (r *ScheduleRepo) BatchUpsertSchedules(schedules []map[string]interface{}) error {
	return db.DB.Model(&models.Schedule{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tech_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_available", "time_slots"}),
	}).Create(&schedules).Error
}

// GetByDate 获取某天的所有排班
func (r *ScheduleRepo) GetByDate(date string) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := db.DB.Where("date = ?", date).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetUnavailableTechIDs 获取某天不可用的技师ID集合
func (r *ScheduleRepo) GetUnavailableTechIDs(date string) (map[uint]bool, error) {
	var schedules []models.Schedule
	if err := db.DB.Select("tech_id", "is_available").Where("date = ?", date).Find(&schedules).Error; err != nil {
		return nil, err
	}

	unavailableMap := make(map[uint]bool)
	for _, s := range schedules {
		if !s.IsAvailable {
			unavailableMap[s.TechID] = true
		}
	}
	return unavailableMap, nil
}

// GetByTechAndDate 获取指定技师指定日期的排班
func (r *ScheduleRepo) GetByTechAndDate(techID string, date string) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := db.DB.Where("tech_id = ? AND date = ?", techID, date).First(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

// GetAvailableTechs 获取指定时间段available的技师
func (r *ScheduleRepo) GetAvailableTechs(date string, startTime, endTime time.Time) ([]models.Technician, error) {
	var techs []models.Technician

	// 1. 查找当天不可用的技师 (请假/休息)
	unavailableSub := db.DB.Model(&models.Schedule{}).
		Select("tech_id").
		Where("date = ? AND is_available = ?", date, false)

	// 2. 查找该时间段有预约冲突的技师 (预约状态 != cancelled/confirmed)
	// 冲突条件: 预约开始时间 < 查询结束时间 AND 预约结束时间 > 查询开始时间
	busySub := db.DB.Model(&models.Appointment{}).
		Select("tech_id").
		Where("status NOT IN (?)", []string{"cancelled", "complete"}).
		Where("start_time < ? AND end_time > ?", endTime, startTime)

	// 3. 查询不在上述两个集合中的技师
	if err := db.DB.Model(&models.Technician{}).
		Where("id NOT IN (?)", unavailableSub).
		Where("id NOT IN (?)", busySub).
		Find(&techs).Error; err != nil {
		return nil, err
	}
	return techs, nil
}
