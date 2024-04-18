package color

import (
	"context"

	"github.com/fatih/color"
	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func Set(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("color.set", 1, 64, args); err != nil {
		return err
	}
	var attrs []color.Attribute
	for _, arg := range args {
		colorCode, err := object.AsInt(arg)
		if err != nil {
			return err
		}
		attrs = append(attrs, color.Attribute(colorCode))
	}
	color.Set(attrs...)
	return object.Nil
}

func Unset(ctx context.Context, args ...object.Object) object.Object {
	color.Unset()
	return object.Nil
}

func CreateColor(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("color.color", 1, 64, args); err != nil {
		return err
	}
	var attrs []color.Attribute
	for _, arg := range args {
		colorCode, err := object.AsInt(arg)
		if err != nil {
			return err
		}
		attrs = append(attrs, color.Attribute(colorCode))
	}
	return NewColor(color.New(attrs...))
}

func Module() *object.Module {
	return object.NewBuiltinsModule("color", map[string]object.Object{
		"color":        object.NewBuiltin("color", CreateColor),
		"set":          object.NewBuiltin("set", Set),
		"unset":        object.NewBuiltin("unset", Unset),
		"reset":        object.NewInt(int64(color.Reset)),
		"bold":         object.NewInt(int64(color.Bold)),
		"dim":          object.NewInt(int64(color.Faint)),
		"italic":       object.NewInt(int64(color.Italic)),
		"underline":    object.NewInt(int64(color.Underline)),
		"blinkslow":    object.NewInt(int64(color.BlinkSlow)),
		"blinkrapid":   object.NewInt(int64(color.BlinkRapid)),
		"reversevideo": object.NewInt(int64(color.ReverseVideo)),
		"concealed":    object.NewInt(int64(color.Concealed)),
		"crossedout":   object.NewInt(int64(color.CrossedOut)),
		"bg_black":     object.NewInt(int64(color.BgBlack)),
		"bg_blue":      object.NewInt(int64(color.BgBlue)),
		"bg_cyan":      object.NewInt(int64(color.BgCyan)),
		"bg_green":     object.NewInt(int64(color.BgGreen)),
		"bg_hiblack":   object.NewInt(int64(color.BgHiBlack)),
		"bg_hiblue":    object.NewInt(int64(color.BgHiBlue)),
		"bg_hicyan":    object.NewInt(int64(color.BgHiCyan)),
		"bg_higreen":   object.NewInt(int64(color.BgHiGreen)),
		"bg_himagenta": object.NewInt(int64(color.BgHiMagenta)),
		"bg_hired":     object.NewInt(int64(color.BgHiRed)),
		"bg_hiwhite":   object.NewInt(int64(color.BgHiWhite)),
		"bg_hiyellow":  object.NewInt(int64(color.BgHiYellow)),
		"bg_magenta":   object.NewInt(int64(color.BgMagenta)),
		"bg_red":       object.NewInt(int64(color.BgRed)),
		"bg_white":     object.NewInt(int64(color.BgWhite)),
		"bg_yellow":    object.NewInt(int64(color.BgYellow)),
		"fg_black":     object.NewInt(int64(color.FgBlack)),
		"fg_blue":      object.NewInt(int64(color.FgBlue)),
		"fg_cyan":      object.NewInt(int64(color.FgCyan)),
		"fg_green":     object.NewInt(int64(color.FgGreen)),
		"fg_hiblack":   object.NewInt(int64(color.FgHiBlack)),
		"fg_hiblue":    object.NewInt(int64(color.FgHiBlue)),
		"fg_hicyan":    object.NewInt(int64(color.FgHiCyan)),
		"fg_higreen":   object.NewInt(int64(color.FgHiGreen)),
		"fg_himagenta": object.NewInt(int64(color.FgHiMagenta)),
		"fg_hired":     object.NewInt(int64(color.FgHiRed)),
		"fg_hiwhite":   object.NewInt(int64(color.FgHiWhite)),
		"fg_hiyellow":  object.NewInt(int64(color.FgHiYellow)),
		"fg_magenta":   object.NewInt(int64(color.FgMagenta)),
		"fg_red":       object.NewInt(int64(color.FgRed)),
		"fg_white":     object.NewInt(int64(color.FgWhite)),
		"fg_yellow":    object.NewInt(int64(color.FgYellow)),
	}, CreateColor)
}
