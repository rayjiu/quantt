package model

import "time"

// SecFundFlow represents the sec_fund_flow table in the database
type SecFundFlow struct {
	SecCode        string    `gorm:"type:varchar;column:sec_code"`                      // 证券代码
	MainBuyerNetIn float64   `gorm:"type:float8;column:main_buyer_netin"`               // 主力买入净流入
	MainBuyerRatio float64   `gorm:"type:float8;column:main_buyer_ratio"`               // 主力买入比例
	L1NetIn        float64   `gorm:"type:varchar;column:l1_netin"`                      // L1净流入
	L1NetInRatio   float32   `gorm:"type:float4;column:l1_netin_ratio"`                 // L1净流入比例
	L2NetIn        float64   `gorm:"type:float8;column:l2_netin"`                       // L2净流入
	L2NetInRatio   float32   `gorm:"type:float4;column:l2_netin_ratio"`                 // L2净流入比例
	L3NetIn        float64   `gorm:"type:float8;column:l3_netin"`                       // L3净流入
	L3NetInRatio   float32   `gorm:"type:float4;column:l3_netin_ratio"`                 // L3净流入比例
	L4NetInt       float64   `gorm:"type:float8;column:l4_netint"`                      // L4净流入
	L4NetRatio     float64   `gorm:"type:float8;column:l4_net_ratio"`                   // L4净流入比例
	MarketDate     int32     `gorm:"type:int4;column:market_date"`                      // 市场日期，存储为整数格式
	UpdateTime     time.Time `gorm:"type:timestamptz;column:update_time;default:now()"` // 更新时间，默认当前时间
}

// TableName specifies the table name for SecFundFlow
func (SecFundFlow) TableName() string {
	return "sec_fund_flow"
}
