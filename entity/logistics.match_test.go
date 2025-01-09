package entity

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLogisticsMatch_LatestScheduleTime(t *testing.T) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format(time.DateOnly)
	tests := []struct {
		name      string
		schedules []LogisticsMatchChannelScheduleTime
		hasErr    bool
	}{
		{
			name: "test1",
			schedules: []LogisticsMatchChannelScheduleTime{
				{
					BjDate:      tomorrow,
					BjStartTime: "09:00",
					BjEndTime:   "12:00",
				},
				{
					BjDate:      "2024-01-02",
					BjStartTime: "09:00",
					BjEndTime:   "12:00",
				},
				{
					BjDate:      "2025-01-02",
					BjStartTime: "09:00",
					BjEndTime:   "12:00",
				},
			},
			hasErr: false,
		},
		{
			name: "test2",
			schedules: []LogisticsMatchChannelScheduleTime{
				{
					BjDate:      tomorrow,
					BjStartTime: "09:00",
					BjEndTime:   "12:00",
				},
				{
					BjDate:      "2024-01-02",
					BjStartTime: "10:00",
					BjEndTime:   "12:00",
				},
				{
					BjDate:      "2024-01-02",
					BjStartTime: "09:00",
					BjEndTime:   "12:00",
				},
				{
					BjDate:      "2025-01-02", // start.After(now) && end.Before(now) {
					BjStartTime: "09:00",
					BjEndTime:   "12:00",
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		m := LogisticsMatch{ChannelScheduleTimeList: tt.schedules}
		_, err := m.LatestScheduleTime()
		assert.Equal(t, tt.hasErr, err != nil, tt.name)
		fmt.Print(err)
	}
}
