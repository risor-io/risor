package carbon

import (
	"context"
	"fmt"

	"github.com/golang-module/carbon/v2"
	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

var _ object.Object = (*Carbon)(nil)

const CARBON object.Type = "carbon.carbon"

type Carbon struct {
	value carbon.Carbon
}

func (c *Carbon) Value() carbon.Carbon {
	return c.value
}

func (c *Carbon) IsTruthy() bool {
	return true
}

func (c *Carbon) Type() object.Type {
	return CARBON
}

func (c *Carbon) Inspect() string {
	return fmt.Sprintf("%s(%s)", CARBON, c.value.String())
}

func (c *Carbon) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: cannot set %q on %s object", name, CARBON)
}

func (c *Carbon) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "add_day":
		return object.NewBuiltin("add_day", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.add_day", 0, args); err != nil {
				return err
			}
			return NewCarbon(c.value.AddDay())
		}), true
	case "add_days":
		return object.NewBuiltin("add_days", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.add_days", 1, args); err != nil {
				return err
			}
			i, inputErr := object.AsInt(args[0])
			if inputErr != nil {
				return inputErr
			}
			return NewCarbon(c.value.AddDays(int(i)))
		}), true
	case "sub_day":
		return object.NewBuiltin("sub_day", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.sub_day", 0, args); err != nil {
				return err
			}
			return NewCarbon(c.value.SubDay())
		}), true
	case "sub_days":
		return object.NewBuiltin("sub_days", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.sub_days", 1, args); err != nil {
				return err
			}
			i, inputErr := object.AsInt(args[0])
			if inputErr != nil {
				return inputErr
			}
			return NewCarbon(c.value.SubDays(int(i)))
		}), true
	case "timezone":
		return object.NewBuiltin("timezone", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.timezone", 0, args); err != nil {
				return err
			}
			return object.NewString(c.value.Timezone())
		}), true
	case "age":
		return object.NewBuiltin("age", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.age", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.Age()))
		}), true
	case "days_in_month":
		return object.NewBuiltin("days_in_month", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.days_in_month", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.DaysInMonth()))
		}), true
	case "days_in_year":
		return object.NewBuiltin("days_in_year", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.days_in_year", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.DaysInYear()))
		}), true
	case "week_of_month":
		return object.NewBuiltin("week_of_month", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.week_of_month", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.WeekOfMonth()))
		}), true
	case "week_of_year":
		return object.NewBuiltin("week_of_year", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.week_of_year", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.WeekOfYear()))
		}), true
	case "timestamp":
		return object.NewBuiltin("timestamp", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.timestamp", 0, args); err != nil {
				return err
			}
			return object.NewInt(c.value.Timestamp())
		}), true
	case "string":
		return object.NewBuiltin("string", func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewString(c.value.String())
		}), true
	case "is_valid":
		return object.NewBuiltin("is_valid", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.is_valid", 0, args); err != nil {
				return err
			}
			return object.NewBool(c.value.IsValid())
		}), true
	case "is_am":
		return object.NewBuiltin("is_am", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.is_am", 0, args); err != nil {
				return err
			}
			return object.NewBool(c.value.IsAM())
		}), true
	case "is_pm":
		return object.NewBuiltin("is_pm", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.is_pm", 0, args); err != nil {
				return err
			}
			return object.NewBool(c.value.IsPM())
		}), true
	case "is_leap_year":
		return object.NewBuiltin("is_leap_year", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.is_leap_year", 0, args); err != nil {
				return err
			}
			return object.NewBool(c.value.IsLeapYear())
		}), true
	case "is_future":
		return object.NewBuiltin("is_future", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.is_future", 0, args); err != nil {
				return err
			}
			return object.NewBool(c.value.IsFuture())
		}), true
	case "is_past":
		return object.NewBuiltin("is_past", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.is_past", 0, args); err != nil {
				return err
			}
			return object.NewBool(c.value.IsPast())
		}), true
	case "is_today":
		return object.NewBuiltin("is_today", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.is_today", 0, args); err != nil {
				return err
			}
			return object.NewBool(c.value.IsToday())
		}), true
	case "is_yesterday":
		return object.NewBuiltin("is_yesterday", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.is_yesterday", 0, args); err != nil {
				return err
			}
			return object.NewBool(c.value.IsYesterday())
		}), true
	case "day":
		return object.NewBuiltin("day", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.day", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.Day()))
		}), true
	case "month":
		return object.NewBuiltin("month", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.month", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.Month()))
		}), true
	case "year":
		return object.NewBuiltin("year", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.year", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.Year()))
		}), true
	case "hour":
		return object.NewBuiltin("hour", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.hour", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.Hour()))
		}), true
	case "minute":
		return object.NewBuiltin("minute", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.minute", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.Minute()))
		}), true
	case "second":
		return object.NewBuiltin("second", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.second", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(c.value.Second()))
		}), true
	case "diff_for_humans":
		return object.NewBuiltin("diff_for_humans", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.diff_for_humans", 1, args); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *Carbon:
				return object.NewString(c.value.DiffForHumans(arg.value))
			case *object.Time:
				return object.NewString(c.value.DiffForHumans(carbon.CreateFromStdTime(arg.Value())))
			default:
				return object.TypeErrorf("type error: expected carbon (got %s)", args[0].Type())
			}
		}), true
	case "std_time":
		return object.NewBuiltin("std_time", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.std_time", 0, args); err != nil {
				return err
			}
			return object.NewTime(c.value.StdTime())
		}), true
	case "to_date_string":
		return object.NewBuiltin("to_date_string", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.to_date_string", 0, args); err != nil {
				return err
			}
			return object.NewString(c.value.ToDateString())
		}), true
	case "to_time_string":
		return object.NewBuiltin("to_time_string", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.to_time_string", 0, args); err != nil {
				return err
			}
			return object.NewString(c.value.ToTimeString())
		}), true
	case "to_datetime_string":
		return object.NewBuiltin("to_datetime_string", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("carbon.to_datetime_string", 0, args); err != nil {
				return err
			}
			return object.NewString(c.value.ToDateTimeString())
		}), true
	default:
		return nil, false
	}
}

func (c *Carbon) String() string {
	return c.value.String()
}

func (c *Carbon) Interface() interface{} {
	return c.value
}

func (c *Carbon) Equals(other object.Object) object.Object {
	otherCarbon, ok := other.(*Carbon)
	if !ok {
		return object.False
	}
	return object.NewBool(otherCarbon.Value().Eq(c.value))
}

func (c *Carbon) Cost() int {
	return 0
}

func (c *Carbon) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.EvalErrorf("eval error: unsupported operation for %s: %v", CARBON, opType)
}

func NewCarbon(v carbon.Carbon) *Carbon {
	return &Carbon{value: v}
}
