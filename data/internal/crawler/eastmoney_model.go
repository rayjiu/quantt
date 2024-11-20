package crawler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rayjiu/quantt/data/internal/db/model"
)

type Diff struct {
	F1   int     `json:"f1"`
	F2   float64 `json:"f2"`
	F3   float64 `json:"f3"`
	F4   float64 `json:"f4"`
	F5   int64   `json:"f5"`
	F6   float64 `json:"f6"`
	F7   float64 `json:"f7"`
	F8   float64 `json:"f8"`
	F9   float64 `json:"f9"`
	F10  float64 `json:"f10"`
	F11  float64 `json:"f11"`
	F12  string  `json:"f12"`
	F13  int     `json:"f13"`
	F14  string  `json:"f14"`
	F15  float64 `json:"f15"`
	F16  float64 `json:"f16"`
	F17  float64 `json:"f17"`
	F18  float64 `json:"f18"`
	F20  int64   `json:"f20"`
	F21  int64   `json:"f21"`
	F22  float64 `json:"f22"`
	F23  string  `json:"f23"`
	F24  float64 `json:"f24"`
	F25  float64 `json:"f25"`
	F62  float64 `json:"f62"`
	F104 int     `json:"f104"`
	F105 int     `json:"f105"`
	F115 string  `json:"f115"`
	F128 string  `json:"f128"`
	F140 string  `json:"f140"`
	F141 int     `json:"f141"`
	F133 string  `json:"f133"`
	F136 float64 `json:"f136"`
	F152 int     `json:"f152"`
}

type Response struct {
	Rc     int    `json:"rc"`
	Rt     int    `json:"rt"`
	Svr    int    `json:"svr"`
	Lt     int    `json:"lt"`
	Full   int    `json:"full"`
	Dlmkts string `json:"dlmkts"`
	Data   *Data  `json:"data"` // Pointer to handle the possible null value
}

type Data struct {
	Total int    `json:"total"`
	Diff  []Diff `json:"diff"`
}

type FundFlowResponse struct {
	Rc     int       `json:"rc"`     // 响应代码
	Rt     int       `json:"rt"`     // 响应时间
	Svr    int64     `json:"svr"`    // 服务器编号
	Lt     int       `json:"lt"`     // 响应状态
	Full   int       `json:"full"`   // 是否为完整数据
	Dlmkts string    `json:"dlmkts"` // 数据市场信息
	Data   *FlowData `json:"data"`   // 数据内容
}

// FlowData represents the nested data field in the JSON
type FlowData struct {
	Code   string   `json:"code"`   // 数据代码
	Klines []string `json:"klines"` // 行情数据，包含日期和多个浮点值
}

// ParseKlines parses the klines data into a slice of Kline structs
func ParseKlines(secCode string, klines []string) ([]*model.SecFundFlow, error) {
	var result []*model.SecFundFlow
	for _, line := range klines {
		fields := strings.Split(line, ",")
		if len(fields) != 15 {
			return nil, fmt.Errorf("unexpected number of fields in kline data: %v", line)
		}

		// Parse date string to uint64 in YYYYMMDD format
		dateStr := fields[0]
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}
		dateInt64, _ := strconv.ParseInt(parsedDate.Format("20060102"), 10, 64)

		// Convert other fields to appropriate types
		kline := model.SecFundFlow{SecCode: secCode, MarketDate: int32(dateInt64)}
		for i, f := range fields[1:] {
			val, err := strconv.ParseFloat(f, 64)
			if err != nil {
				return nil, err
			}
			switch i {
			case 0:
				kline.MainBuyerNetIn = val
			case 1:
				kline.L4NetInt = val
			case 2:
				kline.L3NetIn = val
			case 3:
				kline.L2NetIn = val
			case 4:
				kline.L1NetIn = val
			case 5:
				kline.MainBuyerRatio = val
			case 6:
				kline.L4NetRatio = val
			case 7:
				kline.L3NetInRatio = float32(val)
			case 8:
				kline.L2NetInRatio = float32(val)
			case 9:
				kline.L1NetInRatio = float32(val)
				// case 10:
				// 	kline.Price = val
				// case 11:
				// 	kline.Percentage = val
				// case 12:
				// 	kline.Additional1 = val
				// case 13:
				// 	kline.Additional2 = val
			}
		}
		result = append(result, &kline)
	}
	return result, nil
}

type TradePeriod struct {
	B int64 `json:"b"`
	E int64 `json:"e"`
}

type TradePeriods struct {
	Pre     interface{}   `json:"pre"`   // null in your example, so using `interface{}` to accommodate any type
	After   interface{}   `json:"after"` // null in your example, so using `interface{}` to accommodate any type
	Periods []TradePeriod `json:"periods"`
}

type FundflowData struct {
	Code   string `json:"code"`
	Market int    `json:"market"`
	// Name         string       `json:"name"`
	// TradePeriods TradePeriods `json:"tradePeriods"`
	Klines []string `json:"klines"`
}

type FundflowResponse struct {
	Rc     int    `json:"rc"`
	Rt     int    `json:"rt"`
	Svr    int    `json:"svr"`
	Lt     int    `json:"lt"`
	Full   int    `json:"full"`
	Dlmkts string `json:"dlmkts"`
	Data   Data   `json:"data"`
}

// Root represents the root structure of the JSON response
type KlineResponse struct {
	Rc     int        `json:"rc"`
	Rt     int        `json:"rt"`
	Svr    int        `json:"svr"`
	Lt     int        `json:"lt"`
	Full   int        `json:"full"`
	Dlmkts string     `json:"dlmkts"`
	Data   *KlineData `json:"data"`
}

// Data represents the "data" field in the JSON response
type KlineData struct {
	Code      string   `json:"code"`
	Market    int      `json:"market"`
	Name      string   `json:"name"`
	Decimal   int      `json:"decimal"`
	Dktotal   int      `json:"dktotal"`
	PreKPrice float64  `json:"preKPrice"`
	Klines    []string `json:"klines"`
}

func ParseKline(stockCode string, marketType uint32, klineStr string) (*model.KlineDay, error) {
	fields := strings.Split(klineStr, ",")
	if len(fields) != 11 {
		return nil, errors.New("invalid kline data format")
	}

	// Parse the date string into a time.Time type
	date, err := time.Parse("2006-01-02", fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	dateInt64, _ := strconv.ParseInt(date.Format("20060102"), 10, 64)

	// Parse remaining fields
	open, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid open value: %v", err)
	}
	closeVal, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid close value: %v", err)
	}
	high, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid high value: %v", err)
	}
	low, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid low value: %v", err)
	}
	volume, err := strconv.ParseInt(fields[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid volume value: %v", err)
	}
	amount, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount value: %v", err)
	}
	amplitude, err := strconv.ParseFloat(fields[7], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amplitude value: %v", err)
	}
	changePct, err := strconv.ParseFloat(fields[8], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid changePct value: %v", err)
	}
	change, err := strconv.ParseFloat(fields[9], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid change value: %v", err)
	}
	turnover, err := strconv.ParseFloat(fields[10], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid turnover value: %v", err)
	}

	return &model.KlineDay{
		StockCode:  stockCode,
		MarketType: int16(marketType),
		MarketDate: int32(dateInt64),
		Open:       open,
		Close:      closeVal,
		High:       high,
		Low:        low,
		Volume:     volume,
		Amount:     amount,
		Amplitude:  amplitude,
		ChangePct:  changePct,
		Change:     change,
		Turnover:   turnover,
	}, nil
}
