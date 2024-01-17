package regexp

import (
	"context"
	"fmt"
	"regexp"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const REGEXP object.Type = "regexp"

type Regexp struct {
	value *regexp.Regexp
}

func (r *Regexp) Type() object.Type {
	return REGEXP
}

func (r *Regexp) Inspect() string {
	return fmt.Sprintf("regexp(%q)", r.value.String())
}

func (r *Regexp) String() string {
	return r.Inspect()
}

func (r *Regexp) Interface() interface{} {
	return r.value
}

func (r *Regexp) HashKey() object.HashKey {
	return object.HashKey{Type: r.Type(), StrValue: r.value.String()}
}

func (r *Regexp) Compare(other object.Object) (int, error) {
	typeComp := object.CompareTypes(r, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherRegex := other.(*Regexp)
	if r.value == otherRegex.value {
		return 0, nil
	}
	if r.value.String() > otherRegex.value.String() {
		return 1, nil
	}
	return -1, nil
}

func (r *Regexp) Equals(other object.Object) object.Object {
	switch other := other.(type) {
	case *Regexp:
		if r.value == other.value {
			return object.True
		}
	}
	return object.False
}

func (r *Regexp) MarshalJSON() ([]byte, error) {
	return []byte(r.value.String()), nil
}

func (r *Regexp) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for regexp: %v", opType)
}

func (r *Regexp) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: cannot set attribute %q on regexp object", name)
}

func (r *Regexp) IsTruthy() bool {
	return true
}

func (r *Regexp) Cost() int {
	return 0
}

func (r *Regexp) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "match":
		return object.NewBuiltin("regexp.match",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewArgsError("regexp.match", 1, len(args))
				}
				strValue, err := object.AsString(args[0])
				if err != nil {
					return err
				}
				return object.NewBool(r.value.MatchString(strValue))
			},
		), true
	case "find":
		return object.NewBuiltin("regexp.find",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewArgsError("regexp.find", 1, len(args))
				}
				strValue, err := object.AsString(args[0])
				if err != nil {
					return err
				}
				return object.NewString(r.value.FindString(strValue))
			},
		), true
	case "find_all":
		return object.NewBuiltin("regexp.find_all",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) < 1 || len(args) > 2 {
					return object.NewArgsRangeError("regexp.find_all", 1, 2, len(args))
				}
				strValue, err := object.AsString(args[0])
				if err != nil {
					return err
				}
				n := -1
				if len(args) == 2 {
					i64, err := object.AsInt(args[1])
					if err != nil {
						return err
					}
					n = int(i64)
				}
				var matches []object.Object
				for _, match := range r.value.FindAllString(strValue, n) {
					matches = append(matches, object.NewString(match))
				}
				return object.NewList(matches)
			},
		), true
	case "find_submatch":
		return object.NewBuiltin("regexp.find_submatch",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewArgsError("regexp.find_submatch", 1, len(args))
				}
				strValue, err := object.AsString(args[0])
				if err != nil {
					return err
				}
				var matches []object.Object
				for _, match := range r.value.FindStringSubmatch(strValue) {
					matches = append(matches, object.NewString(match))
				}
				return object.NewList(matches)
			},
		), true
	case "replace_all":
		return object.NewBuiltin("regexp.replace_all",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewArgsError("regexp.replace_all", 2, len(args))
				}
				strValue, err := object.AsString(args[0])
				if err != nil {
					return err
				}
				replaceValue, err := object.AsString(args[1])
				if err != nil {
					return err
				}
				return object.NewString(r.value.ReplaceAllString(strValue, replaceValue))
			},
		), true
	case "split":
		return object.NewBuiltin("regexp.split",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) < 1 || len(args) > 2 {
					return object.NewArgsRangeError("regexp.split", 1, 2, len(args))
				}
				strValue, err := object.AsString(args[0])
				if err != nil {
					return err
				}
				n := -1
				if len(args) == 2 {
					i64, err := object.AsInt(args[1])
					if err != nil {
						return err
					}
					n = int(i64)
				}
				matches := r.value.Split(strValue, n)
				matchObjects := make([]object.Object, 0, len(matches))
				for _, match := range matches {
					matchObjects = append(matchObjects, object.NewString(match))
				}
				return object.NewList(matchObjects)
			},
		), true
	}
	return nil, false
}

func NewRegexp(value *regexp.Regexp) *Regexp {
	return &Regexp{value: value}
}
