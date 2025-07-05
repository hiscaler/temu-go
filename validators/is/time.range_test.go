package is

import (
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestTimeRange(t *testing.T) {
	type object struct {
		StartTime string
		EndTime   string
		Layout    string
		Ok        bool
	}

	objects := []object{
		{"", "", "", false},
		{"2024-01-01", "2024-01-01", time.DateOnly, true},
		{"2024-01-01 12:12:12", "2024-01-01 12:12:12", time.DateTime, true},
		{"2024-01-01 12:12:12", "2024-01-01 12:12:11", time.DateTime, false},
		{"2024-01-01", "2024-01-01 12:12:12", time.DateTime, false},
		{"2024-01-01 0:0:0", "2024-01-01 12:12:12", time.DateTime, false},
		{"2024-01-01", "2024-01-01 12:12:12", time.DateOnly, false},
		{"2024-01-01", "2024-01-02", time.DateOnly, true},
		{"2024-01-02", "2024-01-01", time.DateOnly, false},
		{"qw3e1232024-01-02", "2024-01-01", time.DateOnly, false},
	}

	for _, o := range objects {
		err := validation.Validate("", validation.By(TimeRange(o.StartTime, o.EndTime, o.Layout)))
		assert.Equalf(t, err == nil, o.Ok, "#%s-%s-%s", o.StartTime, o.EndTime, o.Layout)
	}
}
