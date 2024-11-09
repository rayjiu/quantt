package model

// SecQuote represents the sec_quote table in the database
type SecQuote struct {
	SecCode     string  `gorm:"type:varchar;column:sec_code"`     // 股票代码
	LastPrice   float64 `gorm:"type:float8;column:last_price"`    // 最新价格
	ChgRatio    float64 `gorm:"type:float8;column:chg_ratio"`     // 涨跌幅
	ExchgRatio  float64 `gorm:"type:float8;column:exchg_ratio"`   // 换手率
	TotalMktCap float64 `gorm:"type:float8;column:total_mkt_cap"` // 总市值
	MarketType  int16   `gorm:"type:int2;column:market_type"`     // 市场类型
}

// TableName specifies the table name for SecQuote
func (SecQuote) TableName() string {
	return "sec_quote"
}
