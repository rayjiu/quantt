package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrendKData struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	MarketDate     uint32             `bson:"market_date,omitempty"`     //市场日期
	StockId        string             `bson:"stock_id,omitempty"`        //市场类型
	MarketType     uint32             `bson:"market_type,omitempty"`     //市场类型
	YesterdayClose float64            `bson:"yesterday_close,omitempty"` //昨收价
	OpenPrice      float64            `bson:"open_price,omitempty"`      //开盘价
	ClosePrice     float64            `bson:"close_price,omitempty"`     //收盘价
	TrendItem      []TrendItem        `bson:"trend_item,omitempty"`      //分时数据
	TrendType      uint32             `bson:"trend_type,omitempty"`      //分时类型: 1, 5, 30 ,60, 120 分钟
	CandleMode     uint32             `bson:"candle_mode,omitempty"`     //复权类型
	Multiple       int32              `bson:"multiple,omitempty"`        //返回给客户端增大的倍数的系数，比如n=2， 返回给客户端的内容扩大10^2= 100倍
	CreateTime     time.Time          `bson:"create_time,omitempty"`     //创建时间
	UpdateTime     time.Time          `bson:"update_time,omitempty"`     //修改时间
}

type TrendItem struct { //分时信息
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	MarketDate       uint32             `bson:"market_date,omitempty"`        //市场日期
	MarketTime       uint32             `bson:"market_time,omitempty"`        //市场时间
	MarketType       uint32             `bson:"market_type,omitempty"`        //市场类型
	TradeStat        uint32             `bson:"trade_stat,omitempty"`         //交易状态
	StockId          string             `bson:"stock_id,omitempty"`           //股票代码
	LineNo           int32              `bson:"line_no,omitempty"`            //线号
	LastPrice        float64            `bson:"last_price,omitempty"`         //最新价
	AvgPrice         float64            `bson:"avg_price,omitempty"`          //均价
	Position         uint32             `bson:"position,omitempty"`           //持仓量
	TradeAmount      float64            `bson:"trade_val,omitempty"`          //这一分钟内的成交额
	TradeVolume      uint64             `bson:"trade_amount,omitempty"`       //这一分钟内的成交量
	TotalTradeAmount float64            `bson:"total_trade_val,omitempty"`    //总的成交额
	TotalTradeVolume uint64             `bson:"total_trade_amount,omitempty"` //总的成交额
	Open             float64            `bson:"open"`                         //开盘价
	Close            float64            `bson:"close"`                        // 收盘价
	High             float64            `bson:"high"`                         // 最高价
	Low              float64            `bson:"low"`                          // 最低价
	TradeTime        uint32             `bson:"trade_time,omitempty"`         //交易时间
	Multiple         int32              `bson:"multiple,omitempty"`           // 返回给客户端增大的倍数的系数，比如n=2， 返回给客户端的内容扩大10^2= 100倍
	CreateTime       time.Time          `bson:"create_time,omitempty"`        //创建时间
	UpdateTime       time.Time          `bson:"update_time,omitempty"`        //修改时间
	ExchangeRatio    float64            `bson:"exchange_ratio,omitempty"`     //换手率
}
