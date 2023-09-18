package object

import (
	"context"
	"fmt"
	"regexp"

	"github.com/risor-io/risor/op"
)

type Regexp struct {
	*base
	value *regexp.Regexp
}

func (r *Regexp) Type() Type {
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

func (r *Regexp) HashKey() HashKey {
	return HashKey{Type: r.Type(), StrValue: r.value.String()}
}

func (r *Regexp) Compare(other Object) (int, error) {
	typeComp := CompareTypes(r, other)
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

func (r *Regexp) Equals(other Object) Object {
	switch other := other.(type) {
	case *Regexp:
		if r.value == other.value {
			return True
		}
	}
	return False
}

func (r *Regexp) MarshalJSON() ([]byte, error) {
	return []byte(r.value.String()), nil
}

func (r *Regexp) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("eval error: unsupported operation for regexp: %v", opType))
}

func (r *Regexp) GetAttr(name string) (Object, bool) {
	switch name {
	case "match":
		return &Builtin{
			name: "regexp.match",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("regexp.match", 1, len(args))
				}
				strValue, err := AsString(args[0])
				if err != nil {
					return err
				}
				return NewBool(r.value.MatchString(strValue))
			},
		}, true
	case "find":
		return &Builtin{
			name: "regexp.find",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("regexp.find", 1, len(args))
				}
				strValue, err := AsString(args[0])
				if err != nil {
					return err
				}
				return NewString(r.value.FindString(strValue))
			},
		}, true
	case "find_all":
		return &Builtin{
			name: "regexp.find_all",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) < 1 || len(args) > 2 {
					return NewArgsRangeError("regexp.find_all", 1, 2, len(args))
				}
				strValue, err := AsString(args[0])
				if err != nil {
					return err
				}
				n := -1
				if len(args) == 2 {
					i64, err := AsInt(args[1])
					if err != nil {
						return err
					}
					n = int(i64)
				}
				var matches []Object
				for _, match := range r.value.FindAllString(strValue, n) {
					matches = append(matches, NewString(match))
				}
				return NewList(matches)
			},
		}, true
	case "find_submatch":
		return &Builtin{
			name: "regexp.find_submatch",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 1 {
					return NewArgsError("regexp.find_submatch", 1, len(args))
				}
				strValue, err := AsString(args[0])
				if err != nil {
					return err
				}
				var matches []Object
				for _, match := range r.value.FindStringSubmatch(strValue) {
					matches = append(matches, NewString(match))
				}
				return NewList(matches)
			},
		}, true
	case "replace_all":
		return &Builtin{
			name: "regexp.replace_all",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) != 2 {
					return NewArgsError("regexp.replace_all", 2, len(args))
				}
				strValue, err := AsString(args[0])
				if err != nil {
					return err
				}
				replaceValue, err := AsString(args[1])
				if err != nil {
					return err
				}
				return NewString(r.value.ReplaceAllString(strValue, replaceValue))
			},
		}, true
	case "split":
		return &Builtin{
			name: "regexp.split",
			fn: func(ctx context.Context, args ...Object) Object {
				if len(args) < 1 || len(args) > 2 {
					return NewArgsRangeError("regexp.split", 1, 2, len(args))
				}
				strValue, err := AsString(args[0])
				if err != nil {
					return err
				}
				n := -1
				if len(args) == 2 {
					i64, err := AsInt(args[1])
					if err != nil {
						return err
					}
					n = int(i64)
				}
				matches := r.value.Split(strValue, n)
				matchObjects := make([]Object, 0, len(matches))
				for _, match := range matches {
					matchObjects = append(matchObjects, NewString(match))
				}
				return NewList(matchObjects)
			},
		}, true
	}
	return nil, false
}

func NewRegexp(value *regexp.Regexp) *Regexp {
	return &Regexp{value: value}
}
