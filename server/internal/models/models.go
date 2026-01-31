package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User represents a system user with role-based access control.
type User struct {
	BaseModel
	Username     string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Role         string `gorm:"size:32;not null" json:"role"` // "manager" or "operator"
	IsActive     bool   `gorm:"default:true" json:"is_active"`
}

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
	AvatarURL     string         `gorm:"size:255" json:"avatar_url"` // 头像链接
	Skills        datatypes.JSON `gorm:"type:json" json:"skills"`
	Status        int            `gorm:"default:0" json:"status"` // 0:free, 1:booked, 2:leave
	AverageRating float32        `gorm:"type:decimal(3,2);default:0" json:"average_rating"`
	Reason        string         `gorm:"-" json:"reason,omitempty"` // skill_mismatch, leave, busy
}

// ServiceProduct describes a spa service with price and duration.
type ServiceProduct struct {
	BaseModel
	Name     string  `gorm:"size:64;not null" json:"name"`
	Duration int     `gorm:"not null" json:"duration"` // minutes
	Price    float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive bool    `gorm:"default:true" json:"is_active"`
	ImageURL string  `gorm:"size:255" json:"image_url"` // 服务图片
}

// Appointment captures booking details and pricing.
type Appointment struct {
	BaseModel
	MemberID       uint           `gorm:"index;not null" json:"member_id"`
	Member         Member         `gorm:"foreignKey:MemberID" json:"member"`
	TechID         uint           `gorm:"index;not null" json:"tech_id"`
	Technician     Technician     `gorm:"foreignKey:TechID" json:"technician"`
	ServiceID      uint           `gorm:"index;not null" json:"service_id"`
	ServiceProduct ServiceProduct `gorm:"foreignKey:ServiceID" json:"service_item"`
	StartTime      time.Time      `gorm:"index;not null" json:"start_time"`
	EndTime        time.Time      `gorm:"index;not null" json:"end_time"`
	Status         string         `gorm:"size:24;default:'pending'" json:"status"` // pending/completed/waitlist/cancelled
	OriginPrice    float64        `gorm:"type:decimal(10,2);not null" json:"origin_price"`
	ActualPrice    float64        `gorm:"type:decimal(10,2);not null" json:"actual_price"`
	PaymentMethod  string         `gorm:"size:32" json:"payment_method"`                    // balance/cash/mixed
	PaidBalance    float64        `gorm:"type:decimal(10,2);default:0" json:"paid_balance"` // 余额支付金额
	PaidCash       float64        `gorm:"type:decimal(10,2);default:0" json:"paid_cash"`    // 现金支付金额
}

// Schedule represents a technician's daily availability
type Schedule struct {
	BaseModel
	TechID      uint           `gorm:"uniqueIndex:idx_tech_date;not null" json:"tech_id"`
	Technician  Technician     `gorm:"foreignKey:TechID" json:"technician"`
	Date        datatypes.Date `gorm:"uniqueIndex:idx_tech_date;not null" json:"date"`
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

// PhysicalProduct represents physical products for sale in the store.
type PhysicalProduct struct {
	BaseModel
	Name        string  `gorm:"size:128;not null" json:"name"`
	Stock       int     `gorm:"not null;default:0" json:"stock"`                         // 库存数量
	RetailPrice float64 `gorm:"type:decimal(10,2);not null" json:"retail_price"`         // 零售价
	CostPrice   float64 `gorm:"type:decimal(10,2);not null;default:0" json:"cost_price"` // 进货价
	Description string  `gorm:"size:500" json:"description"`                             // 商品描述
	IsActive    bool    `gorm:"default:true" json:"is_active"`                           // 是否上架
	ImageURL    string  `gorm:"size:255" json:"image_url"`                               // 商品图片
}

// InventoryLog records all inventory changes for physical products.
type InventoryLog struct {
	BaseModel
	ProductID    uint            `gorm:"index;not null" json:"product_id"`
	Product      PhysicalProduct `gorm:"foreignKey:ProductID" json:"product"`
	OperatorID   uint            `gorm:"index;not null" json:"operator_id"` // 操作员ID
	Operator     User            `gorm:"foreignKey:OperatorID" json:"operator"`
	ChangeAmount int             `gorm:"not null" json:"change_amount"`       // 变动数量（正数为入库，负数为出库）
	ActionType   string          `gorm:"size:32;not null" json:"action_type"` // "restock"(到货), "sale"(销售), "adjustment"(纠错)
	BeforeStock  int             `gorm:"not null" json:"before_stock"`        // 变动前库存
	AfterStock   int             `gorm:"not null" json:"after_stock"`         // 变动后库存
	OrderID      *uint           `gorm:"index" json:"order_id"`               // 关联订单ID（销售时）
	Remark       string          `gorm:"size:255" json:"remark"`              // 备注
}
