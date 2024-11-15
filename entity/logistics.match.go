package entity

import (
	"fmt"
	"time"
)

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
}

type LogisticsMatchChannelScheduleTime struct {
	BjDate      string `json:"bjDate"`
	BjStartTime string `json:"bjStartTime"`
	BjEndTime   string `json:"bjEndTime"`
}

// Range 时间范围
func (st LogisticsMatchChannelScheduleTime) Range() (start, end time.Time, err error) {
	layout := "2006-01-02 15:04"
	start, err = time.ParseInLocation(layout, fmt.Sprintf("%s %s", st.BjDate, st.BjStartTime), loc)
	if err != nil {
		return
	}

	end, err = time.ParseInLocation(layout, fmt.Sprintf("%s %s", st.BjDate, st.BjEndTime), loc)

	return
}

// LogisticsMatch 推荐物流商匹配
type LogisticsMatch struct {
	ExpressCompanyId        int                                 `json:"expressCompanyId"`        // 快递公司 ID
	ExpressCompanyName      string                              `json:"expressCompanyName"`      // 快递公司名称
	StandbyExpress          bool                                `json:"standbyExpress"`          // 是否是备用快递公司
	TmsChannelId            int                                 `json:"tmsChannelId"`            // TMS 快递产品类型 ID
	MinSupplierChargeAmount float64                             `json:"minSupplierChargeAmount"` // 最小预估商家承担运费（单位元）
	MaxSupplierChargeAmount float64                             `json:"maxSupplierChargeAmount"` // 最大预估商家承担运费（单位元）
	MinChargeAmount         float64                             `json:"minChargeAmount"`         // 最小预估运费（单位元）
	MaxChargeAmount         float64                             `json:"maxChargeAmount"`         // 最小预估运费（单位元）
	ChannelScheduleTimeList []LogisticsMatchChannelScheduleTime `json:"channelScheduleTimeList"` // 可预约揽收时间
	LatestAppointmentTime   int64                               `json:"latestAppointmentTime"`   // 最迟预约时间
	PredictId               int64                               `json:"predictId"`               // 预测 ID
}

// LatestScheduleTime 获取最近可用的时间
func (lm LogisticsMatch) LatestScheduleTime() *LogisticsMatchChannelScheduleTime {
	if len(lm.ChannelScheduleTimeList) == 0 {
		return nil
	}

	now := time.Now()
	for _, scheduleTime := range lm.ChannelScheduleTimeList {
		start, end, err := scheduleTime.Range()
		if err != nil {
			return nil
		}

		if start.After(now) && end.Before(now) {
			return &scheduleTime
		}
	}

	return nil
}
