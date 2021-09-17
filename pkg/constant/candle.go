package constant

type Period int32

const (
	MIN1 Period = iota
	MIN5
	MIN15
	MIN30
	HOUR1
	HOUR4
	DAY1
	WEEK1
	MONTH1
	YEAR1
)

var HuobiPeriodMap = map[Period]string{
	MIN1:   "1min",
	MIN5:   "5min",
	MIN15:  "15min",
	MIN30:  "30min",
	HOUR1:  "60min",
	HOUR4:  "4hour",
	DAY1:   "1day",
	WEEK1:  "1week",
	MONTH1: "1mon",
	YEAR1:  "1year",
}

var HuobiPeriodNameMap = map[string]Period{
	"1min":  MIN1,
	"5min":  MIN5,
	"15min": MIN15,
	"30min": MIN30,
	"60min": HOUR1,
	"4hour": HOUR4,
	"1day":  DAY1,
	"1week": WEEK1,
	"1mon":  MONTH1,
	"1year": YEAR1,
}
