package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BaseModel replaces gorm.Model with JSON tags
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Member represents a customer profile with referral metadata.
type Member struct {
	BaseModel
	Name                   string  `gorm:"size:64;not null" json:"name"`
	Phone                  string  `gorm:"size:32;uniqueIndex;not null" json:"phone"`
	Level                  string  `gorm:"size:32;default:basic" json:"level"`
	YearlyTotalConsumption float64 `gorm:"type:decimal(12,2);default:0" json:"yearly_total_consumption"`
	Balance                float64 `gorm:"type:decimal(12,2);default:0" json:"balance"`
	InvitationCode         string  `gorm:"size:32;uniqueIndex" json:"invitation_code"`
	ReferrerID             *uint   `json:"referrer_id"`
}

// Technician holds skill tags and availability state.
type Technician struct {
	BaseModel
	Name          string         `gorm:"size:64;not null" json:"name"`
	Skills        datatypes.JSON `gorm:"type:json" json:"skills"`
	Status        int            `gorm:"default:0" json:"status"` // 0:free, 1:booked, 2:leave
	AverageRating float32        `gorm:"type:decimal(3,2);default:0" json:"average_rating"`
}

// ServiceItem describes a spa service with price and duration.
type ServiceItem struct {
	BaseModel
	Name     string  `gorm:"size:64;not null" json:"name"`
	Duration int     `gorm:"not null" json:"duration"` // minutes
	Price    float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive bool    `gorm:"default:true" json:"is_active"`
}

// Appointment captures booking details and pricing.
type Appointment struct {
	BaseModel
	MemberID    uint        `gorm:"index;not null" json:"member_id"`
	Member      Member      `gorm:"foreignKey:MemberID" json:"member"`
	TechID      uint        `gorm:"index;not null" json:"tech_id"`
	Technician  Technician  `gorm:"foreignKey:TechID" json:"technician"`
	ServiceID   uint        `gorm:"index;not null" json:"service_id"`
	ServiceItem ServiceItem `gorm:"foreignKey:ServiceID" json:"service_item"`
	StartTime   time.Time   `gorm:"index;not null" json:"start_time"`
	EndTime     time.Time   `gorm:"index;not null" json:"end_time"`
	Status      string      `gorm:"size:24;default:'pending'" json:"status"` // pending/completed/waitlist/cancelled
	OriginPrice float64     `gorm:"type:decimal(10,2);not null" json:"origin_price"`
	ActualPrice float64     `gorm:"type:decimal(10,2);not null" json:"actual_price"`
}

// Schedule represents a technician's daily availability and booked slots.
type Schedule struct {
	BaseModel
	TechID      uint           `gorm:"index;not null" json:"tech_id"`
	Technician  Technician     `gorm:"foreignKey:TechID" json:"technician,omitempty"`
	Date        time.Time      `gorm:"index;not null" json:"date"`
	TimeSlots   datatypes.JSON `gorm:"type:json" json:"time_slots"`
	IsAvailable bool           `gorm:"default:true" json:"is_available"`
}

// FissionLog stores commission payouts for referral fission events.
type FissionLog struct {
	BaseModel
	InviterID        uint    `gorm:"index;not null" json:"inviter_id"`
	InviteeID        uint    `gorm:"index;not null" json:"invitee_id"`
	CommissionAmount float64 `gorm:"type:decimal(10,2);not null" json:"commission_amount"`
	OrderID          *uint   `gorm:"index" json:"order_id"`
}
