package carbon

import (
	"context"
	"testing"

	"github.com/golang-module/carbon/v2"
	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestModule(t *testing.T) {
	m := Module()
	require.NotNil(t, m)

	parseObj, ok := m.GetAttr("parse")
	require.True(t, ok)
	fn, ok := parseObj.(*object.Builtin)
	require.True(t, ok)
	result := fn.Call(context.Background(), object.NewString("2021-03-14"))
	require.NotNil(t, result)

	carbonObj, ok := result.(*Carbon)
	require.True(t, ok)
	require.Equal(t, "2021-03-14 00:00:00", carbonObj.Value().String())

	fnObj, ok := carbonObj.GetAttr("add_days")
	require.True(t, ok)
	addDays, ok := fnObj.(*object.Builtin)
	require.True(t, ok)

	result = addDays.Call(context.Background(), object.NewInt(3))
	require.NotNil(t, result)
	require.Equal(t, "carbon.carbon(2021-03-17 00:00:00)", result.Inspect())
}

func TestCarbonMethods(t *testing.T) {

	c := NewCarbon(carbon.Parse("2021-03-14"))

	testCases := []struct {
		name     string
		method   string
		args     []object.Object
		expected string
	}{
		{
			name:     "add_day",
			method:   "add_day",
			args:     []object.Object{},
			expected: "carbon.carbon(2021-03-15 00:00:00)",
		},
		{
			name:     "sub_day",
			method:   "sub_day",
			args:     []object.Object{},
			expected: "carbon.carbon(2021-03-13 00:00:00)",
		},
		{
			name:     "sub_days",
			method:   "sub_days",
			args:     []object.Object{object.NewInt(3)},
			expected: "carbon.carbon(2021-03-11 00:00:00)",
		},
		{
			name:     "days_in_month",
			method:   "days_in_month",
			args:     []object.Object{},
			expected: "31",
		},
		{
			name:     "days_in_year",
			method:   "days_in_year",
			args:     []object.Object{},
			expected: "365",
		},
		{
			name:     "week_of_month",
			method:   "week_of_month",
			args:     []object.Object{},
			expected: "2",
		},
		{
			name:     "week_of_year",
			method:   "week_of_year",
			args:     []object.Object{},
			expected: "10",
		},
		{
			name:     "timestamp",
			method:   "timestamp",
			args:     []object.Object{},
			expected: "1615698000",
		},
		{
			name:     "is_am",
			method:   "is_am",
			args:     []object.Object{},
			expected: "true",
		},
		{
			name:     "is_pm",
			method:   "is_pm",
			args:     []object.Object{},
			expected: "false",
		},
		{
			name:     "is_leap_year",
			method:   "is_leap_year",
			args:     []object.Object{},
			expected: "false",
		},
		{
			name:     "is_future",
			method:   "is_future",
			args:     []object.Object{},
			expected: "false",
		},
		{
			name:     "is_past",
			method:   "is_past",
			args:     []object.Object{},
			expected: "true",
		},
		{
			name:     "is_today",
			method:   "is_today",
			args:     []object.Object{},
			expected: "false",
		},
		{
			name:     "is_yesterday",
			method:   "is_yesterday",
			args:     []object.Object{},
			expected: "false",
		},
		{
			name:     "day",
			method:   "day",
			args:     []object.Object{},
			expected: "14",
		},
		{
			name:     "month",
			method:   "month",
			args:     []object.Object{},
			expected: "3",
		},
		{
			name:     "year",
			method:   "year",
			args:     []object.Object{},
			expected: "2021",
		},
		{
			name:     "to_date_string",
			method:   "to_date_string",
			args:     []object.Object{},
			expected: "\"2021-03-14\"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fnObj, ok := c.GetAttr(tc.method)
			require.True(t, ok)
			fn, ok := fnObj.(*object.Builtin)
			require.True(t, ok)
			result := fn.Call(context.Background(), tc.args...)
			require.Equal(t, tc.expected, result.Inspect())
		})
	}
}
