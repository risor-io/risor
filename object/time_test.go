package object

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeAddDate(t *testing.T) {
	baseTime := NewTime(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC))

	tests := []struct {
		years   int64
		months  int64
		days    int64
		want    time.Time
		wantErr bool
	}{
		{1, 0, 0, time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{0, 1, 0, time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC), false},
		{0, 0, 1, time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), false},
		{1, 1, 1, time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC), false},
		{-1, 0, 0, time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{-1, 0, -1, time.Date(2022, 9, 30, 0, 0, 0, 0, time.UTC), false},
		{0, 0, 0, time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("AddDate(%d, %d, %d)", tt.years, tt.months, tt.days), func(t *testing.T) {
			result := baseTime.AddDate(context.Background(), NewInt(tt.years), NewInt(tt.months), NewInt(tt.days))

			if tt.wantErr {
				_, ok := result.(*Error)
				assert.True(t, ok, "expected error, got %v", result)
				return
			}

			assert.Equal(t, TIME, result.Type(), "expected TIME, got %v", result.Type())

			resultTime := result.(*Time).Value()
			assert.True(t, resultTime.Equal(tt.want), "expected %v, got %v", tt.want, resultTime)
		})
	}
}
