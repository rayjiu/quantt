package model

import (
	"time"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

// Sector 表示 sec_base_info 表的模型
type Sector struct {
	SectorID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"` // 主键字段
	SecCode    string    `gorm:"type:varchar(255);index"`                        // sec_code 字段，建立索引
	SecType    int16     `gorm:"type:smallint;index"`                            // sec_type 字段，建立索引
	SecName    string    `gorm:"type:varchar(255);"`                             // sec_name 字段
	UpdateTime time.Time `gorm:"type:timestamptz;"`
}

func (Sector) TableName() string {
	return "sec_base_info"
}
