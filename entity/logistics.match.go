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
	start, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", st.BjDate, st.BjStartTime), loc)
	if err != nil {
		return
	}

	end, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", st.BjDate, st.BjEndTime), loc)

	return
}

// LogisticsMatch 推荐物流商匹配
type LogisticsMatch struct {
	MaxChargeAmount         float64                             `json:"maxChargeAmount"`
	PredictId               int64                               `json:"predictId"`
	MaxSupplierChargeAmount float64                             `json:"maxSupplierChargeAmount"`
	StandbyExpress          bool                                `json:"standbyExpress"`
	MinSupplierChargeAmount float64                             `json:"minSupplierChargeAmount"`
	TmsChannelId            int                                 `json:"tmsChannelId"`
	LatestAppointmentTime   int                                 `json:"latestAppointmentTime"`
	ExpressCompanyId        int                                 `json:"expressCompanyId"`
	MinChargeAmount         float64                             `json:"minChargeAmount"`
	ExpressCompanyName      string                              `json:"expressCompanyName"`
	ChannelScheduleTimeList []LogisticsMatchChannelScheduleTime `json:"channelScheduleTimeList"`
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
