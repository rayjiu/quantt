package crawler

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
