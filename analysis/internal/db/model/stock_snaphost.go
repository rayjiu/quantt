package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockSnapshot struct { //股票快照信息stock
	Id                     primitive.ObjectID `bson:"_id,omitempty"`
	MarketType             uint32             `bson:"market_type,omitempty"`               //市场类型
	StockId                string             `bson:"stock_id,omitempty"`                  //股票代码
	MarketDate             uint32             `bson:"market_date,omitempty"`               //交易日期
	TradeTime              uint32             `bson:"trade_time,omitempty"`                //交易时间
	TotalMins              uint32             `bson:"total_mins,omitempty"`                //总交易时间
	CurrDisplayPrice       float64            `bson:"curr_display_price,omitempty"`        //当前显示价格
	CurrCollectPrice       float64            `bson:"curr_collect_price,omitempty"`        //当前集合竞价价格
	YesterdayClosePrice    float64            `bson:"yesterday_close_price,omitempty"`     //昨收价
	LowestPrice            float64            `bson:"lowest_price,omitempty"`              //最低价
	LastPrice              float64            `bson:"last_price,omitempty"`                //最新价
	OpenPrice              float64            `bson:"open_price,omitempty"`                //开盘价
	ClosePrice             float64            `bson:"close_price,omitempty"`               //收盘价
	HighestPrice           float64            `bson:"highest_price,omitempty"`             //最高价
	TradeVal               float64            `bson:"trade_val,omitempty"`                 //成交额
	TradeAmount            uint64             `bson:"trade_amount,omitempty"`              //成交量
	TradeStatus            int32              `bson:"trade_status,omitempty"`              //交易状态
	Buy1Price              float64            `bson:"b1p, omitempty"`                      // 买1价
	Buy1Amt                uint32             `bson:"b1a, omitempty"`                      // 买1量
	Sell1Price             float64            `bson:"s1p, omitempty"`                      // 卖1价
	Sell1Amt               uint32             `bson:"s1a, omitempty"`                      // 卖1量
	ChangePrice            float64            `bson:"chg_price,omitempty"`                 // 涨跌额，最新价减去昨收价
	ChangeRatio            float64            `bson:"chg_ratio,omitempty"`                 // 涨跌幅
	YearChangeRatio        float64            `bson:"y_chg_ratio,omitempty"`               // 今年以来涨跌幅
	Day20ChangeRatio       float64            `bson:"d20_chg_ratio,omitempty"`             // 近20涨跌幅
	TotalVal               float64            `bson:"total_val,omitempty"`                 //总市值
	FlowTotalVal           float64            `bson:"flow_total_val,omitempty"`            //流通总市值
	FreeFlowTotalVal       float64            `bson:"free_flow_total_val,omitempty"`       //自由流通总市值
	YesterTotalVal         float64            `bson:"yes_flow_tval,omitempty"`             //昨日流通总市值，用于计算板块涨跌幅
	LineNo                 int32              `bson:"line_no,omitempty"`                   //线号
	ExchangeRatio          float64            `bson:"exchange_ratio,omitempty"`            //换手率
	RealExchangeRatio      float64            `bson:"real_exchange_ratio,omitempty"`       //真实换手率
	TradeNum               uint32             `bson:"trade_num,omitempty"`                 //成交笔数
	YesterdayCloseIOPV     float64            `bson:"yesterday_close_iopv,omitempty"`      //昨收基金净值
	Iopv                   float64            `bson:"iopv,omitempty"`                      //基金净值
	DynPERate              float64            `bson:"dyn_pe_rate,omitempty"`               //动态市盈率
	PeRate                 float64            `bson:"pe_rate,omitempty"`                   //静态市盈率
	PeTTM                  float64            `bson:"pe_ttm,omitempty"`                    //滚动市盈率
	PbRate                 float64            `bson:"pb_rate,omitempty"`                   //市净率
	RiseCount              uint32             `bson:"rise_count,omitempty"`                //涨家数
	FallCount              uint32             `bson:"fall_count,omitempty"`                //跌家数
	UnchangeCount          uint32             `bson:"unchange_count,omitempty"`            //平家数
	AvgPrice               float64            `bson:"avg_price,omitempty"`                 //均价
	InPlate                uint32             `bson:"in_plate,omitempty"`                  //内盘
	OutPlate               uint32             `bson:"out_plate,omitempty"`                 //外盘
	EntrustRate            float64            `bson:"entrust_rate,omitempty"`              //委比
	AmountRate             float64            `bson:"amount_rate,omitempty"`               //量比
	AfterTradeDate         uint32             `bson:"after_trade_date,omitempty"`          //盘后交易日期
	AfterTradeTime         uint32             `bson:"after_trade_time,omitempty"`          //盘后交易时间
	AfterTradeAmount       uint32             `bson:"after_trade_amount,omitempty"`        //盘后交易量
	AfterTradeVal          float64            `bson:"after_trade_val,omitempty"`           //盘后交易额
	YesterdaySettlePrice   float64            `bson:"yesterday_settle_price,omitempty"`    //昨日结算价
	YesterdayPosition      uint32             `bson:"yesterday_position,omitempty"`        //昨日持仓量
	SettlePrice            float64            `bson:"settle_price,omitempty"`              //结算价
	Position               uint32             `bson:"position,omitempty"`                  //持仓量
	TotalEntrustBuyAmount  float64            `bson:"total_entrust_buy_amount,omitempty"`  //总委买比
	TotalEntrustSaleAmount float64            `bson:"total_entrust_sale_amount,omitempty"` //总委卖比
	WeightAvgBuyPrice      float64            `bson:"weight_avg_buy_price,omitempty"`      //加权平均买价
	WeightAvgSalePrice     float64            `bson:"weight_avg_sale_price,omitempty"`     //加权平均卖价
	AfterTradePrice        float64            `bson:"after_trade_price,omitempty"`         //盘后成交价
	PremiumRate            float64            `bson:"premium_rate,omitempty"`              // 溢价率
	EquitySwapPrice        float64            `bson:"swap_price,omitempty"`                // 转股价
	EquitySwapValue        float64            `bson:"swap_value,omitempty"`                // 转鼓价值
	CallAucTradeVal        float64            `bson:"call_auc_tradeval,omitempty"`         // 集合竞价成交额
	HHisPx                 string             `bson:"hhis_px,omitempty"`                   //历史最高价字符串, "前复权,后复权,不复权"
	H52WsPx                string             `bson:"h52wsPx,omitempty"`                   // 52周最高价字符串, "前复权,后复权,不复权"
	L52WsPx                string             `bson:"l52wsPx,omitempty"`                   // 52周最低价字符串, "前复权,后复权,不复权"
	StockDividendRatio     float64            `bson:"stock_divi_ratio,omitempty"`          // 股息率
	ExpectedRevenue        float64            `bson:"exp_revenue,omitempty"`               // 预期收益
	AvailableUseDate       uint32             `bson:"a_use_date,omitempty"`                // 可用日期
	// GearPrices             []GearPrice        `bson:"gear_prices,omitempty"`               //买卖五档
	UpOrDownLimit int32     `bson:"ul"`                    // -1 跌停，1 涨停
	Multiple      int32     `bson:"multiple,omitempty"`    // 返回给客户端增大的倍数的系数，比如n=2， 返回给客户端的内容扩大10^2= 100倍
	CreateTime    time.Time `bson:"create_time,omitempty"` //创建时间
	UpdateTime    time.Time `bson:"update_time,omitempty"` //修改时间
}
