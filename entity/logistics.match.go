package entity

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

type LogisticsMatchChannelScheduleTime struct {
	BjDate      string `json:"bjDate"`
	BjStartTime string `json:"bjStartTime"`
	BjEndTime   string `json:"bjEndTime"`
}

// LogisticsChannelAppointmentTime 物流渠道预约时间
type LogisticsChannelAppointmentTime struct {
	Start time.Time
	End   time.Time
}

// Range 时间范围
func (st LogisticsMatchChannelScheduleTime) Range() (t LogisticsChannelAppointmentTime, err error) {
	layout := "2006-01-02 15:04"
	start, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", st.BjDate, st.BjStartTime), time.Local)
	if err != nil {
		return
	}

	end, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", st.BjDate, st.BjEndTime), time.Local)
	if err != nil {
		return
	}

	return LogisticsChannelAppointmentTime{Start: start, End: end}, nil
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
func (lm LogisticsMatch) LatestScheduleTime() (t LogisticsChannelAppointmentTime, err error) {
	if len(lm.ChannelScheduleTimeList) == 0 {
		return
	}

	schedules := lm.ChannelScheduleTimeList
	sort.Slice(schedules, func(i, j int) bool {
		return (schedules[i].BjDate + schedules[i].BjStartTime) < (schedules[j].BjDate + schedules[j].BjStartTime)
	})
	now := time.Now()
	for _, scheduleTime := range schedules {
		st, e := scheduleTime.Range()
		if e != nil {
			return t, e
		}

		if (now.Before(st.Start) || now.Equal(st.Start)) && (now.Before(st.End) || now.Equal(st.End)) {
			return st, nil
		}
	}
	return t, errors.New("not exists")
}
