package entity

type LogisticsMatchChannelScheduleTime struct {
	BjDate      string `json:"bjDate"`
	BjStartTime string `json:"bjStartTime"`
	BjEndTime   string `json:"bjEndTime"`
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
